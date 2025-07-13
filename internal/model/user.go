package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id"`
	Email        string         `json:"email"`
	PasswordHash string         `json:"password_hash"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

type UserDTO struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Users []*User

func (user *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (userForm *UserForm) ToModel(hashedPassword string) *User {
	return &User{
		Email:        userForm.Email,
		PasswordHash: hashedPassword,
	}
}
