package main

import (
	"echo-server/db"
	"echo-server/handlers"
	"echo-server/models"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return cv.customMessage(i, err)
	}
	return nil
}

var customErrorMessages = map[string]string{
	"Student.Name.required":   "Name is required",
	"Student.Name.min":        "Name must be at least 2 characters long",
	"Student.Name.max":        "Name cannot be longer than 100 characters",
	"Student.Name.alphaspace": "Name can only contain alphabetic characters and spaces",
	"Student.Age.required":    "Age is required",
	"Student.Age.min":         "Age must be greater than 1",
	"Student.Age.max":         "Age cannot be greater than 120",
	"Student.Email.required":  "Email is required",
	"Student.Email.email":     "Email must be a valid email address",
}

func (cv *CustomValidator) customMessage(i any, err error) error {
	var errMessages []string
	structName := reflect.TypeOf(i).Elem().Name()
	log.Println("Struct Name:", structName)

	for _, e := range err.(validator.ValidationErrors) {
		key := fmt.Sprintf("%s.%s.%s", structName, e.Field(), e.Tag())
		if msg, exists := customErrorMessages[key]; exists {
			errMessages = append(errMessages, msg)
			continue
		}
		errMessages = append(errMessages, fmt.Sprintf("%s failed on %s: ", e.Field(), e.Tag()))
	}

	return fmt.Errorf(strings.Join(errMessages, "; "))

}

func main() {
	e := echo.New()

	e.Validator = &CustomValidator{validator: validator.New()}

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
