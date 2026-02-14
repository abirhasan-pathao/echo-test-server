package repository

import (
	"echo-server/app/students/model"

	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {

	return &StudentRepository{db: db}
}

func (r *StudentRepository) CreateStudent(student *model.Student) error {

	return r.db.Create(student).Error
}

func (r *StudentRepository) GetStudentByID(id uint) (*model.Student, error) {
	var student model.Student
	err := r.db.First(&student, id).Error

	return &student, err
}

func (r *StudentRepository) GetAllStudents() ([]model.Student, error) {
	var students []model.Student
	err := r.db.Find(&students).Error

	return students, err
}

func (r *StudentRepository) UpdateStudent(id uint, updatedStudent *model.Student) (*model.Student, error) {
	var student model.Student
	if err := r.db.First(&student, id).Error; err != nil {
		return nil, err
	}

	student.Name = updatedStudent.Name
	student.Age = updatedStudent.Age
	student.Email = updatedStudent.Email

	if err := r.db.Save(&student).Error; err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *StudentRepository) DeleteStudent(id uint) error {

	return r.db.Delete(&model.Student{}, id).Error
}

func (r *StudentRepository) AvarageAge() (float64, error) {
	var avgAge float64
	err := r.db.Model(&model.Student{}).Select("AVG(age)").Scan(&avgAge).Error

	return avgAge, err
}
