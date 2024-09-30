package service

import (
	"context"
	"gin-demo/internal/middleware"
	"gin-demo/internal/models"
	"gin-demo/internal/repository"

	"go.uber.org/zap"
)

type UserService struct {
	userRepo *repository.UserRepository
	logger   *zap.Logger
}

func NewUserService(userRepo *repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	requestID := middleware.GetRequestID(ctx)
	s.logger.Info("Service: Getting user by ID",
		zap.String("request_id", requestID),
		zap.String("id", id))
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		s.logger.Error("Service: Failed to get user by ID",
			zap.String("request_id", requestID),
			zap.String("id", id),
			zap.Error(err))
		return nil, err
	}
	s.logger.Info("Service: User retrieved successfully",
		zap.String("request_id", requestID),
		zap.String("id", id))
	return user, nil
}

// 可以添加其他方法，如 CreateUser, UpdateUser, DeleteUser 等

// 添加获取用户列表的方法
func (s *UserService) GetUsers(ctx context.Context, query models.UserQuery) ([]models.User, int64, error) {
	requestID := middleware.GetRequestID(ctx)
	s.logger.Info("Service: Getting users",
		zap.String("request_id", requestID),
		zap.Any("query", query))
	users, total, err := s.userRepo.GetUsers(ctx, query)
	if err != nil {
		s.logger.Error("Service: Failed to get users",
			zap.String("request_id", requestID),
			zap.Error(err))
		return nil, 0, err
	}
	s.logger.Info("Service: Users retrieved successfully",
		zap.String("request_id", requestID),
		zap.Int("total", int(total)),
		zap.Int("retrieved", len(users)))
	return users, total, nil
}
