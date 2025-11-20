package repo

import (
	"TestTaskInAvito/internal/domain"
	"context"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team domain.Team) error
	GetByName(ctx context.Context, name string) (*domain.Team, error)
}

type UserRepository interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetActiveTeamMembers(ctx context.Context, teamName string) ([]domain.User, error)
	CreateOrUpdate(ctx context.Context, user domain.User) error
}

type PRRepository interface {
	Create(ctx context.Context, pr domain.PullRequest) error
	GetByID(ctx context.Context, id string) (*domain.PullRequest, error)
	Update(ctx context.Context, pr domain.PullRequest) error
	GetByReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error)
}
