package domain

import "time"

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID                string     `db:"id"`
	Name              string     `db:"name"`
	Author            string     `db:"author_id"` // user_id автора
	Status            PRStatus   `db:"status"`
	AssignedReviewers []string   `db:"-"` // это отдельная таблица
	CreatedAt         time.Time  `db:"created_at"`
	MergedAt          *time.Time `db:"merged_at"`
}
