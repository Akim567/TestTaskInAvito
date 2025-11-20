package service

import (
	"context"

	"TestTaskInAvito/internal/domain"
	"TestTaskInAvito/internal/repo"
)

// PRService — управление PR'ами.
type PRService interface {
	CreatePR(ctx context.Context, pr domain.PullRequest) (*domain.PullRequest, error)
	MergePR(ctx context.Context, prID string) (*domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*domain.PullRequest, string, error)
}

type prService struct {
	pr    repo.PRRepository
	users repo.UserRepository
	teams repo.TeamRepository
	tx    TxManager
}

func NewPRService(
	pr repo.PRRepository,
	users repo.UserRepository,
	teams repo.TeamRepository,
	tx TxManager,
) PRService {
	return &prService{
		pr:    pr,
		users: users,
		teams: teams,
		tx:    tx,
	}
}

// ===== реализация интерфейса =====

func (s *prService) CreatePR(ctx context.Context, pr domain.PullRequest) (*domain.PullRequest, error) {
	// TODO: реализовать:
	// - проверить, что PR ещё нет
	// - найти автора, его команду
	// - выбрать случайных ревьюверов
	// - сохранить PR и ревьюверов
	panic("not implemented")
}

func (s *prService) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	// TODO: реализовать идемпотентный merge
	panic("not implemented")
}

func (s *prService) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*domain.PullRequest, string, error) {
	// TODO: реализовать правила переназначения
	panic("not implemented")
}
