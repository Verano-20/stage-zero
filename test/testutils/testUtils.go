package testutils

import (
	"net/http/httptest"

	"github.com/Verano-20/stage-zero/internal/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	JwtSecret = []byte("test-secret-key")
	UserForm1 = model.UserForm{Email: "test1@example.com", Password: "password1"}
	UserForm2 = model.UserForm{Email: "test2@example.com", Password: "password2"}
	Simple1   = model.Simple{ID: 1, Name: "Simple 1"}
	Simple2   = model.Simple{ID: 2, Name: "Simple 2"}
)

func CreateTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	return ctx, recorder
}

func CreateTestContextWithAuthHeader(authHeader string) (*gin.Context, *httptest.ResponseRecorder) {
	ctx, recorder := CreateTestContext()
	ctx.Request = httptest.NewRequest("GET", "/test", nil)
	ctx.Request.Header.Set("Authorization", authHeader)
	return ctx, recorder
}

func GetUserWithPasswordHashFromForm(userForm model.UserForm) *model.User {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(userForm.Password), bcrypt.DefaultCost)
	return userForm.ToModel(string(passwordHash))
}
