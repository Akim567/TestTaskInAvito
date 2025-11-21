package service

import (
	"TestTaskInAvito/internal/repo"
	"context"
	"database/sql"
	"errors"

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

// SetIsActive меняет активность пользователя и возвращает обновлённого юзера.
func (s *userService) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	var updated *domain.User

	err := s.tx.Do(ctx, func(txCtx context.Context) error {
		// проверяем, что пользователь существует
		_, err := s.users.GetByID(txCtx, userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.NewNotFoundError("user")
			}
			return err
		}

		if err := s.users.SetIsActive(txCtx, userID, isActive); err != nil {
			return err
		}

		updated, err = s.users.GetByID(txCtx, userID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return updated, nil
}

// GetUserReviews возвращает все PR, где пользователь является ревьювером.
func (s *userService) GetUserReviews(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	// сначала убеждаемся, что пользователь существует
	_, err := s.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("user")
		}
		return nil, err
	}

	prs, err := s.pr.GetByReviewer(ctx, userID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
