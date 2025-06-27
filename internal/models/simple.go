package models

import (
	"time"

	"gorm.io/gorm"
)

type Simple struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Explicitly tell GORM this is a read-only model
func (Simple) TableName() string {
	return "simples"
}
