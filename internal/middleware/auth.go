package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Verano-20/go-crud/internal/logger"
	"github.com/Verano-20/go-crud/internal/repository"
	"github.com/Verano-20/go-crud/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	jwtSecret      []byte
	userRepository *repository.UserRepository
}

func NewAuthMiddleware(jwtSecret []byte, db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:      jwtSecret,
		userRepository: repository.NewUserRepository(db),
	}
}

func (m *AuthMiddleware) AuthenticateRequest(ctx *gin.Context) {
	log := logger.GetFromContext(ctx)

	log.Debug("Authentcating request...")

	token, err := m.validateToken(ctx)
	if err != nil {
		log.Warn("Token validation failed", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
		ctx.Abort()
		return
	}

	if err := m.validateClaims(ctx, token); err != nil {
		log.Warn("Token claims validation failed", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
		ctx.Abort()
		return
	}

	log.Debug("Authentication successful")
	ctx.Next()
}

func (m *AuthMiddleware) validateToken(ctx *gin.Context) (*jwt.Token, error) {
	log := logger.GetFromContext(ctx)

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		log.Warn("Missing authorization header")
		return nil, errors.New("authorization header required")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		log.Warn("Invalid authorization header format", zap.String("header", authHeader))
		return nil, errors.New("invalid authorization header format")
	}

	tokenString := tokenParts[1]
	log.Debug("Parsing JWT token...")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Warn("Invalid JWT signing method",
				zap.String("method", token.Method.Alg()))
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return m.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		log.Warn("JWT token parsing failed", zap.Error(err))
		return nil, errors.New("invalid token")
	}

	log.Debug("JWT token parsed successfully")
	return token, nil
}

func (m *AuthMiddleware) validateClaims(ctx *gin.Context, token *jwt.Token) error {
	log := logger.GetFromContext(ctx)

	log.Debug("Validating JWT token claims...")
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Warn("Invalid JWT token claims format")
		return errors.New("invalid token claims")
	}

	if claims["exp"] == nil || claims["sub"] == nil {
		log.Warn("Missing required JWT claims",
			zap.Bool("has_exp", claims["exp"] != nil),
			zap.Bool("has_sub", claims["sub"] != nil))
		return errors.New("invalid token claims")
	}

	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		log.Warn("JWT token expired",
			zap.Float64("exp", claims["exp"].(float64)),
			zap.Int64("now", time.Now().Unix()))
		return errors.New("token expired")
	}

	userID := uint(claims["sub"].(float64))
	user, err := m.userRepository.GetByID(ctx, userID)
	if err != nil {
		log.Warn("User not found during token validation",
			zap.Uint("user_id", userID),
			zap.Error(err))
		return errors.New("invalid user id")
	}

	log.Debug("JWT token claims validated successfully",
		zap.Float64("exp", claims["exp"].(float64)),
		zap.Float64("now", float64(time.Now().Unix())),
		zap.Uint("user_id", user.ID),
		zap.String("email", user.Email))

	ctx.Set("user_id", user.ID)
	ctx.Set("user_email", user.Email)

	return nil
}
