package service

import (
	"errors"
	"time"

	"github.com/Verano-20/go-crud/internal/logger"
	"github.com/Verano-20/go-crud/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	ValidateUserCredentials(ctx *gin.Context, userForm model.UserForm) (user *model.User, err error)
	GenerateTokenString(ctx *gin.Context, user *model.User, jwtSecret []byte) (tokenString string, err error)
}

type authService struct {
	UserService UserService
}

var _ AuthService = &authService{}

func NewAuthService(userService UserService) AuthService {
	return &authService{UserService: userService}
}

func (s *authService) ValidateUserCredentials(ctx *gin.Context, userForm model.UserForm) (user *model.User, err error) {
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

func (s *authService) GenerateTokenString(ctx *gin.Context, user *model.User, jwtSecret []byte) (tokenString string, err error) {
	log := logger.GetFromContext(ctx)
	log.Debug("Generating JWT token...", zap.Object("user", user))

	if jwtSecret == nil {
		err = errors.New("jwtSecret is nil")
		log.Error("JWT secret is nil", zap.Object("user", user), zap.Error(err))
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	if tokenString, err = token.SignedString(jwtSecret); err != nil {
		log.Error("Failed to generate JWT token", zap.Object("user", user), zap.Error(err))
		return "", err
	}

	log.Debug("JWT token generated successfully", zap.Object("user", user))
	return tokenString, nil
}
