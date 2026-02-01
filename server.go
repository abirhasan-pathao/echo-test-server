package main

import (
	"echo-server/db"
	"echo-server/handlers"
	"echo-server/models"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Server Running...")
	})

	db.Connect()
	db.DB.AutoMigrate(&models.Student{})

	e.POST("/students", handlers.CreateStudent)
	e.GET("/students", handlers.GetAllStudents)
	e.GET("/students/:id", handlers.GetStudentByID)
	e.PUT("/students/:id", handlers.UpdateStudent)
	e.DELETE("/students/:id", handlers.DeleteStudent)

	e.Logger.Info("Server Starting on Port 8080")

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("Shutting down the server", "error", err)
	}

}
