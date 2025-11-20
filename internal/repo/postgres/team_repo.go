package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"TestTaskInAvito/internal/domain"
	"TestTaskInAvito/internal/repo"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
)

// teamRepo — реализация TeamRepository.
type teamRepo struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

// NewTeamRepository — конструктор.
func NewTeamRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) repo.TeamRepository {
	return &teamRepo{
		db:     db,
		getter: getter,
	}
}

// CreateTeam сохраняет новую команду.
// domain.Team содержит список членов, но сама таблица teams хранит только name.
func (r *teamRepo) CreateTeam(ctx context.Context, team domain.Team) error {
	dbOrTx := r.getter.DefaultTrOrDB(ctx, r.db)

	_, err := dbOrTx.ExecContext(ctx,
		`INSERT INTO teams (name) VALUES ($1)`,
		team.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetByName возвращает команду + её членов.
func (r *teamRepo) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	dbOrTx := r.getter.DefaultTrOrDB(ctx, r.db)

	// Проверяем, что сама команда существует
	var teamName string
	err := dbOrTx.GetContext(ctx, &teamName,
		`SELECT name FROM teams WHERE name = $1`,
		name,
	)
	if err != nil {
		return nil, err
	}

	// Забираем пользователей команды
	var members []domain.User
	err = dbOrTx.SelectContext(ctx, &members,
		`SELECT id, username, team_name, is_active FROM users WHERE team_name = $1`,
		name,
	)
	if err != nil {
		return nil, err
	}

	return &domain.Team{
		Name:    teamName,
		Members: members,
	}, nil
}
