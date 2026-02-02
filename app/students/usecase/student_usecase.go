package usecase

import "echo-server/app/students/model"

type StudentRepository interface {
	CreateStudent(student *model.Student) error
	GetStudentByID(id uint) (*model.Student, error)
	GetAllStudents() ([]model.Student, error)
	UpdateStudent(id uint, updatedStudent *model.Student) (*model.Student, error)
	DeleteStudent(id uint) error
}

type StudentUsecase struct {
	studentRepo StudentRepository
}

func NewStudentUsecase(studentRepo StudentRepository) *StudentUsecase {
	return &StudentUsecase{studentRepo: studentRepo}
}

func (u *StudentUsecase) CreateStudent(student *model.Student) error {
	return u.studentRepo.CreateStudent(student)
}

func (u *StudentUsecase) GetStudentByID(id uint) (*model.Student, error) {
	return u.studentRepo.GetStudentByID(id)
}

func (u *StudentUsecase) GetAllStudents() ([]model.Student, error) {
	return u.studentRepo.GetAllStudents()
}

func (u *StudentUsecase) UpdateStudent(id uint, updatedStudent *model.Student) (*model.Student, error) {
	return u.studentRepo.UpdateStudent(id, updatedStudent)
}

func (u *StudentUsecase) DeleteStudent(id uint) error {
	return u.studentRepo.DeleteStudent(id)
}
