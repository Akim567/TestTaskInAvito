package postgres

import (
	"context"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"TestTaskInAvito/internal/domain"
	"TestTaskInAvito/internal/repo"
)

type prRepo struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewPRRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) repo.PRRepository {
	return &prRepo{
		db:     db,
		getter: getter,
	}
}

// Create создаёт PR и его ревьюверов.
func (r *prRepo) Create(ctx context.Context, pr domain.PullRequest) error {
	dbOrTx := r.getter.DefaultTrOrDB(ctx, r.db)

	_, err := dbOrTx.ExecContext(ctx,
		`INSERT INTO pull_requests (id, name, author_id, status)
         VALUES ($1, $2, $3, $4)`,
		pr.ID, pr.Name, pr.Author, pr.Status,
	)
	if err != nil {
		return err
	}

	for _, reviewerID := range pr.AssignedReviewers {
		_, err := dbOrTx.ExecContext(ctx,
			`INSERT INTO pull_request_reviewers (pr_id, reviewer_id)
             VALUES ($1, $2)`,
			pr.ID, reviewerID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByID возвращает PR с ревьюверами.
func (r *prRepo) GetByID(ctx context.Context, id string) (*domain.PullRequest, error) {
	dbOrTx := r.getter.DefaultTrOrDB(ctx, r.db)

	var pr domain.PullRequest
	err := dbOrTx.GetContext(ctx, &pr,
		`SELECT id, name, author_id, status, created_at, merged_at
         FROM pull_requests
         WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	reviewers, err := r.getReviewers(ctx, dbOrTx, id)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

// Update обновляет PR и список ревьюверов.
func (r *prRepo) Update(ctx context.Context, pr domain.PullRequest) error {
	dbOrTx := r.getter.DefaultTrOrDB(ctx, r.db)

	_, err := dbOrTx.ExecContext(ctx,
		`UPDATE pull_requests
         SET name = $1,
             author_id = $2,
             status = $3,
             created_at = $4,
             merged_at = $5
         WHERE id = $6`,
		pr.Name, pr.Author, pr.Status, pr.CreatedAt, pr.MergedAt, pr.ID,
	)
	if err != nil {
		return err
	}

	// пересобираем ревьюверов
	_, err = dbOrTx.ExecContext(ctx,
		`DELETE FROM pull_request_reviewers WHERE pr_id = $1`,
		pr.ID,
	)
	if err != nil {
		return err
	}

	for _, reviewerID := range pr.AssignedReviewers {
		_, err := dbOrTx.ExecContext(ctx,
			`INSERT INTO pull_request_reviewers (pr_id, reviewer_id)
             VALUES ($1, $2)`,
			pr.ID, reviewerID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByReviewer возвращает все PR, где userID назначен ревьювером.
func (r *prRepo) GetByReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	dbOrTx := r.getter.DefaultTrOrDB(ctx, r.db)

	var prs []domain.PullRequest
	err := dbOrTx.SelectContext(ctx, &prs,
		`SELECT p.id, p.name, p.author_id, p.status, p.created_at, p.merged_at
         FROM pull_requests p
         JOIN pull_request_reviewers prr ON prr.pr_id = p.id
         WHERE prr.reviewer_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	for i := range prs {
		reviewers, err := r.getReviewers(ctx, dbOrTx, prs[i].ID)
		if err != nil {
			return nil, err
		}
		prs[i].AssignedReviewers = reviewers
	}

	return prs, nil
}

// getReviewers возвращает список reviewer_id по pr_id.
func (r *prRepo) getReviewers(ctx context.Context, dbOrTx sqlx.ExtContext, prID string) ([]string, error) {
	var reviewerIDs []string
	err := sqlx.SelectContext(ctx, dbOrTx, &reviewerIDs,
		`SELECT reviewer_id
         FROM pull_request_reviewers
         WHERE pr_id = $1`,
		prID,
	)
	if err != nil {
		return nil, err
	}

	return reviewerIDs, nil
}
