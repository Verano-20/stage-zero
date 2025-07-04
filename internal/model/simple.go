package model

import (
	"time"
)

type Simple struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type SimpleDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SimpleForm struct {
	Name string `json:"name"`
}

type Simples []*Simple

func (simple *Simple) ToDTO() *SimpleDTO {
	return &SimpleDTO{
		ID:        simple.ID,
		Name:      simple.Name,
		CreatedAt: simple.CreatedAt,
		UpdatedAt: simple.UpdatedAt,
	}
}

func (simpleForm *SimpleForm) ToModel() *Simple {
	return &Simple{
		Name: simpleForm.Name,
	}
}

func (simples Simples) ToDTOs() []*SimpleDTO {
	simpleDTOs := make([]*SimpleDTO, len(simples))
	for i, simple := range simples {
		simpleDTOs[i] = simple.ToDTO()
	}
	return simpleDTOs
}

// Explicitly tell GORM this is a read-only model
func (Simple) TableName() string {
	return "simples"
}
