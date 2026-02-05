package cmd

import (
	"echo-server/app/students/delivery"
	"echo-server/app/students/repository"
	"echo-server/app/students/usecase"
	"echo-server/app/utils"
	"echo-server/config"
	"echo-server/infrastructure/db"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
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
		studentRepo := repository.NewStudentRepository(dbConn)
		studentUsecase := usecase.NewStudentUsecase(studentRepo)
		studentController := delivery.NewStudentController(studentUsecase)

		e := echo.New()
		e.Use(middleware.RequestLogger())

		e.Validator = &utils.CustomValidator{Validator: validator.New()}

		e.GET("/", func(c *echo.Context) error {
			return c.JSON(http.StatusOK, "Server is Running...")
		})

		e.POST("/students", studentController.CreateStudent)
		e.GET("/students/:id", studentController.GetStudentByID)
		e.GET("/students", studentController.GetAllStudents)
		e.PUT("/students/:id", studentController.UpdateStudent)
		e.DELETE("/students/:id", studentController.DeleteStudent)
		e.GET("/students/average-age", studentController.GetAverageAge)

		log.Println("Starting " + config.Environment.Env + " server on :" + config.Environment.Port)
		log.Fatal(e.Start(":" + config.Environment.Port))
	},
}

func init() {
	rootCmd.AddCommand(startServerCmd)
}
