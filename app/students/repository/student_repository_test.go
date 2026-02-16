package repository

import (
	"echo-server/app/students/model"
	"fmt"
	"math"
	"os"
	"testing"
)

var dsn string

func init() {
	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_SSLMODE"),
	)
}

func Test_CreateStudent(t *testing.T) {
	tests := []struct {
		name    string
		student *model.Student
		setup   func(repo *StudentRepository)
		expErr  bool
	}{
		{
			name: "Create a new student",
			student: &model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john.oe@example.com",
			},
			expErr: false,
		},
		{
			name: "duplicate email",
			student: &model.Student{
				Name:  "Jane Doe",
				Age:   22,
				Email: "john.doe@example.com",
			},
			setup: func(repo *StudentRepository) {
				repo.CreateStudent(&model.Student{
					Name:  "John Doe",
					Age:   20,
					Email: "john.doe@example.com",
				})
			},
			expErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := SetupTestDB(t, dsn)
			tx := db.Begin()
			defer tx.Rollback()

			repo := NewStudentRepository(tx)
			if tt.setup != nil {
				tt.setup(repo)
			}

			err := repo.CreateStudent(tt.student)
			if err != nil {
				if !tt.expErr {
					t.Errorf("CreateStudent() error = %v", err)
				}
				return
			}
		})
	}
}

func Test_GetStudentByID(t *testing.T) {
	tests := []struct {
		name   string
		id     uint
		setup  func(repo *StudentRepository) uint
		expErr bool
	}{
		{
			name: "Get existing student by ID",
			setup: func(repo *StudentRepository) uint {
				student := &model.Student{
					Name:  "John Doe",
					Age:   20,
					Email: "john.doe@example.com",
				}
				repo.CreateStudent(student)
				return student.ID
			},
			expErr: false,
		},
		{
			name: "Get non-existing student by ID",
			setup: func(repo *StudentRepository) uint {
				return 999
			},
			expErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := SetupTestDB(t, dsn)
			tx := db.Begin()
			defer tx.Rollback()

			repo := NewStudentRepository(tx)
			if tt.setup != nil {
				tt.id = tt.setup(repo)
			}

			student, err := repo.GetStudentByID(tt.id)
			if err != nil {
				if !tt.expErr {
					t.Errorf("GetStudentByID() error = %v", err)
				}
				return
			}
			if student == nil {
				t.Error("GetStudentByID() returned nil student")
			}
		})
	}
}

func Test_GetAllStudents(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(repo *StudentRepository)
		expCount      int
		expErr        bool
		assertResults func(t *testing.T, students []model.Student)
	}{
		{
			name:     "Get all students when none exist",
			expCount: 0,
			expErr:   false,
		},
		{
			name: "Get all students with records",
			setup: func(repo *StudentRepository) {
				repo.CreateStudent(&model.Student{Name: "John Doe", Age: 20, Email: "john.doe@example.com"})
				repo.CreateStudent(&model.Student{Name: "Jane Doe", Age: 22, Email: "jane.doe@example.com"})
			},
			expCount: 2,
			expErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := SetupTestDB(t, dsn)
			tx := db.Begin()
			defer tx.Rollback()

			repo := NewStudentRepository(tx)
			if tt.setup != nil {
				tt.setup(repo)
			}

			students, err := repo.GetAllStudents()
			if err != nil {
				if !tt.expErr {
					t.Errorf("GetAllStudents() error = %v", err)
				}
				return
			}

			if len(students) != tt.expCount {
				t.Errorf("GetAllStudents() count = %d, expected %d", len(students), tt.expCount)
			}
		})
	}
}

