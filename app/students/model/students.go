package model

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name  string `json:"name" form:"name" validate:"required,min=2,max=100,alphaspace"`
	Age   int    `json:"age" form:"age" validate:"required,min=1,max=120"`
	Email string `json:"email" form:"email" gorm:"unique" validate:"required,email"`
}
