package controller

import (
	"errors"
	"net/http"

	"github.com/Verano-20/stage-zero/internal/config"
	"github.com/Verano-20/stage-zero/internal/err"
	"github.com/Verano-20/stage-zero/internal/model"
	"github.com/Verano-20/stage-zero/internal/response"
	"github.com/Verano-20/stage-zero/internal/service"
	"github.com/Verano-20/stage-zero/internal/telemetry"
	"github.com/Verano-20/stage-zero/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserService service.UserService
	AuthService service.AuthService
}

func NewAuthController(userService service.UserService, authService service.AuthService) *AuthController {
	return &AuthController{UserService: userService, AuthService: authService}
}

// SignUp godoc
// @Summary Sign up a new user
// @Description Create a new user with email and password. The email must be unique. The password must be at least 8 characters long.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body model.UserForm true "User details"
// @Success 201 {object} model.UserDTO "User created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request format or validation failed"
// @Failure 409 {object} response.ErrorResponse "User already exists"
// @Failure 500 {object} response.ErrorResponse "Internal server error during user creation"
// @Router /auth/signup [post]
func (c *AuthController) SignUp(ctx *gin.Context) {
	metrics := telemetry.GetMetrics()

	var userForm model.UserForm
	if formErr := ctx.ShouldBindJSON(&userForm); formErr != nil {
		utils.HandleBindingErrors(ctx, formErr, "signup")
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
	ctx.JSON(http.StatusCreated, response.ApiResponse{Message: "User created successfully", Data: user.ToDTO()})
}

// Login godoc
// @Summary Authenticate user and generate JWT token
// @Description Authenticate a user with email and password credentials. Returns a JWT token upon successful authentication that can be used for subsequent API calls.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body model.UserForm true "User login credentials"
// @Success 200 {object} response.ApiResponse "Authentication successful, returns JWT token"
// @Failure 400 {object} response.ErrorResponse "Invalid request format or validation failed"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Failure 500 {object} response.ErrorResponse "Internal server error during authentication"
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	metrics := telemetry.GetMetrics()
	config := config.Get()

	var userForm model.UserForm
	if formErr := ctx.ShouldBindJSON(&userForm); formErr != nil {
		utils.HandleBindingErrors(ctx, formErr, "login")
		return
	}

	user, err := c.AuthService.ValidateUserCredentials(ctx, userForm)
	if err != nil {
		metrics.RecordAuthAttempt(ctx, false, "login")
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	tokenString, err := c.AuthService.GenerateTokenString(ctx, user, config.GetJwtSecret())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	metrics.RecordAuthAttempt(ctx, true, "login")
	ctx.JSON(http.StatusOK, response.ApiResponse{Message: "Login successful", Data: map[string]string{"token": tokenString}})
}
