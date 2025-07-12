package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type UserForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (userForm *UserForm) ToModel() *User {
	return &User{
		Email:    userForm.Email,
		Password: userForm.Password,
	}
}
