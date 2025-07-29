package model

import (
	"time"

	"go.uber.org/zap/zapcore"
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

func (user *User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddUint("id", user.ID)
	enc.AddString("email", user.Email)
	enc.AddTime("created_at", user.CreatedAt)
	enc.AddTime("updated_at", user.UpdatedAt)
	if user.DeletedAt.Valid {
		enc.AddTime("deleted_at", user.DeletedAt.Time)
	} else {
		enc.AddString("deleted_at", "null")
	}
	return nil
}

func (userForm *UserForm) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("email", userForm.Email)
	if userForm.Password != "" {
		enc.AddString("password", "[PROVIDED]")
		enc.AddInt("password_length", len(userForm.Password))
	} else {
		enc.AddString("password", "[NOT PROVIDED]")
	}
	return nil
}
