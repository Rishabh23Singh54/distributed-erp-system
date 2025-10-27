package service

import (
	"context"
	"strconv"

	"github.com/Rishabh23Singh54/distributed-erp-system/services/user-service/internal/domain"
	"github.com/Rishabh23Singh54/distributed-erp-system/services/user-service/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, user *domain.User) error {
	return s.repo.Update(ctx, user)
}

// func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
// 	return s.repo.List(ctx, limit, offset)
// }

// func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
// 	return s.repo.Delete(ctx, userID)
// }
