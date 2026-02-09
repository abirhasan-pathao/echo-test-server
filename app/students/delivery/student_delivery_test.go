package delivery

import (
	"echo-server/app/students/model"
	"echo-server/app/utils/validation"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
)

type MockStudentUsecase struct {
	CreateStudentFunc  func(student *model.Student) error
	GetStudentByIDFunc func(id uint) (*model.Student, error)
	GetAllStudentsFunc func() ([]model.Student, error)
	UpdateStudentFunc  func(id uint, updatedStudent *model.Student) (*model.Student, error)
	DeleteStudentFunc  func(id uint) error
	AvarageAgeFunc     func() (float64, error)
}

func (m *MockStudentUsecase) CreateStudent(student *model.Student) error {
	return m.CreateStudentFunc(student)
}

func (m *MockStudentUsecase) GetStudentByID(id uint) (*model.Student, error) {
	return m.GetStudentByIDFunc(id)
}

func (m *MockStudentUsecase) GetAllStudents() ([]model.Student, error) {
	return m.GetAllStudentsFunc()
}

func (m *MockStudentUsecase) UpdateStudent(id uint, updatedStudent *model.Student) (*model.Student, error) {
	return m.UpdateStudentFunc(id, updatedStudent)
}

func (m *MockStudentUsecase) DeleteStudent(id uint) error {
	return m.DeleteStudentFunc(id)
}

func (m *MockStudentUsecase) AvarageAge() (float64, error) {
	return m.AvarageAgeFunc()
}

func Test_CreateStudent(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		method      string
		mockStudent model.Student
		mockError   error
		wantBody    string
		wantCode    int
	}{
		{
			name:   "Create student successfully",
			url:    "/students",
			method: "POST",
			mockStudent: model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john.doe@example.com",
			},
			wantBody: "",
			wantCode: 201,
		},
		{
			name:   "Validation error - missing name",
			url:    "/students",
			method: "POST",
			mockStudent: model.Student{
				Age:   20,
				Email: "john.doe@example.com",
			},
			wantBody: "",
			wantCode: 400,
		},
		{
			name:   "Validation error - invalid email",
			url:    "/students",
			method: "POST",
			mockStudent: model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "invalid-email",
			},
			wantBody: "",
			wantCode: 400,
		},
		{
			name:   "Internal server error",
			url:    "/students",
			method: "POST",
			mockStudent: model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john.doe@example.com",
			},
			mockError: fmt.Errorf("database error"),
			wantBody:  "",
			wantCode:  500,
		},
		{
			name:   "Conflict error - duplicate email",
			url:    "/students",
			method: "POST",
			mockStudent: model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john.doe@example.com",
			},
			mockError: fmt.Errorf("duplicate key value"),
			wantBody:  "",
			wantCode:  409,
		},
	}

	for _, tc := range tests {
		mockUC := &MockStudentUsecase{
			CreateStudentFunc: func(student *model.Student) error {
				if tc.mockError != nil {
					return tc.mockError
				}
				return nil
			},
		}
		ctrl := NewStudentController(mockUC)

		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.mockStudent)
			req := httptest.NewRequest(tc.method, tc.url, strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e := echo.New()
			e.POST("/students", ctrl.CreateStudent)
			validation.AddValidator()
			e.ServeHTTP(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("CreateStudent() gotCode = %v, want %v", rec.Code, tc.wantCode)
			}

		})
	}
}

func Test_GetStudentByID(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		id          string
		method      string
		mockStudent *model.Student
		mockError   error
		wantBody    string
		wantCode    int
	}{
		{
			name:   "Get student by ID successfully",
			url:    "/students/1",
			id:     "1",
			method: "GET",
			mockStudent: &model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john.doe@example.com",
			},
			wantBody: `{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"name":"John Doe","age":20,"email":"john.doe@example.com"}`,
			wantCode: 200,
		},
		{
			name:      "Student not found",
			url:       "/students/999",
			id:        "999",
			method:    "GET",
			mockError: fmt.Errorf("record not found"),
			wantBody:  `{"error":"Student not found"}`,
			wantCode:  404,
		},
		{
			name:     "Invalid student ID",
			url:      "/students/abc",
			id:       "abc",
			method:   "GET",
			wantBody: `{"error":"Invalid student ID"}`,
			wantCode: 400,
		},
		{
			name:      "Internal server error",
			url:       "/students/1",
			id:        "1",
			method:    "GET",
			mockError: fmt.Errorf("database error"),
			wantBody:  `{"error":"database error"}`,
			wantCode:  500,
		},
		{
			name:     "Invalid student ID - negative",
			url:      "/students/-1",
			id:       "-1",
			method:   "GET",
			wantBody: `{"error":"Invalid student ID"}`,
			wantCode: 400,
		},
		{
			name:     "Invalid student ID - zero",
			url:      "/students/0",
			id:       "0",
			method:   "GET",
			wantBody: `{"error":"Invalid student ID"}`,
			wantCode: 400,
		},
	}

	for _, tc := range tests {
		mockUC := &MockStudentUsecase{
			GetStudentByIDFunc: func(id uint) (*model.Student, error) {
				if tc.mockError != nil {
					return nil, tc.mockError
				}
				return tc.mockStudent, nil
			},
		}
		ctrl := NewStudentController(mockUC)

		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, nil)
			rec := httptest.NewRecorder()
			e := echo.New()
			e.GET("/students/:id", ctrl.GetStudentByID)
			e.ServeHTTP(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("GetStudentByID() gotCode = %v, want %v", rec.Code, tc.wantCode)
			}

			if strings.TrimSpace(rec.Body.String()) != tc.wantBody {
				t.Errorf("GetStudentByID() gotBody = %v, want %v", rec.Body.String(), tc.wantBody)
			}
		})
	}
}

