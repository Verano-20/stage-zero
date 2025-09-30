package service

import (
	"errors"
	"testing"

	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/service"
	"github.com/Verano-20/go-crud/test/mocks/repository"
	"github.com/Verano-20/go-crud/test/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	simple1 = model.Simple{ID: 1, Name: "Simple 1"}
	simple2 = model.Simple{ID: 2, Name: "Simple 2"}
)

func createSimpleServiceAndMockRepo() (service.SimpleService, *repository.MockSimpleRepository) {
	mockRepo := repository.NewMockSimpleRepository()
	target := service.NewSimpleService(mockRepo)
	return target, mockRepo
}

/*
 * Create Simple Tests
 */

func TestCreateSimple_Success(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceAndMockRepo()
	// expect
	simpleRepository.On("Create", ctx, mock.MatchedBy(func(simple *model.Simple) bool {
		return simple.Name == simple1.Name
	})).Return(&simple1, nil)
	// when
	result, err := target.CreateSimple(ctx, *simple1.ToForm())
	// then
	assert.NoError(t, err)
	assert.Equal(t, &simple1, result)
	simpleRepository.AssertExpectations(t)
}

func TestCreateSimple_Error(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceAndMockRepo()
	expectedError := errors.New("database error")
	// expect
	simpleRepository.On("Create", ctx, mock.MatchedBy(func(simple *model.Simple) bool {
		return simple.Name == simple1.Name
	})).Return(nil, expectedError)
	// when
	result, err := target.CreateSimple(ctx, *simple1.ToForm())
	// then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	simpleRepository.AssertExpectations(t)
}

/*
 * Get All Simples Tests
 */

func TestGetAllSimples_Success(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceAndMockRepo()
	expectedSimples := model.Simples{&simple1, &simple2}
	// expect
	simpleRepository.On("GetAll", ctx).Return(expectedSimples, nil)
	// when
	result, err := target.GetAllSimples(ctx)
	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedSimples, result)
	simpleRepository.AssertExpectations(t)
}

func TestGetAllSimples_Error(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceAndMockRepo()
	expectedError := errors.New("database error")
	// expect
	simpleRepository.On("GetAll", ctx).Return(nil, expectedError)
	// when
	result, err := target.GetAllSimples(ctx)
	// then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	simpleRepository.AssertExpectations(t)
}

/*
 * Get Simple By ID Tests
 */

func TestGetSimpleByID_Success(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceAndMockRepo()
	// expect
	simpleRepository.On("GetByID", ctx, simple1.ID).Return(&simple1, nil)
	// when
	result, err := target.GetSimpleByID(ctx, uint64(simple1.ID))
	// then
	assert.NoError(t, err)
	assert.Equal(t, &simple1, result)
	simpleRepository.AssertExpectations(t)
}

func TestGetSimpleByID_Error(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceAndMockRepo()
	expectedError := errors.New("database error")
	// expect
	simpleRepository.On("GetByID", ctx, simple1.ID).Return(nil, expectedError)
	// when
	result, err := target.GetSimpleByID(ctx, uint64(simple1.ID))
	// then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	simpleRepository.AssertExpectations(t)
}

/*
 * Update Simple Tests
 */

func TestUpdateSimple_Success(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceAndMockRepo()
	// expect
	simpleRepository.On("Update", ctx, &simple1).Return(&simple1, nil)
	// when
	result, err := target.UpdateSimple(ctx, &simple1, *simple1.ToForm())
	// then
	assert.NoError(t, err)
	assert.Equal(t, &simple1, result)
	simpleRepository.AssertExpectations(t)
}

func TestUpdateSimple_Error(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceAndMockRepo()
	expectedError := errors.New("database error")
	// expect
	simpleRepository.On("Update", ctx, &simple1).Return(nil, expectedError)
	// when
	result, err := target.UpdateSimple(ctx, &simple1, *simple1.ToForm())
	// then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	simpleRepository.AssertExpectations(t)
}

/*
 * Delete Simple Tests
 */

func TestDeleteSimple_Success(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceAndMockRepo()
	// expect
	simpleRepository.On("Delete", ctx, simple1.ID).Return(nil)
	// when
	err := target.DeleteSimple(ctx, &simple1)
	// then
	assert.NoError(t, err)
	simpleRepository.AssertExpectations(t)
}

func TestDeleteSimple_Error(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceAndMockRepo()
	expectedError := errors.New("database error")
	// expect
	simpleRepository.On("Delete", ctx, simple1.ID).Return(expectedError)
	// when
	err := target.DeleteSimple(ctx, &simple1)
	// then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	simpleRepository.AssertExpectations(t)
}
