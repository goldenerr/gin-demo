package service

import (
	"context"
	"gin-demo/internal/models"
	"gin-demo/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 可以添加其他方法，如 CreateUser, UpdateUser, DeleteUser 等

// 添加获取用户列表的方法
func (s *UserService) GetUsers(ctx context.Context, query models.UserQuery) ([]models.User, int64, error) {
	users, total, err := s.userRepo.GetUsers(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
