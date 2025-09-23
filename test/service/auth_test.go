package service

import (
	"net/http/httptest"
	"testing"

	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/test/mocks/container"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	testContainer = container.NewContainerWithMockRepositories()
	user1         = model.UserForm{Email: "test1@example.com", Password: "password1"}
	user2         = model.UserForm{Email: "test2@example.com", Password: "password2"}
)

func TestValidateUserCredentials_Success(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", "/test", nil)
	// and
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(user1.Password), bcrypt.DefaultCost)
	userRepository := testContainer.UserRepository
	userRepository.Create(ctx, user1.ToModel(string(passwordHash)))
	target := testContainer.AuthService
	// when
	user, err := target.ValidateUserCredentials(ctx, user1)
	// then
	assert.Nil(t, err)
	assert.Equal(t, user1.Email, user.Email)
}

func TestValidateUserCredentials_Failure(t *testing.T) {
	tests := []struct {
		testName       string
		userForm       model.UserForm
		passwordToHash string
		expectedError  string
	}{
		{
			testName:       "No Email",
			userForm:       model.UserForm{Password: user1.Password},
			passwordToHash: user1.Password,
			expectedError:  "user not found",
		},
		{
			testName:       "Invalid Email",
			userForm:       user2,
			passwordToHash: user1.Password,
			expectedError:  "user not found",
		},
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
			gin.SetMode(gin.TestMode)
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest("GET", "/test", nil)
			// and
			passwordHash, _ := bcrypt.GenerateFromPassword([]byte(test.passwordToHash), bcrypt.DefaultCost)
			userRepository := testContainer.UserRepository
			userRepository.Create(ctx, user1.ToModel(string(passwordHash)))
			target := testContainer.AuthService
			// when
			user, err := target.ValidateUserCredentials(ctx, test.userForm)
			// then
			assert.Contains(t, err.Error(), test.expectedError)
			assert.Nil(t, user)
		})
	}
}
