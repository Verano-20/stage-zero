package service

import (
	"github.com/Verano-20/stage-zero/internal/model"
	"github.com/Verano-20/stage-zero/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

var _ service.UserService = &MockUserService{}

func NewMockUserService() *MockUserService {
	return &MockUserService{}
}

func (m *MockUserService) CreateUser(ctx *gin.Context, userForm model.UserForm) (user *model.User, createErr error) {
	args := m.Called(ctx, userForm)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx *gin.Context, email string) (user *model.User, err error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}
