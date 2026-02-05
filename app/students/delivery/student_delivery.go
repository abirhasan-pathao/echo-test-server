package delivery

import (
	"echo-server/app/students/model"
	"echo-server/app/students/usecase"
	"net/http"
	"strconv"
	"strings"

	"github.com/gookit/validate"
	"github.com/labstack/echo/v5"
)

type StudentController struct {
	studentUsecase *usecase.StudentUsecase
}

func NewStudentController(studentUsecase *usecase.StudentUsecase) *StudentController {
	return &StudentController{studentUsecase: studentUsecase}
}

func (ctrl *StudentController) CreateStudent(c *echo.Context) error {
	var student model.Student
	if err := c.Bind(&student); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	v := validate.Struct(student)
	v.StopOnError = false
	if !v.Validate() {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "validation failed",
			"errors":  v.Errors.All(),
		})
	}

	err := ctrl.studentUsecase.CreateStudent(&student)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Email already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, student)
}

func (ctrl *StudentController) GetStudentByID(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid student ID"})
	}
	student, err := ctrl.studentUsecase.GetStudentByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Student not found"})
	}
	return c.JSON(http.StatusOK, student)
}

func (ctrl *StudentController) GetAllStudents(c *echo.Context) error {
	students, err := ctrl.studentUsecase.GetAllStudents()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, students)
}

func (ctrl *StudentController) UpdateStudent(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid student ID"})
	}

	var updatedStudent model.Student
	if err := c.Bind(&updatedStudent); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	v := validate.Struct(updatedStudent)
	v.StopOnError = false
	if !v.Validate() {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "validation failed",
			"errors":  v.Errors.All(),
		})
	}

	student, err := ctrl.studentUsecase.UpdateStudent(uint(id), &updatedStudent)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Email already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, student)
}

func (ctrl *StudentController) DeleteStudent(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid student ID"})
	}

	err = ctrl.studentUsecase.DeleteStudent(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (ctrl *StudentController) GetAverageAge(c *echo.Context) error {
	averageAge, err := ctrl.studentUsecase.AvarageAge()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if averageAge == -1 {
		return c.JSON(http.StatusNoContent, map[string]string{"error": "No students found"})
	}
	return c.JSON(http.StatusOK, map[string]float64{"average_age": averageAge})
}
