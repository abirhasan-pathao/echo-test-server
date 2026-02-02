package utils

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.Validator.Struct(i); err != nil {
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
