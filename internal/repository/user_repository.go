package repository

import (
	"context"
	"gin-demo/internal/middleware"
	"gin-demo/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	requestID := middleware.GetRequestID(ctx)
	var user models.User
	r.logger.Info("Repository: Getting user by ID",
		zap.String("request_id", requestID),
		zap.String("id", id))
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		r.logger.Error("Repository: Failed to get user by ID",
			zap.String("request_id", requestID),
			zap.String("id", id),
			zap.Error(err))
		return nil, err
	}
	r.logger.Info("Repository: User retrieved by ID",
		zap.String("request_id", requestID),
		zap.String("id", id),
		zap.Any("user", user))
	return &user, nil
}

// 添加获取用户列表的方法
func (r *UserRepository) GetUsers(ctx context.Context, query models.UserQuery) ([]models.User, int64, error) {
	requestID := middleware.GetRequestID(ctx)
	var users []models.User
	var total int64

	db := r.db.WithContext(ctx).Model(&models.User{})

	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Email != "" {
		db = db.Where("email LIKE ?", "%"+query.Email+"%")
	}

	r.logger.Info("Repository: Counting users",
		zap.String("request_id", requestID),
		zap.Any("query", query))
	err := db.Count(&total).Error
	if err != nil {
		r.logger.Error("Repository: Failed to count users",
			zap.String("request_id", requestID),
			zap.Error(err))
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.Size
	r.logger.Info("Repository: Getting users",
		zap.String("request_id", requestID),
		zap.Any("query", query),
		zap.Int("offset", offset),
		zap.Int("limit", query.Size))
	err = db.Offset(offset).Limit(query.Size).Find(&users).Error
	if err != nil {
		r.logger.Error("Repository: Failed to get users",
			zap.String("request_id", requestID),
			zap.Error(err))
		return nil, 0, err
	}

	r.logger.Info("Repository: Users retrieved",
		zap.String("request_id", requestID),
		zap.Int("count", len(users)),
		zap.Int64("total", total))
	return users, total, nil
}

// 实现其他数据库操作...
