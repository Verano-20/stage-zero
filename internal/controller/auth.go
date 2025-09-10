package controller

import (
	"errors"
	"net/http"

	"github.com/Verano-20/go-crud/internal/err"
	"github.com/Verano-20/go-crud/internal/logger"
	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/response"
	"github.com/Verano-20/go-crud/internal/service"
	"github.com/Verano-20/go-crud/internal/telemetry"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthController struct {
	UserService *service.UserService
	AuthService *service.AuthService
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{UserService: service.NewUserService(db), AuthService: service.NewAuthService(db)}
}

// SignUp godoc
// @Summary Sign up a new user
// @Description Create a new user with email and password. The email must be unique.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body model.UserForm true "User details" example({"email": "user@example.com", "password": "securePassword123"})
// @Success 201 {object} model.UserDTO "User account created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request format or validation failed" example({"error": "Invalid request format"})
// @Failure 409 {object} response.ErrorResponse "Email address already exists" example({"error": "Email already exists"})
// @Failure 500 {object} response.ErrorResponse "Internal server error during user creation" example({"error": "Failed to create user"})
// @Router /auth/signup [post]
func (c *AuthController) SignUp(ctx *gin.Context) {
	metrics := telemetry.GetMetrics()
	log := logger.GetFromContext(ctx)

	var userForm model.UserForm
	if formErr := ctx.ShouldBindJSON(&userForm); formErr != nil {
		log.Warn("Invalid signup request format", zap.Error(formErr))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid request format"})
		return
	}

	user, createErr := c.UserService.CreateUser(ctx, userForm)
	if createErr != nil {
		metrics.RecordAuthAttempt(ctx, false, "signup")
		var apiError *err.ApiError
		if errors.As(createErr, &apiError) {
			switch apiError.Type {
			case err.ErrorTypePasswordHash:
				ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to process password"})
				return
			case err.ErrorTypeEmailExists:
				ctx.JSON(http.StatusConflict, response.ErrorResponse{Error: "User already exists"})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to create user"})
		return
	}

	metrics.RecordAuthAttempt(ctx, true, "signup")
	ctx.JSON(http.StatusCreated, user.ToDTO())
}

// Login godoc
// @Summary Authenticate user and generate JWT token
// @Description Authenticate a user with email and password credentials. Returns a JWT token upon successful authentication that can be used for subsequent API calls.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body model.UserForm true "User login credentials" example({"email": "user@example.com", "password": "securePassword123"})
// @Success 200 {object} response.ApiResponse "Authentication successful, returns JWT token" example({"message": "Login successful", "data": {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}})
// @Failure 400 {object} response.ErrorResponse "Invalid request format" example({"error": "Invalid request format"})
// @Failure 401 {object} response.ErrorResponse "Invalid email or password" example({"error": "Invalid credentials"})
// @Failure 500 {object} response.ErrorResponse "Internal server error during authentication" example({"error": "Failed to generate token"})
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	metrics := telemetry.GetMetrics()
	log := logger.GetFromContext(ctx)

	var userForm model.UserForm
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		log.Warn("Invalid login request format", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid request format"})
		return
	}

	user, err := c.AuthService.ValidateUserCredentials(ctx, userForm)
	if err != nil {
		metrics.RecordAuthAttempt(ctx, false, "login")
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	tokenString, err := c.AuthService.GenerateTokenString(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	metrics.RecordAuthAttempt(ctx, true, "login")
	ctx.JSON(http.StatusOK, response.ApiResponse{Message: "Login successful", Data: map[string]string{"token": tokenString}})
}
