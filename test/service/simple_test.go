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

func createSimpleServiceWithMockDependencies(t *testing.T) (service.SimpleService, *repository.MockSimpleRepository) {
	mockRepo := repository.NewMockSimpleRepository()
	defer mockRepo.AssertExpectations(t)
	target := service.NewSimpleService(mockRepo)
	return target, mockRepo
}

/*
 * Create Simple Tests
 */

func TestCreateSimple_Success(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceWithMockDependencies(t)
	// expect
	simpleRepository.On("Create", ctx, mock.MatchedBy(func(simple *model.Simple) bool {
		return simple.Name == testutils.Simple1.Name
	})).Return(&testutils.Simple1, nil).Once()
	// when
	result, err := target.CreateSimple(ctx, *testutils.Simple1.ToForm())
	// then
	assert.NoError(t, err)
	assert.Equal(t, &testutils.Simple1, result)
	simpleRepository.AssertExpectations(t)
}

func TestCreateSimple_Error(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceWithMockDependencies(t)
	expectedError := errors.New("database error")
	// expect
	simpleRepository.On("Create", ctx, mock.MatchedBy(func(simple *model.Simple) bool {
		return simple.Name == testutils.Simple1.Name
	})).Return(nil, expectedError).Once()
	// when
	result, err := target.CreateSimple(ctx, *testutils.Simple1.ToForm())
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
	target, simpleRepository := createSimpleServiceWithMockDependencies(t)
	expectedSimples := model.Simples{&testutils.Simple1, &testutils.Simple2}
	// expect
	simpleRepository.On("GetAll", ctx).Return(expectedSimples, nil).Once()
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
	target, simpleRepository := createSimpleServiceWithMockDependencies(t)
	expectedError := errors.New("database error")
	// expect
	simpleRepository.On("GetAll", ctx).Return(nil, expectedError).Once()
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
	target, simpleRepository := createSimpleServiceWithMockDependencies(t)
	// expect
	simpleRepository.On("GetByID", ctx, testutils.Simple1.ID).Return(&testutils.Simple1, nil).Once()
	// when
	result, err := target.GetSimpleByID(ctx, uint64(testutils.Simple1.ID))
	// then
	assert.NoError(t, err)
	assert.Equal(t, &testutils.Simple1, result)
	simpleRepository.AssertExpectations(t)
}

func TestGetSimpleByID_Error(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceWithMockDependencies(t)
	expectedError := errors.New("database error")
	// expect
	simpleRepository.On("GetByID", ctx, testutils.Simple1.ID).Return(nil, expectedError).Once()
	// when
	result, err := target.GetSimpleByID(ctx, uint64(testutils.Simple1.ID))
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
	target, simpleRepository := createSimpleServiceWithMockDependencies(t)
	// expect
	simpleRepository.On("Update", ctx, &testutils.Simple1).Return(&testutils.Simple1, nil).Once()
	// when
	result, err := target.UpdateSimple(ctx, &testutils.Simple1, *testutils.Simple1.ToForm())
	// then
	assert.NoError(t, err)
	assert.Equal(t, &testutils.Simple1, result)
	simpleRepository.AssertExpectations(t)
}

func TestUpdateSimple_Error(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceWithMockDependencies(t)
	expectedError := errors.New("database error")
	// expect
	simpleRepository.On("Update", ctx, &testutils.Simple1).Return(nil, expectedError).Once()
	// when
	result, err := target.UpdateSimple(ctx, &testutils.Simple1, *testutils.Simple1.ToForm())
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
	target, simpleRepository := createSimpleServiceWithMockDependencies(t)
	// expect
	simpleRepository.On("Delete", ctx, testutils.Simple1.ID).Return(nil).Once()
	// when
	err := target.DeleteSimple(ctx, &testutils.Simple1)
	// then
	assert.NoError(t, err)
	simpleRepository.AssertExpectations(t)
}

func TestDeleteSimple_Error(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, simpleRepository := createSimpleServiceWithMockDependencies(t)
	expectedError := errors.New("database error")
	// expect
	simpleRepository.On("Delete", ctx, testutils.Simple1.ID).Return(expectedError).Once()
	// when
	err := target.DeleteSimple(ctx, &testutils.Simple1)
	// then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	simpleRepository.AssertExpectations(t)
}
