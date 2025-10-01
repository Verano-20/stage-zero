package err

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	Type string
	Err  error
}

func (e *ApiError) Error() string {
	return e.Err.Error()
}

func (e *ApiError) Unwrap() error {
	return e.Err
}

const (
	ErrorTypePasswordHash = "password_hash_failure"
	ErrorTypeEmailExists  = "email_already_exists"
)

func NewPasswordHashError(err error) *ApiError {
	return &ApiError{
		Type: ErrorTypePasswordHash,
		Err:  err,
	}
}

func NewEmailExistsError(err error) *ApiError {
	return &ApiError{
		Type: ErrorTypeEmailExists,
		Err:  err,
	}
}

func GetValidationErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}
