package service

import (
	"context"

	"TestTaskInAvito/internal/domain"
	"TestTaskInAvito/internal/repo"
)

// TeamService описывает операции для работы с командами
type TeamService interface {
	CreateTeam(ctx context.Context, team domain.Team) (*domain.Team, error)
	GetTeam(ctx context.Context, name string) (*domain.Team, error)
}

type teamService struct {
	teams repo.TeamRepository
	users repo.UserRepository
	tx    TxManager
}

func NewTeamService(teams repo.TeamRepository, users repo.UserRepository, tx TxManager) TeamService {
	return &teamService{
		teams: teams,
		users: users,
		tx:    tx,
	}
}
