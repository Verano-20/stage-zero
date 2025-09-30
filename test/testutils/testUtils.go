package testutils

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
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
