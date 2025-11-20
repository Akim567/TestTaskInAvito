package domain

import "time"

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID                string
	Name              string
	Author            string
	Status            PRStatus
	AssignedReviewers []string // <= 2
	CreatedAt         time.Time
	MergedAt          *time.Time
}
