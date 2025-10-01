package model

import (
	"time"

	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type Simple struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type SimpleDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SimpleForm struct {
	Name string `json:"name" validate:"required,max=255"`
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

func (simple *Simple) ToForm() *SimpleForm {
	return &SimpleForm{
		Name: simple.Name,
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

func (s *Simple) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddUint("id", s.ID)
	enc.AddString("name", s.Name)
	enc.AddTime("created_at", s.CreatedAt)
	enc.AddTime("updated_at", s.UpdatedAt)
	return nil
}

func (s *SimpleForm) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", s.Name)
	return nil
}
