package handlers

import (
	"echo-server/db"
	"echo-server/models"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

func CreateStudent(c *echo.Context) error {
	var student models.Student

	contentTYpe := c.Request().Header.Get("Content-Type")
	log.Println("Content-Type:", contentTYpe)

	if err := c.Bind(&student); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := c.Validate(&student); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var existingStudent models.Student
	if err := db.DB.Where("email = ?", student.Email).First(&existingStudent).Error; err == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Student with this emailalready exists"})
	}

	err := db.DB.Create(&student).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, student)
}

func CreateStudentFormData(c *echo.Context) error {
	var student models.Student

	student.Name = c.FormValue("name")
	age, err := strconv.Atoi(c.FormValue("age"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid age"})
	}
	student.Age = age
	student.Email = c.FormValue("email")

	var existingStudent models.Student
	if err := db.DB.Where("email = ?", student.Email).First(&existingStudent).Error; err == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Student with this email already exists"})
	}

	err = db.DB.Create(&student).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, student)
}

func GetAllStudents(c *echo.Context) error {
	var students []models.Student
	query := db.DB
	if name := c.QueryParam("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if email := c.QueryParam("email"); email != "" {
		query = query.Where("email = ?", email)
	}
	query.Find(&students)

	return c.JSON(http.StatusOK, students)
}

func GetStudentByID(c *echo.Context) error {
	// log.Println("GetStudentByID called")
	id := c.Param("id")
	var student models.Student

	if err := db.DB.First(&student, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Student not found"})
	}

	return c.JSON(http.StatusOK, student)
}

func UpdateStudent(c *echo.Context) error {
	id := c.Param("id")
	var student models.Student

	if err := db.DB.First(&student, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Student not found"})
	}

	var updatedData models.Student
	if err := c.Bind(&updatedData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	student.Name = updatedData.Name
	student.Email = updatedData.Email
	student.Age = updatedData.Age

	if err := db.DB.Save(&student).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, student)
}

func DeleteStudent(c *echo.Context) error {
	id := c.Param("id")
	var student models.Student

	if err := db.DB.First(&student, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Student not found"})
	}

	if err := db.DB.Unscoped().Delete(&student).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Student deleted successfully"})
}
