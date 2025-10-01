package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Verano-20/stage-zero/internal/middleware"
	"github.com/Verano-20/stage-zero/internal/model"
	"github.com/Verano-20/stage-zero/test/mocks/repository"
	"github.com/Verano-20/stage-zero/test/testutils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

var (
	user1 = model.User{ID: 1, Email: "test1@example.com"}
	user2 = model.User{ID: 2, Email: "test2@example.com"}
)

func createMiddlewareAndMockRepo(t *testing.T) (*middleware.AuthMiddleware, *repository.MockUserRepository) {
	userRepository := repository.NewMockUserRepository()
	defer userRepository.AssertExpectations(t)
	target := middleware.NewAuthMiddleware([]byte("test-secret-key"), userRepository)
	return target, userRepository
}

func TestAuthenticateRequest_Success(t *testing.T) {
	// given
	validAuthHeader := "Bearer " + createHmacSignedToken(int64Ptr(time.Now().Add(time.Minute*1).Unix()), uintPtr(user1.ID))
	ctx, recorder := testutils.CreateTestContextWithAuthHeader(validAuthHeader)
	target, userRepository := createMiddlewareAndMockRepo(t)
	// expect
	userRepository.On("GetByID", ctx, user1.ID).Return(&user1, nil).Once()
	// when
	target.AuthenticateRequest(ctx)
	// then
	assert.False(t, ctx.IsAborted())
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, user1.ID, ctx.GetUint("user_id"))
	assert.Equal(t, user1.Email, ctx.GetString("user_email"))
}

func TestAuthenticateRequest_Failure(t *testing.T) {
	tests := []struct {
		testName             string
		authHeader           string
		expectedStatusCode   int
		expectedErrorMessage string
	}{
		{
			testName:             "Missing Authorization Header",
			authHeader:           "",
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "authorization header required",
		},
		{
			testName:             "Invalid Authorization Header",
			authHeader:           "not a bearer token",
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "invalid authorization header format",
		},
		{
			testName:             "Invalid Token Signing Method",
			authHeader:           "Bearer " + createRsaSignedToken(),
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "invalid token",
		},
		{
			testName:             "Invalid Token",
			authHeader:           "Bearer invalid-token",
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "invalid token",
		},
		{
			testName:             "Expired Token",
			authHeader:           "Bearer " + createHmacSignedToken(int64Ptr(time.Now().Add(-time.Minute*1).Unix()), uintPtr(user1.ID)),
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "token expired",
		},
		{
			testName:             "Sub Missing",
			authHeader:           "Bearer " + createHmacSignedToken(int64Ptr(time.Now().Add(time.Minute*1).Unix()), nil),
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "invalid token claims",
		},
		{
			testName:             "User Not Found",
			authHeader:           "Bearer " + createHmacSignedToken(int64Ptr(time.Now().Add(time.Minute*1).Unix()), uintPtr(user2.ID)),
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "invalid user id",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// given
			ctx, recorder := testutils.CreateTestContextWithAuthHeader(test.authHeader)
			target, userRepository := createMiddlewareAndMockRepo(t)
			// expect
			userRepository.On("GetByID", ctx, user1.ID).Return(&user1, nil).Maybe()
			userRepository.On("GetByID", ctx, user2.ID).Return(nil, errors.New("user not found")).Maybe()
			// when
			target.AuthenticateRequest(ctx)
			// then
			assert.True(t, ctx.IsAborted())
			assert.Equal(t, test.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), test.expectedErrorMessage)
		})
	}
}

func createRsaSignedToken() string {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": user1.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func createHmacSignedToken(exp *int64, sub *uint) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": exp,
	})

	tokenString, err := token.SignedString([]byte("test-secret-key"))
	if err != nil {
		panic(err)
	}
	return tokenString
}

func int64Ptr(i int64) *int64 { return &i }
func uintPtr(i uint) *uint    { return &i }
