package controller

import (
	"net/http"
	"time"

	"github.com/Verano-20/go-crud/internal/config"
	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
// @Description Create a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.UserForm true "User email and password"
// @Success 200 {object} model.UserDTO "User created successfully"
// @Router /signup [post]
func (c *AuthController) SignUp(ctx *gin.Context) {
	var userForm model.UserForm
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userForm.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := c.UserRepository.Create(userForm.ToModel(string(passwordHash)))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user.ToDTO())
}

// Login godoc
// @Summary Login a user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.UserForm true "User email and password"
// @Success 204
// @Router /login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var userForm model.UserForm
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.UserRepository.GetByEmail(userForm.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userForm.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(config.GetJwtSecret())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}
