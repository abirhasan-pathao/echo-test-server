package usecase

import (
	"echo-server/app/students/model"
	"fmt"
	"testing"
)

type MockStudentRepository struct {
	students map[uint]*model.Student
	nextID   uint
	avg      float64
	err      error
}

func NewMockStudentRepository() *MockStudentRepository {
	return &MockStudentRepository{
		students: make(map[uint]*model.Student),
		nextID:   1,
	}
}

func (m *MockStudentRepository) CreateStudent(student *model.Student) error {
	if m.err != nil {
		return m.err
	}

	student.ID = m.nextID
	m.nextID++
	m.students[student.ID] = student
	return nil
}

func (m *MockStudentRepository) GetStudentByID(id uint) (*model.Student, error) {
	if m.err != nil {
		return nil, m.err
	}

	student, exists := m.students[id]
	if !exists {
		return nil, fmt.Errorf("student not found")
	}
	return student, nil
}

func (m *MockStudentRepository) GetAllStudents() ([]model.Student, error) {
	if m.err != nil {
		return nil, m.err
	}

	var students []model.Student
	for _, student := range m.students {
		students = append(students, *student)
	}
	return students, nil
}

func (m *MockStudentRepository) UpdateStudent(id uint, updatedStudent *model.Student) (*model.Student, error) {
	if m.err != nil {
		return nil, m.err
	}

	student, exists := m.students[id]
	if !exists {
		return nil, fmt.Errorf("student not found")
	}

	student.Name = updatedStudent.Name
	student.Age = updatedStudent.Age
	student.Email = updatedStudent.Email

	return student, nil
}

func (m *MockStudentRepository) DeleteStudent(id uint) error {
	if m.err != nil {
		return m.err
	}

	delete(m.students, id)
	return nil
}

func (m *MockStudentRepository) AvarageAge() (float64, error) {
	if m.err != nil {
		return 0, m.err
	}

	return m.avg, nil
}

