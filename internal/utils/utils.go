package utils

import (
	"net/http"

	"github.com/Verano-20/stage-zero/internal/err"
	"github.com/Verano-20/stage-zero/internal/logger"
	"github.com/Verano-20/stage-zero/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// Handles binding errors and sets a Bad Request response with validation errors or a generic error message.
// ctx *gin.Context - The Gin context containing the request and response.
// formErr error - The error returned by the binding process.
// action string - The action being performed (e.g., "signup", "login", "create", "update", "delete").
func HandleBindingErrors(ctx *gin.Context, formErr error, action string) {
	log := logger.GetFromContext(ctx)
	log.Warn("Invalid request payload", zap.String("action", action), zap.Error(formErr))
	if validationErrors, ok := formErr.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, e := range validationErrors {
			errors[e.Field()] = err.GetValidationErrorMessage(e)
		}
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Validation failed", Details: errors})
		return
	}

	ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid request format"})
}
