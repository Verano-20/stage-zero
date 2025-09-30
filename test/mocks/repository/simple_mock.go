package repository

import (
	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockSimpleRepository struct {
	mock.Mock
}

var _ repository.SimpleRepository = &MockSimpleRepository{}

func NewMockSimpleRepository() *MockSimpleRepository {
	return &MockSimpleRepository{}
}

func (m *MockSimpleRepository) Create(ctx *gin.Context, simple *model.Simple) (*model.Simple, error) {
	args := m.Called(ctx, simple)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Simple), args.Error(1)
}

func (m *MockSimpleRepository) GetAll(ctx *gin.Context) (model.Simples, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(model.Simples), args.Error(1)
}

func (m *MockSimpleRepository) GetByID(ctx *gin.Context, id uint) (*model.Simple, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Simple), args.Error(1)
}

func (m *MockSimpleRepository) Update(ctx *gin.Context, simple *model.Simple) (*model.Simple, error) {
	args := m.Called(ctx, simple)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Simple), args.Error(1)
}

func (m *MockSimpleRepository) Delete(ctx *gin.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
