package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"TestTaskInAvito/internal/domain"
	"TestTaskInAvito/internal/repo"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
)

// userRepo — реализация UserRepository.
type userRepo struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

// NewUserRepository — конструктор.
func NewUserRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) repo.UserRepository {
	return &userRepo{
		db:     db,
		getter: getter,
	}
}

// SetIsActive обновляет активность пользователя.
func (r *userRepo) SetIsActive(ctx context.Context, userID string, isActive bool) error {
	dbOrTx := r.getter.DefaultTrOrDB(ctx, r.db)

	_, err := dbOrTx.ExecContext(ctx,
		`UPDATE users SET is_active = $1 WHERE id = $2`,
		isActive, userID,
	)
	return err
}

// GetByID возвращает пользователя по ID.
func (r *userRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	dbOrTx := r.getter.DefaultTrOrDB(ctx, r.db)

	var user domain.User
	err := dbOrTx.GetContext(ctx, &user,
		`SELECT id, username, team_name, is_active FROM users WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetActiveTeamMembers возвращает активных членов команды.
func (r *userRepo) GetActiveTeamMembers(ctx context.Context, teamName string) ([]domain.User, error) {
	dbOrTx := r.getter.DefaultTrOrDB(ctx, r.db)

	var members []domain.User
	err := dbOrTx.SelectContext(ctx, &members,
		`SELECT id, username, team_name, is_active FROM users WHERE team_name = $1 AND is_active = true`,
		teamName,
	)
	if err != nil {
		return nil, err
	}

	return members, nil
}

// CreateOrUpdate создаёт или обновляет пользователя.
func (r *userRepo) CreateOrUpdate(ctx context.Context, user domain.User) error {
	dbOrTx := r.getter.DefaultTrOrDB(ctx, r.db)

	_, err := dbOrTx.ExecContext(ctx,
		`INSERT INTO users (id, username, team_name, is_active) 
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (id) DO UPDATE SET 
		 username = EXCLUDED.username,
		 team_name = EXCLUDED.team_name,
		 is_active = EXCLUDED.is_active`,
		user.ID, user.Username, user.TeamName, user.IsActive,
	)
	return err
}
