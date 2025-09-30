package service

import (
	"errors"
	"strings"
	"testing"

	apiError "github.com/Verano-20/go-crud/internal/err"
	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/service"
	"github.com/Verano-20/go-crud/test/mocks/repository"
	"github.com/Verano-20/go-crud/test/testutils"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createUserServiceWithMockDependencies(t *testing.T) (service.UserService, *repository.MockUserRepository) {
	mockRepo := repository.NewMockUserRepository()
	defer mockRepo.AssertExpectations(t)
	target := service.NewUserService(mockRepo)
	return target, mockRepo
}

/*
 * Create User Tests
 */

func TestCreateUser_Success(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, userRepository := createUserServiceWithMockDependencies(t)
	expectedUser := testutils.GetUserWithPasswordHashFromForm(testutils.UserForm1)
	// expect
	userRepository.On("Create", ctx, mock.MatchedBy(func(user *model.User) bool {
		return user.Email == expectedUser.Email
	})).Return(expectedUser, nil).Once()
	// when
	result, err := target.CreateUser(ctx, testutils.UserForm1)
	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
}

func TestCreateUser_PasswordHashError(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, userRepository := createUserServiceWithMockDependencies(t)
	userForm := model.UserForm{Email: testutils.UserForm1.Email, Password: strings.Repeat("A", 73)} // length > 72 causes error
	// when
	result, err := target.CreateUser(ctx, userForm)
	// then
	apiErr, ok := err.(*apiError.ApiError)
	assert.True(t, ok)
	assert.Equal(t, apiError.ErrorTypePasswordHash, apiErr.Type)
	assert.Contains(t, err.Error(), "password length exceeds 72 bytes")
	assert.Nil(t, result)
	// and
	userRepository.AssertNotCalled(t, "Create", ctx, mock.Anything)
}

func TestCreateUser_EmailExists(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, userRepository := createUserServiceWithMockDependencies(t)
	userForm := testutils.UserForm1
	errorMessage := "email already exists"
	expectedError := &pgconn.PgError{Code: "23505", Message: errorMessage}
	// expect
	userRepository.On("Create", ctx, mock.MatchedBy(func(user *model.User) bool {
		return user.Email == userForm.Email
	})).Return(nil, expectedError).Once()
	// when
	result, err := target.CreateUser(ctx, userForm)
	// then
	apiErr, ok := err.(*apiError.ApiError)
	assert.True(t, ok)
	assert.Equal(t, apiError.ErrorTypeEmailExists, apiErr.Type)
	assert.Contains(t, err.Error(), errorMessage)
	assert.Nil(t, result)
}

func TestCreateUser_DatabaseError(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, userRepository := createUserServiceWithMockDependencies(t)
	expectedError := errors.New("database error")
	// expect
	userRepository.On("Create", ctx, mock.MatchedBy(func(user *model.User) bool {
		return user.Email == testutils.UserForm1.Email
	})).Return(nil, expectedError).Once()
	// when
	result, err := target.CreateUser(ctx, testutils.UserForm1)
	// then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
}

/*
 * Get User By Email Tests
 */

func TestGetUserByEmail_Success(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, userRepository := createUserServiceWithMockDependencies(t)
	expectedUser := testutils.GetUserWithPasswordHashFromForm(testutils.UserForm1)
	// expect
	userRepository.On("GetByEmail", ctx, testutils.UserForm1.Email).Return(expectedUser, nil).Once()
	// when
	result, err := target.GetUserByEmail(ctx, testutils.UserForm1.Email)
	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
}

func TestGetUserByEmail_Error(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, userRepository := createUserServiceWithMockDependencies(t)
	expectedError := errors.New("database error")
	// expect
	userRepository.On("GetByEmail", ctx, testutils.UserForm1.Email).Return(nil, expectedError).Once()
	// when
	result, err := target.GetUserByEmail(ctx, testutils.UserForm1.Email)
	// then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
}
