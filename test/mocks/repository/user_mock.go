package repository

import (
	"errors"

	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/repository"
	"github.com/gin-gonic/gin"
)

type mockUserRepository struct {
	users map[uint]*model.User
}

func NewMockUserRepository(users map[uint]*model.User) repository.UserRepository {
	return &mockUserRepository{users: users}
}

var _ repository.UserRepository = &mockUserRepository{}

func (m *mockUserRepository) Create(ctx *gin.Context, user *model.User) (*model.User, error) {
	m.users[user.ID] = user
	return user, nil
}

func (m *mockUserRepository) GetByID(ctx *gin.Context, id uint) (*model.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (m *mockUserRepository) GetByEmail(ctx *gin.Context, email string) (*model.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}
