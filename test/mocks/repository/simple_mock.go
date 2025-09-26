package repository

import (
	"errors"
	"maps"
	"slices"

	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/repository"
	"github.com/gin-gonic/gin"
)

type mockSimpleRepository struct {
	simples map[uint]*model.Simple
}

func NewMockSimpleRepository(simples map[uint]*model.Simple) repository.SimpleRepository {
	return &mockSimpleRepository{simples: simples}
}

var _ repository.SimpleRepository = &mockSimpleRepository{}

func (m *mockSimpleRepository) Create(ctx *gin.Context, simple *model.Simple) (*model.Simple, error) {
	simple.ID = uint(len(m.simples) + 1)
	m.simples[simple.ID] = simple
	return simple, nil
}

func (m *mockSimpleRepository) GetAll(ctx *gin.Context) (model.Simples, error) {
	return model.Simples(slices.Collect(maps.Values(m.simples))), nil
}

func (m *mockSimpleRepository) GetByID(ctx *gin.Context, id uint) (*model.Simple, error) {
	simple, ok := m.simples[id]
	if !ok {
		return nil, errors.New("simple not found")
	}
	return simple, nil
}

func (m *mockSimpleRepository) Update(ctx *gin.Context, simple *model.Simple) (*model.Simple, error) {
	_, ok := m.simples[simple.ID]
	if !ok {
		return nil, errors.New("simple not found")
	}
	m.simples[simple.ID] = simple
	return simple, nil
}

func (m *mockSimpleRepository) Delete(ctx *gin.Context, id uint) error {
	_, ok := m.simples[id]
	if !ok {
		return errors.New("simple not found")
	}
	delete(m.simples, id)
	return nil
}
