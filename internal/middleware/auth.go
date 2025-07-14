package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Verano-20/go-crud/internal/config"
	"github.com/Verano-20/go-crud/internal/repository"
	"github.com/Verano-20/go-crud/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	userRepository *repository.UserRepository
}

func NewAuthMiddleware(db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{userRepository: repository.NewUserRepository(db)}
}

func (m *AuthMiddleware) AuthenticateRequest(ctx *gin.Context) {
	token, err := m.validateToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
		ctx.Abort()
		return
	}

	if err := m.validateClaims(token); err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
		ctx.Abort()
		return
	}

	ctx.Next()
}

func (m *AuthMiddleware) validateToken(ctx *gin.Context) (*jwt.Token, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header required")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return nil, errors.New("invalid authorization header format")
	}

	tokenString := tokenParts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(config.GetJwtSecret()), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func (m *AuthMiddleware) validateClaims(token *jwt.Token) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	if claims["exp"] == nil || claims["sub"] == nil {
		return errors.New("invalid token claims")
	}

	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		return errors.New("token expired")
	}

	if _, err := m.userRepository.GetByID(uint(claims["sub"].(float64))); err != nil {
		return errors.New("invalid user id")
	}

	return nil
}
