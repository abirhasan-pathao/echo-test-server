package cmd

import (
	"context"
	cookiedelivery "echo-server/app/cookie/delivery"
	cookieusecase "echo-server/app/cookie/usecase"
	"echo-server/app/students/delivery"
	"echo-server/app/students/model"
	"echo-server/app/students/repository"
	"echo-server/app/students/usecase"
	"echo-server/app/utils/validation"
	"echo-server/config"
	"echo-server/infrastructure/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/spf13/cobra"
)

var startServerCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the Echo Server",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.LoadConfig()
		if err != nil {
			log.Fatal("Failed to load config: ", err)
		}

		dbConn := db.NewPostgresDB(config.Database)
		dbConn.AutoMigrate(&model.Student{})

		studentRepo := repository.NewStudentRepository(dbConn)
		studentUsecase := usecase.NewStudentUsecase(studentRepo)
		studentController := delivery.NewStudentController(studentUsecase)

		cookieUsecase := cookieusecase.NewCookieUsecase()
		cookieController := cookiedelivery.NewCookieController(cookieUsecase)

		e := echo.New()
		e.Use(middleware.RequestLogger())

		validation.AddValidator()

		e.GET("/", func(c *echo.Context) error {
			return c.JSON(http.StatusOK, "Server is Running...")
		})

		e.POST("/students", studentController.CreateStudent)
		e.GET("/students/:id", studentController.GetStudentByID)
		e.GET("/students", studentController.GetAllStudents)
		e.PUT("/students/:id", studentController.UpdateStudent)
		e.DELETE("/students/:id", studentController.DeleteStudent)
		e.GET("/students/average-age", studentController.GetAverageAge)

		e.GET("/cookie", cookieController.GetCookie)
		e.GET("/milk", cookieController.GetMilk)

		srv := &http.Server{
			Addr:    ":" + config.Environment.Port,
			Handler: e,
		}

		go func() {
			log.Println("Starting " + config.Environment.Env + " server on :" + config.Environment.Port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server error: %v", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

		sqlDB, err := dbConn.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("Database connection closed")
		}

		log.Println("Server exited gracefully")

	},
}

func init() {
	rootCmd.AddCommand(startServerCmd)
}
