package model

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name  string `gorm:"not null" json:"name" validate:"required,min=3"`
	Age   int    `gorm:"not null" json:"age" validate:"required,gt=0"`
	Email string `gorm:"unique;not null" json:"email" validate:"required,email"`
}
