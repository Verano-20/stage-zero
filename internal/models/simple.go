package models

import (
	"gorm.io/gorm"
)

type Simple struct {
	gorm.Model
	Name string `json:"name" binding:"required"`
} 