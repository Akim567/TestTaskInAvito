package http

import (
	"TestTaskInAvito/internal/domain"
	"time"
)

// ===== DTO из OpenAPI =====

type TeamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamDTO struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PullRequestDTO struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

type PullRequestShortDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

// Ответ для /users/getReview
type UserReviewsResponseDTO struct {
	UserID       string                `json:"user_id"`
	PullRequests []PullRequestShortDTO `json:"pull_requests"`
}

// Обёртки ответов, как в OpenAPI

type TeamResponseDTO struct {
	Team TeamDTO `json:"team"`
}

type PRResponseDTO struct {
	PR PullRequestDTO `json:"pr"`
}

type ReassignResponseDTO struct {
	PR         PullRequestDTO `json:"pr"`
	ReplacedBy string         `json:"replaced_by"`
}

// ===== ErrorResponse DTO =====

type ErrorBodyDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponseDTO struct {
	Error ErrorBodyDTO `json:"error"`
}

// ===== DTO запросов =====

type SetIsActiveRequestDTO struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type CreatePRRequestDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePRRequestDTO struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignRequestDTO struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

// ===== маппинг domain ⇄ DTO =====

func teamToDTO(t domain.Team) TeamDTO {
	members := make([]TeamMemberDTO, 0, len(t.Members))
	for _, m := range t.Members {
		members = append(members, TeamMemberDTO{
			UserID:   m.ID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}
	return TeamDTO{
		TeamName: t.Name,
		Members:  members,
	}
}

func teamFromDTO(dto TeamDTO) domain.Team {
	members := make([]domain.User, 0, len(dto.Members))
	for _, m := range dto.Members {
		members = append(members, domain.User{
			ID:       m.UserID,
			Username: m.Username,
			TeamName: dto.TeamName,
			IsActive: m.IsActive,
		})
	}
	return domain.Team{
		Name:    dto.TeamName,
		Members: members,
	}
}

func userToDTO(u domain.User) UserDTO {
	return UserDTO{
		UserID:   u.ID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func prToDTO(pr domain.PullRequest) PullRequestDTO {
	return PullRequestDTO{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.Author,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         &pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func prToShortDTO(pr domain.PullRequest) PullRequestShortDTO {
	return PullRequestShortDTO{
		PullRequestID:   pr.ID,
		PullRequestName: pr.Name,
		AuthorID:        pr.Author,
		Status:          string(pr.Status),
	}
}
