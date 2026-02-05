package model

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name  string `gorm:"not null" json:"name" validate:"required|minLen:2|maxLen:100" message:"required:{field} is required|minLen:{field} must have at least 2 chars|maxLen:{field} can have at most 100 chars"`
	Age   int    `gorm:"not null" json:"age" validate:"required|int|min:1|max:120" message:"required:{field} is required|min:{field} must be at least 1|max:{field} must be at most 120"`
	Email string `gorm:"unique;not null" json:"email" validate:"required|email" message:"required:{field} is required|email:{field} must be a valid email"`
}
