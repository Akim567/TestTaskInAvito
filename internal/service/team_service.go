package service

import (
	"context"
	"database/sql"
	"errors"

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

// ===== реализация интерфейса =====

func (s *teamService) CreateTeam(ctx context.Context, team domain.Team) (*domain.Team, error) {
	var createdTeam *domain.Team

	err := s.tx.Do(ctx, func(txCtx context.Context) error {
		existingTeam, err := s.teams.GetByName(txCtx, team.Name)
		if err == nil && existingTeam != nil {
			return domain.NewTeamExistsError(team.Name)
		}
		// Если ошибка не "not found", возвращаем её
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		// Создаём команду
		if err := s.teams.CreateTeam(txCtx, team); err != nil {
			return err
		}

		// Создаём/обновляем пользователей
		for _, member := range team.Members {
			// Устанавливаем team_name для каждого пользователя
			member.TeamName = team.Name
			if err := s.users.CreateOrUpdate(txCtx, member); err != nil {
				return err
			}
		}

		// Получаем созданную команду с пользователями
		createdTeam, err = s.teams.GetByName(txCtx, team.Name)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdTeam, nil
}

func (s *teamService) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	team, err := s.teams.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("team")
		}
		return nil, err
	}

	return team, nil
}