func Test_GetAllStudents(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		method       string
		mockStudents []model.Student
		mockError    error
		wantBody     string
		wantCode     int
	}{
		{
			name:   "Get all students successfully",
			url:    "/students",
			method: "GET",
			mockStudents: []model.Student{
				{Name: "John Doe", Age: 20, Email: "john@example.com"},
				{Name: "Jane Doe", Age: 22, Email: "jane@example.com"},
			},
			wantBody: `[{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"name":"John Doe","age":20,"email":"john@example.com"},{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"name":"Jane Doe","age":22,"email":"jane@example.com"}]`,
			wantCode: 200,
		},
		{
			name:         "Get all students - empty list",
			url:          "/students",
			method:       "GET",
			mockStudents: []model.Student{},
			wantBody:     `[]`,
			wantCode:     200,
		},
		{
			name:      "Internal server error",
			url:       "/students",
			method:    "GET",
			mockError: fmt.Errorf("database error"),
			wantBody:  `{"error":"database error"}`,
			wantCode:  500,
		},
	}

	for _, tc := range tests {
		mockUC := &MockStudentUsecase{
			GetAllStudentsFunc: func() ([]model.Student, error) {
				if tc.mockError != nil {
					return nil, tc.mockError
				}
				return tc.mockStudents, nil
			},
		}
		ctrl := NewStudentController(mockUC)

		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, nil)
			rec := httptest.NewRecorder()
			e := echo.New()
			e.GET("/students", ctrl.GetAllStudents)
			e.ServeHTTP(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("GetAllStudents() gotCode = %v, want %v", rec.Code, tc.wantCode)
			}

			if strings.TrimSpace(rec.Body.String()) != tc.wantBody {
				t.Errorf("GetAllStudents() gotBody = %v, want %v", rec.Body.String(), tc.wantBody)
			}
		})
	}
}

func Test_UpdateStudent(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		method       string
		inputStudent model.Student
		mockStudent  *model.Student
		mockError    error
		wantBody     string
		wantCode     int
	}{
		{
			name:   "Update student successfully",
			url:    "/students/1",
			method: "PUT",
			inputStudent: model.Student{
				Name:  "John Updated",
				Age:   25,
				Email: "john.updated@example.com",
			},
			mockStudent: &model.Student{
				Name:  "John Updated",
				Age:   25,
				Email: "john.updated@example.com",
			},
			wantBody: `{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"name":"John Updated","age":25,"email":"john.updated@example.com"}`,
			wantCode: 200,
		},
		{
			name:   "Invalid student ID - non-numeric",
			url:    "/students/abc",
			method: "PUT",
			inputStudent: model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john@example.com",
			},
			wantBody: `{"error":"Invalid student ID"}`,
			wantCode: 400,
		},
		{
			name:   "Invalid student ID - zero",
			url:    "/students/0",
			method: "PUT",
			inputStudent: model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john@example.com",
			},
			wantBody: `{"error":"Invalid student ID"}`,
			wantCode: 400,
		},
		{
			name:   "Invalid student ID - negative",
			url:    "/students/-1",
			method: "PUT",
			inputStudent: model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john@example.com",
			},
			wantBody: `{"error":"Invalid student ID"}`,
			wantCode: 400,
		},
		{
			name:   "Validation error - missing name",
			url:    "/students/1",
			method: "PUT",
			inputStudent: model.Student{
				Age:   20,
				Email: "john@example.com",
			},
			wantCode: 400,
		},
		{
			name:   "Validation error - invalid email",
			url:    "/students/1",
			method: "PUT",
			inputStudent: model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "invalid-email",
			},
			wantCode: 400,
		},
		{
			name:   "Duplicate email error",
			url:    "/students/1",
			method: "PUT",
			inputStudent: model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john@example.com",
			},
			mockError: fmt.Errorf("duplicate key value"),
			wantBody:  `{"error":"Email already exists"}`,
			wantCode:  409,
		},
		{
			name:   "Internal server error",
			url:    "/students/1",
			method: "PUT",
			inputStudent: model.Student{
				Name:  "John Doe",
				Age:   20,
				Email: "john@example.com",
			},
			mockError: fmt.Errorf("database error"),
			wantBody:  `{"error":"database error"}`,
			wantCode:  500,
		},
	}

	for _, tc := range tests {
		mockUC := &MockStudentUsecase{
			UpdateStudentFunc: func(id uint, student *model.Student) (*model.Student, error) {
				if tc.mockError != nil {
					return nil, tc.mockError
				}
				return tc.mockStudent, nil
			},
		}
		ctrl := NewStudentController(mockUC)

		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.inputStudent)
			req := httptest.NewRequest(tc.method, tc.url, strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e := echo.New()
			e.PUT("/students/:id", ctrl.UpdateStudent)
			validation.AddValidator()
			e.ServeHTTP(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("UpdateStudent() gotCode = %v, want %v", rec.Code, tc.wantCode)
			}

			if tc.wantBody != "" && strings.TrimSpace(rec.Body.String()) != tc.wantBody {
				t.Errorf("UpdateStudent() gotBody = %v, want %v", rec.Body.String(), tc.wantBody)
			}
		})
	}
}

