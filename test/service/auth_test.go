package service

import (
	"errors"
	"testing"
	"time"

	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/service"
	mockService "github.com/Verano-20/go-crud/test/mocks/service"
	"github.com/Verano-20/go-crud/test/testutils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSecret = []byte("test-secret-key")
	user1     = model.UserForm{Email: "test1@example.com", Password: "password1"}
	user2     = model.UserForm{Email: "test2@example.com", Password: "password2"}
)

func createAuthServiceWithMockDependencies(t *testing.T) (service.AuthService, *mockService.MockUserService) {
	userService := mockService.NewMockUserService()
	defer userService.AssertExpectations(t)
	target := service.NewAuthService(userService)
	return target, userService
}

/*
 * Validate User Credentials Tests
 */

func TestValidateUserCredentials_Success(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, userService := createAuthServiceWithMockDependencies(t)
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(user1.Password), bcrypt.DefaultCost)
	// expect
	userService.On("GetUserByEmail", ctx, user1.Email).Return(user1.ToModel(string(passwordHash)), nil).Once()
	// when
	user, err := target.ValidateUserCredentials(ctx, user1)
	// then
	assert.NoError(t, err)
	assert.Equal(t, user1.Email, user.Email)
}

func TestValidateUserCredentials_GetUserFailure(t *testing.T) {
	tests := []struct {
		testName      string
		userForm      model.UserForm
		expectedError string
	}{
		{
			testName:      "No Email",
			userForm:      model.UserForm{Password: user1.Password},
			expectedError: "user not found",
		},
		{
			testName:      "Invalid Email",
			userForm:      user2,
			expectedError: "user not found",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// given
			ctx, _ := testutils.CreateTestContext()
			target, userService := createAuthServiceWithMockDependencies(t)
			// expect
			userService.On("GetUserByEmail", ctx, test.userForm.Email).Return(nil, errors.New(test.expectedError)).Once()
			// when
			user, err := target.ValidateUserCredentials(ctx, test.userForm)
			// then
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.expectedError)
			assert.Nil(t, user)
		})
	}
}

func TestValidateUserCredentials_PasswordFailure(t *testing.T) {
	tests := []struct {
		testName       string
		userForm       model.UserForm
		passwordToHash string
		expectedError  string
	}{
		{
			testName:       "No Password",
			userForm:       model.UserForm{Email: user1.Email},
			passwordToHash: user1.Password,
			expectedError:  "hashedPassword is not the hash of the given password",
		},
		{
			testName:       "Invalid Password",
			userForm:       user1,
			passwordToHash: user1.Password + "1",
			expectedError:  "hashedPassword is not the hash of the given password",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// given
			ctx, _ := testutils.CreateTestContext()
			target, userService := createAuthServiceWithMockDependencies(t)
			passwordHash, _ := bcrypt.GenerateFromPassword([]byte(test.passwordToHash), bcrypt.DefaultCost)
			// expect
			userService.On("GetUserByEmail", ctx, test.userForm.Email).Return(user1.ToModel(string(passwordHash)), nil)
			// when
			user, err := target.ValidateUserCredentials(ctx, test.userForm)
			// then
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.expectedError)
			assert.Nil(t, user)
		})
	}
}

/*
 * Generate Token String Tests
 */

func TestGenerateTokenString_Success(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, _ := createAuthServiceWithMockDependencies(t)
	user := user1.ToModel(string(user1.Password))
	user.ID = 1234
	// when
	tokenString, err := target.GenerateTokenString(ctx, user, jwtSecret)
	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
	// and
	parser := jwt.NewParser()
	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.True(t, token.Valid)
	assert.Equal(t, float64(user.ID), token.Claims.(jwt.MapClaims)["sub"])
	// and
	iat := time.Unix(int64(token.Claims.(jwt.MapClaims)["iat"].(float64)), 0)
	exp := time.Unix(int64(token.Claims.(jwt.MapClaims)["exp"].(float64)), 0)
	assert.Equal(t, exp, iat.Add(time.Hour*24))
}

func TestGenerateTokenString_Failure_NilJwtSecret(t *testing.T) {
	// given
	ctx, _ := testutils.CreateTestContext()
	target, _ := createAuthServiceWithMockDependencies(t)
	// when
	tokenString, err := target.GenerateTokenString(ctx, user1.ToModel(string(user1.Password)), nil)
	// then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "jwtSecret is nil")
	assert.Empty(t, tokenString)
}
