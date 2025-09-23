package service

import (
	"errors"

	"github.com/Verano-20/go-crud/internal/err"
	"github.com/Verano-20/go-crud/internal/logger"
	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepository repository.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{UserRepository: repository.NewUserRepository(db)}
}

func (s *UserService) CreateUser(ctx *gin.Context, userForm model.UserForm) (user *model.User, createErr error) {
	log := logger.GetFromContext(ctx)

	log.Debug("Creating User...", zap.Object("user", &userForm))

	passwordHash, hashErr := bcrypt.GenerateFromPassword([]byte(userForm.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		log.Error("Failed to hash password", zap.Object("user", &userForm), zap.Error(hashErr))
		return nil, err.NewPasswordHashError(hashErr)
	}

	user, dbErr := s.UserRepository.Create(ctx, userForm.ToModel(string(passwordHash)))
	if dbErr != nil {
		var pgErr *pgconn.PgError
		// Check if the error is a unique constraint violation (SQLSTATE 23505)
		if errors.As(dbErr, &pgErr) && pgErr.Code == "23505" {
			log.Warn("User creation failed - email already in use", zap.Object("user", &userForm), zap.Error(dbErr))
			return nil, err.NewEmailExistsError(dbErr)
		}
		log.Error("Failed to create User in database", zap.Object("user", &userForm), zap.Error(dbErr))
		return nil, dbErr
	}

	log.Debug("User created successfully", zap.Object("user", user))
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx *gin.Context, email string) (user *model.User, err error) {
	log := logger.GetFromContext(ctx)

	log.Debug("Getting User by email...", zap.String("email", email))

	user, err = s.UserRepository.GetByEmail(ctx, email)
	if err != nil {
		log.Warn("Failed to find User with email", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	log.Debug("User retrieved successfully", zap.Object("user", user))
	return user, nil
}