func Test_DeleteStudent(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		method    string
		mockError error
		wantBody  string
		wantCode  int
	}{
		{
			name:     "Delete student successfully",
			url:      "/students/1",
			method:   "DELETE",
			wantBody: "",
			wantCode: 204,
		},
		{
			name:     "Invalid student ID - non-numeric",
			url:      "/students/abc",
			method:   "DELETE",
			wantBody: `{"error":"Invalid student ID"}`,
			wantCode: 400,
		},
		{
			name:     "Invalid student ID - zero",
			url:      "/students/0",
			method:   "DELETE",
			wantBody: `{"error":"Invalid student ID"}`,
			wantCode: 400,
		},
		{
			name:     "Invalid student ID - negative",
			url:      "/students/-1",
			method:   "DELETE",
			wantBody: `{"error":"Invalid student ID"}`,
			wantCode: 400,
		},
		{
			name:      "Internal server error",
			url:       "/students/1",
			method:    "DELETE",
			mockError: fmt.Errorf("database error"),
			wantBody:  `{"error":"database error"}`,
			wantCode:  500,
		},
	}

	for _, tc := range tests {
		mockUC := &MockStudentUsecase{
			DeleteStudentFunc: func(id uint) error {
				return tc.mockError
			},
		}
		ctrl := NewStudentController(mockUC)

		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, nil)
			rec := httptest.NewRecorder()
			e := echo.New()
			e.DELETE("/students/:id", ctrl.DeleteStudent)
			e.ServeHTTP(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("DeleteStudent() gotCode = %v, want %v", rec.Code, tc.wantCode)
			}

			if tc.wantBody != "" && strings.TrimSpace(rec.Body.String()) != tc.wantBody {
				t.Errorf("DeleteStudent() gotBody = %v, want %v", rec.Body.String(), tc.wantBody)
			}
		})
	}
}

func Test_GetAverageAge(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		method    string
		mockAge   float64
		mockError error
		wantBody  string
		wantCode  int
	}{
		{
			name:     "Get average age successfully",
			url:      "/students/average-age",
			method:   "GET",
			mockAge:  21.5,
			wantBody: `{"average_age":21.5}`,
			wantCode: 200,
		},
		{
			name:     "No students found",
			url:      "/students/average-age",
			method:   "GET",
			mockAge:  -1,
			wantBody: `{"error":"No students found"}`,
			wantCode: 204,
		},
		{
			name:      "Internal server error",
			url:       "/students/average-age",
			method:    "GET",
			mockError: fmt.Errorf("database error"),
			wantBody:  `{"error":"database error"}`,
			wantCode:  500,
		},
	}

	for _, tc := range tests {
		mockUC := &MockStudentUsecase{
			AvarageAgeFunc: func() (float64, error) {
				if tc.mockError != nil {
					return 0, tc.mockError
				}
				return tc.mockAge, nil
			},
		}
		ctrl := NewStudentController(mockUC)

		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, nil)
			rec := httptest.NewRecorder()
			e := echo.New()
			e.GET("/students/average-age", ctrl.GetAverageAge)
			e.ServeHTTP(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("GetAverageAge() gotCode = %v, want %v", rec.Code, tc.wantCode)
			}

			if tc.wantBody != "" && strings.TrimSpace(rec.Body.String()) != tc.wantBody {
				t.Errorf("GetAverageAge() gotBody = %v, want %v", rec.Body.String(), tc.wantBody)
			}
		})
	}
}
