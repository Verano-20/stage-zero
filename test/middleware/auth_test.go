package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Verano-20/go-crud/internal/middleware"
	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/test/mocks/container"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

var (
	testContainer = container.NewContainerWithMockRepositories()
	user1         = model.User{ID: 1234567890, Email: "test1@example.com"}
	user2         = model.User{ID: 1234567891, Email: "test2@example.com"}
)

func TestAuthenticateRequest_Success(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", "/test", nil)
	validAuthHeader := "Bearer " + createHmacSignedToken(int64Ptr(time.Now().Add(time.Minute*1).Unix()), uintPtr(user1.ID))
	ctx.Request.Header.Set("Authorization", validAuthHeader)
	// and
	userRepository := testContainer.UserRepository
	userRepository.Create(ctx, &user1)
	target := middleware.NewAuthMiddleware([]byte("test-secret-key"), userRepository)
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
			gin.SetMode(gin.TestMode)
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest("GET", "/test", nil)
			ctx.Request.Header.Set("Authorization", test.authHeader)
			// and
			userRepository := testContainer.UserRepository
			userRepository.Create(ctx, &user1)
			target := middleware.NewAuthMiddleware([]byte("test-secret-key"), userRepository)
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
