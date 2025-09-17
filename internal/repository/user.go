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

	metrics.RecordDBQuery(ctx, "create_user", time.Since(start).Seconds())
	metrics.UpdateUserCount(ctx, 1)
	return user, nil
}

func (r *UserRepository) GetByID(ctx *gin.Context, id uint) (*model.User, error) {
	metrics := telemetry.GetMetrics()
	start := time.Now()

	user := &model.User{}
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	metrics.RecordDBQuery(ctx, "get_user_by_id", time.Since(start).Seconds())
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx *gin.Context, email string) (*model.User, error) {
	metrics := telemetry.GetMetrics()
	start := time.Now()

	user := &model.User{}
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	metrics.RecordDBQuery(ctx, "get_user_by_email", time.Since(start).Seconds())
	return user, nil
}
