package repository

import (
	"time"

	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/telemetry"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx *gin.Context, user *model.User) (*model.User, error) {
	metrics := telemetry.GetMetrics()
	start := time.Now()

	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	metrics.RecordDBQuery(ctx, "create", time.Since(start).Seconds())
	return user, nil
}

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	user := &model.User{}
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return user, nil
}
