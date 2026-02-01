package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Name  string `json:"name" form:"name"`
	Age   int    `json:"age" form:"age"`
	Email string `json:"email" form:"email" gorm:"unique"`
}
