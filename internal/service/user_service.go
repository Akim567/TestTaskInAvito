package service

import (
	"TestTaskInAvito/internal/repo"
	"context"

	"TestTaskInAvito/internal/domain"
)

// UserService — управление активностью пользователя и его PR-ревью.
type UserService interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
	GetUserReviews(ctx context.Context, userID string) ([]domain.PullRequest, error)
}

type userService struct {
	users repo.UserRepository
	pr    repo.PRRepository
	tx    TxManager
}

func NewUserService(
	users repo.UserRepository,
	pr repo.PRRepository,
	tx TxManager,
) UserService {
	return &userService{
		users: users,
		pr:    pr,
		tx:    tx,
	}
}