func Test_UpdateStudent(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(repo *StudentRepository) uint
		id             uint
		updatedStudent *model.Student
		expErr         bool
		assertFn       func(t *testing.T, student *model.Student)
	}{
		{
			name: "Update existing student",
			setup: func(repo *StudentRepository) uint {
				student := &model.Student{Name: "John Doe", Age: 20, Email: "john.doe@example.com"}
				repo.CreateStudent(student)

				return student.ID
			},
			updatedStudent: &model.Student{Name: "John Updated", Age: 25, Email: "john.updated@example.com"},
			expErr:         false,
			assertFn: func(t *testing.T, student *model.Student) {
				if student.Name != "John Updated" {
					t.Errorf("UpdateStudent() name = %s, expected %s", student.Name, "John Updated")
				}
				if student.Age != 25 {
					t.Errorf("UpdateStudent() age = %d, expected %d", student.Age, 25)
				}
				if student.Email != "john.updated@example.com" {
					t.Errorf("UpdateStudent() email = %s, expected %s", student.Email, "john.updated@example.com")
				}
			},
		},
		{
			name: "Update non-existing student",
			setup: func(repo *StudentRepository) uint {
				return 999
			},
			updatedStudent: &model.Student{Name: "Ghost", Age: 30, Email: "ghost@example.com"},
			expErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := SetupTestDB(t, dsn)
			tx := db.Begin()
			defer tx.Rollback()

			repo := NewStudentRepository(tx)
			if tt.setup != nil {
				tt.id = tt.setup(repo)
			}

			student, err := repo.UpdateStudent(tt.id, tt.updatedStudent)
			if err != nil {
				if !tt.expErr {
					t.Errorf("UpdateStudent() error = %v", err)
				}
				return
			}

			if student == nil {
				t.Error("UpdateStudent() returned nil student")
				return
			}

			if tt.assertFn != nil {
				tt.assertFn(t, student)
			}
		})
	}
}

func Test_DeleteStudent(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(repo *StudentRepository) uint
		id     uint
		expErr bool
	}{
		{
			name: "Delete existing student",
			setup: func(repo *StudentRepository) uint {
				student := &model.Student{Name: "John Doe", Age: 20, Email: "john.doe@example.com"}
				repo.CreateStudent(student)

				return student.ID
			},
			expErr: false,
		},
		{
			name: "Delete non-existing student",
			setup: func(repo *StudentRepository) uint {
				return 999
			},
			expErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := SetupTestDB(t, dsn)
			tx := db.Begin()
			defer tx.Rollback()

			repo := NewStudentRepository(tx)
			if tt.setup != nil {
				tt.id = tt.setup(repo)
			}

			err := repo.DeleteStudent(tt.id)
			if err != nil {
				if !tt.expErr {
					t.Errorf("DeleteStudent() error = %v", err)
				}
				return
			}

			if !tt.expErr {
				_, err := repo.GetStudentByID(tt.id)
				if err == nil {
					t.Error("DeleteStudent() did not delete student")
				}
			}
		})
	}
}

func Test_AvarageAge(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(repo *StudentRepository)
		expAverage float64
		expErr     bool
		tolerance  float64
	}{
		{
			name:       "Average age with no students",
			expAverage: 0,
			expErr:     false,
			tolerance:  0.0001,
		},
		{
			name: "Average age with students",
			setup: func(repo *StudentRepository) {
				repo.CreateStudent(&model.Student{Name: "John Doe", Age: 20, Email: "john.doe@example.com"})
				repo.CreateStudent(&model.Student{Name: "Jane Doe", Age: 22, Email: "jane.doe@example.com"})
			},
			expAverage: 21,
			expErr:     false,
			tolerance:  0.0001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := SetupTestDB(t, dsn)
			tx := db.Begin()
			defer tx.Rollback()

			repo := NewStudentRepository(tx)
			if tt.setup != nil {
				tt.setup(repo)
			}

			avgAge, err := repo.AvarageAge()
			if err != nil {
				if !tt.expErr {
					t.Errorf("AvarageAge() error = %v", err)
				}
				return
			}

			if math.Abs(avgAge-tt.expAverage) > tt.tolerance {
				t.Errorf("AvarageAge() = %f, expected %f", avgAge, tt.expAverage)
			}
		})
	}
}
