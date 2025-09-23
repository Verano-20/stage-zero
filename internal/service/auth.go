package service

import (
	"time"

	"github.com/Verano-20/go-crud/internal/config"
	"github.com/Verano-20/go-crud/internal/logger"
	"github.com/Verano-20/go-crud/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserService *UserService
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{UserService: userService}
}

func (s *AuthService) ValidateUserCredentials(ctx *gin.Context, userForm model.UserForm) (user *model.User, err error) {
	log := logger.GetFromContext(ctx)

	log.Debug("Validating User credentials...", zap.Object("userForm", &userForm))

	if user, err = s.UserService.GetUserByEmail(ctx, userForm.Email); err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userForm.Password)); err != nil {
		log.Warn("Login failed - invalid password", zap.Object("userForm", &userForm), zap.Error(err))
		return nil, err
	}

	log.Debug("User credentials valid", zap.Object("user", user))
	return user, nil
}

func (s *AuthService) GenerateTokenString(ctx *gin.Context, user *model.User) (tokenString string, err error) {
	config := config.Get()
	log := logger.GetFromContext(ctx)

	log.Debug("Generating JWT token...", zap.Object("user", user))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	if tokenString, err = token.SignedString(config.GetJwtSecret()); err != nil {
		log.Error("Failed to generate JWT token", zap.Object("user", user), zap.Error(err))
		return "", err
	}

	log.Debug("JWT token generated successfully", zap.Object("user", user))
	return tokenString, nil
}
