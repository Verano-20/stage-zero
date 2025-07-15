package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/Verano-20/go-crud/internal/config"
	"github.com/Verano-20/go-crud/internal/logger"
	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/repository"
	"github.com/Verano-20/go-crud/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	UserRepository *repository.UserRepository
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{UserRepository: repository.NewUserRepository(db)}
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
	log := logger.GetFromContext(ctx)

	var userForm model.UserForm
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		log.Warn("Invalid signup request format", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid request format"})
		return
	}

	log.Debug("Processing User signup...", zap.Object("user", &userForm))

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userForm.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password", zap.Object("user", &userForm), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to process password"})
		return
	}

	user, err := c.UserRepository.Create(userForm.ToModel(string(passwordHash)))
	if err != nil {
		var pgErr *pgconn.PgError
		// Check if the error is a unique constraint violation (SQLSTATE 23505)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Warn("User signup failed - email already in use", zap.Object("user", &userForm), zap.Error(err))
			ctx.JSON(http.StatusConflict, response.ErrorResponse{Error: "User already exists"})
			return
		}
		log.Error("Failed to create User in database", zap.Object("user", &userForm), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to create User"})
		return
	}

	log.Debug("User signup successful", zap.Object("user", user))

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
	log := logger.GetFromContext(ctx)

	var userForm model.UserForm
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		log.Warn("Invalid login request format", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid request format"})
		return
	}

	log.Debug("Processing User login...", zap.Object("user", &userForm))

	user, err := c.UserRepository.GetByEmail(userForm.Email)
	if err != nil {
		log.Warn("Login failed - User not found", zap.Object("user", &userForm), zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userForm.Password)); err != nil {
		log.Warn("Login failed - invalid password", zap.Object("user", &userForm), zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(config.GetJwtSecret())
	if err != nil {
		log.Error("Failed to generate JWT token", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	log.Debug("User login successful", zap.Object("user", user))

	ctx.JSON(http.StatusOK, response.ApiResponse{Message: "Login successful", Data: map[string]string{"token": tokenString}})
}