func Test_CreateStudent(t *testing.T) {
	tests := []struct {
		name      string
		student   *model.Student
		repoErr   error
		expectErr bool
	}{
		{
			name: "Success",
			student: &model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john.doe@example.com",
			},
			repoErr:   nil,
			expectErr: false,
		},
		{
			name: "Repository error",
			student: &model.Student{
				Name:  "Jane Doe",
				Age:   22,
				Email: "jane.doe@example.com",
			},
			repoErr:   fmt.Errorf("repository error"),
			expectErr: true,
		},
	}

	for _, tc := range tests {
		mockRepo := NewMockStudentRepository()
		usecase := NewStudentUsecase(mockRepo)

		t.Run(tc.name, func(t *testing.T) {
			mockRepo.err = tc.repoErr

			err := usecase.CreateStudent(tc.student)
			if tc.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func Test_GetStudentByID(t *testing.T) {
	tests := []struct {
		name          string
		student       *model.Student
		setupRepo     func(*MockStudentRepository, *model.Student)
		expectStudent *model.Student
		expectErr     bool
	}{
		{
			name: "Success",
			student: &model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john.doe@example.com",
			},
			setupRepo: func(repo *MockStudentRepository, student *model.Student) {
				repo.CreateStudent(student)
			},
			expectErr: false,
		},
		{
			name: "Student not found",
			setupRepo: func(repo *MockStudentRepository, student *model.Student) {
				// No students added to the repository
			},
			expectStudent: nil,
			expectErr:     true,
		},
		{
			name: "Repository error",
			student: &model.Student{
				Name:  "Jane Doe",
				Age:   22,
				Email: "jane.doe@example.com",
			},
			setupRepo: func(repo *MockStudentRepository, student *model.Student) {
				repo.err = fmt.Errorf("repository error")
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		mockRepo := NewMockStudentRepository()
		usecase := NewStudentUsecase(mockRepo)

		t.Run(tc.name, func(t *testing.T) {
			tc.setupRepo(mockRepo, tc.student)

			student, err := usecase.GetStudentByID(1)
			if tc.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.expectErr && student != tc.student {
				t.Errorf("expected student %v but got %v", tc.student, student)
			}
		})
	}
}

func Test_GetAllStudents(t *testing.T) {
	tests := []struct {
		name           string
		students       []*model.Student
		setupRepo      func(*MockStudentRepository, []*model.Student)
		expectStudents []model.Student
		expectErr      bool
	}{
		{
			name: "Success",
			students: []*model.Student{
				{Name: "John Doe", Age: 20, Email: "john.doe@example.com"},
				{Name: "Jane Doe", Age: 22, Email: "jane.doe@example.com"},
				{Name: "Alice Smith", Age: 19, Email: "alice.smith@example.com"},
			},
			setupRepo: func(repo *MockStudentRepository, students []*model.Student) {
				for _, student := range students {
					repo.CreateStudent(student)
				}
			},
			expectErr: false,
		},
		{
			name: "No students found",
			setupRepo: func(repo *MockStudentRepository, students []*model.Student) {
				// No students added to the repository
			},
			expectStudents: []model.Student{},
			expectErr:      false,
		},
		{
			name:     "Repository error",
			students: []*model.Student{},
			setupRepo: func(repo *MockStudentRepository, students []*model.Student) {
				repo.err = fmt.Errorf("repository error")
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		mockRepo := NewMockStudentRepository()
		usecase := NewStudentUsecase(mockRepo)

		t.Run(tc.name, func(t *testing.T) {
			tc.setupRepo(mockRepo, tc.students)

			students, err := usecase.GetAllStudents()
			if tc.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.expectErr && len(students) != len(tc.students) {
				t.Errorf("expected %d students but got %d", len(tc.students), len(students))
			}
		})
	}
}

func Test_UpdateStudent(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		original  *model.Student
		updated   *model.Student
		setupRepo func(*MockStudentRepository, *model.Student)
		expectErr bool
	}{
		{
			name: "Success",
			id:   1,
			original: &model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john.doe@example.com",
			},
			updated: &model.Student{
				Name:  "John Smith",
				Age:   21,
				Email: "john.smith@example.com",
			},
			setupRepo: func(repo *MockStudentRepository, student *model.Student) {
				repo.CreateStudent(student)
			},
			expectErr: false,
		},
		{
			name: "Student not found",
			id:   10,
			setupRepo: func(repo *MockStudentRepository, student *model.Student) {
				//
			},
			expectErr: true,
		},
		{
			name: "Repository error",
			id:   1,
			original: &model.Student{
				Name:  "Jane Doe",
				Age:   22,
				Email: "jane.doe@example.com",
			},
			updated: &model.Student{
				Name:  "Jane Smith",
				Age:   23,
				Email: "jane.smith@example.com",
			},
			setupRepo: func(repo *MockStudentRepository, student *model.Student) {
				repo.CreateStudent(student)
				repo.err = fmt.Errorf("repository error")
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		mockRepo := NewMockStudentRepository()
		usecase := NewStudentUsecase(mockRepo)

		t.Run(tc.name, func(t *testing.T) {
			tc.setupRepo(mockRepo, tc.original)

			_, err := usecase.UpdateStudent(tc.id, tc.updated)
			u, errGet := usecase.GetStudentByID(tc.id)
			if tc.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.expectErr && errGet != nil {
				t.Errorf("unexpected error on get: %v", errGet)
			}
			if !tc.expectErr && (u.Name != tc.updated.Name || u.Age != tc.updated.Age || u.Email != tc.updated.Email) {
				t.Errorf("expected updated student %v but got %v", tc.updated, u)
			}
		})
	}
}

func Test_DeleteStudent(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		student   *model.Student
		setupRepo func(*MockStudentRepository, *model.Student)
		expectErr bool
	}{
		{
			name: "Success",
			id:   1,
			student: &model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john.doe@example.com",
			},
			setupRepo: func(repo *MockStudentRepository, student *model.Student) {
				repo.CreateStudent(student)
			},
			expectErr: false,
		},
		{
			name: "Deleting non-existent student",
			id:   10,
			setupRepo: func(repo *MockStudentRepository, student *model.Student) {
				//
			},
			expectErr: false,
		},
		{
			name: "Repository error",
			id:   1,
			student: &model.Student{
				Name:  "Jane Doe",
				Age:   22,
				Email: "jane.doe@example.com",
			},
			setupRepo: func(repo *MockStudentRepository, student *model.Student) {
				repo.CreateStudent(student)
				repo.err = fmt.Errorf("repository error")
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		mockRepo := NewMockStudentRepository()
		usecase := NewStudentUsecase(mockRepo)

		t.Run(tc.name, func(t *testing.T) {
			tc.setupRepo(mockRepo, tc.student)

			err := usecase.DeleteStudent(tc.id)
			if tc.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func Test_AvarageAge(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*MockStudentRepository)
		expectAvg float64
		expectErr bool
	}{
		{
			name: "Success",
			setupRepo: func(repo *MockStudentRepository) {
				repo.avg = 21.0
			},
			expectAvg: 21.0,
			expectErr: false,
		},
		{
			name: "Repository error",
			setupRepo: func(repo *MockStudentRepository) {
				repo.err = fmt.Errorf("repository error")
			},
			expectAvg: 0,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		mockRepo := NewMockStudentRepository()
		usecase := NewStudentUsecase(mockRepo)

		t.Run(tc.name, func(t *testing.T) {
			tc.setupRepo(mockRepo)

			avg, err := usecase.AvarageAge()
			if tc.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.expectErr && avg != tc.expectAvg {
				t.Errorf("expected average age %v but got %v", tc.expectAvg, avg)
			}
		})
	}
}
