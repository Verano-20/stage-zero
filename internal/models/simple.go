package models

import (
	"time"
)

type Simple struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt time.Time      `json:"deleted_at"`
}

type SimpleForm struct {
	Name string `json:"name"`
}

// Explicitly tell GORM this is a read-only model
func (Simple) TableName() string {
	return "simples"
}
