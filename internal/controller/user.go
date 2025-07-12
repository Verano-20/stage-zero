package controller

import (
	"net/http"

	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	UserRepository *repository.UserRepository
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{UserRepository: repository.NewUserRepository(db)}
}

// SignUp godoc
// @Summary Sign up a new user
// @Description Create a new user with email and password
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.UserForm true "User email and password"
// @Success 200 {object} model.UserDTO "User created successfully"
// @Router /signup [post]
func (c *UserController) SignUp(ctx *gin.Context) {
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
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.UserForm true "User email and password"
// @Success 204
// @Router /login [post]
func (c *UserController) Login(ctx *gin.Context) {
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userForm.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// TODO: Generate and return JWT token

	ctx.JSON(http.StatusNoContent, nil)
}
