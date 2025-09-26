package service

import (
	"net/http/httptest"
	"testing"

	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/test/mocks/container"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	simple1 = model.SimpleForm{Name: "Simple 1"}
)

func TestCreateSimple_Success(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	testContainer := container.GetTestContainer()
	// and
	target := testContainer.SimpleService
	simpleRepository := testContainer.SimpleRepository
	simples, err := simpleRepository.GetAll(ctx)
	assert.Nil(t, err)
	assert.Empty(t, simples)
	// when
	simple, err := target.CreateSimple(ctx, simple1)
	// then
	assert.Nil(t, err)
	assert.Equal(t, uint(1), simple.ID)
	assert.Equal(t, simple1.Name, simple.Name)
	// and
	simples, err = simpleRepository.GetAll(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(simples))
	assert.Equal(t, simple.ID, simples[0].ID)
	assert.Equal(t, simple.Name, simples[0].Name)
}
