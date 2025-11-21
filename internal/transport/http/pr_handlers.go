package http

import (
	"encoding/json"

	"TestTaskInAvito/internal/domain"
	"TestTaskInAvito/internal/service"
	stdhttp "net/http"
)

type PRHandler struct {
	svc service.PRService
}

func NewPRHandler(svc service.PRService) *PRHandler {
	return &PRHandler{svc: svc}
}

// POST /pullRequest/create
func (h *PRHandler) CreatePR(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req CreatePRRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}

	// Входная доменная модель: сервис сам проставит статус, даты, ревьюверов
	in := domain.PullRequest{
		ID:     req.PullRequestID,
		Name:   req.PullRequestName,
		Author: req.AuthorID,
	}

	created, err := h.svc.CreatePR(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := PRResponseDTO{
		PR: prToDTO(*created),
	}
	writeJSON(w, stdhttp.StatusCreated, resp)
}

// POST /pullRequest/merge
func (h *PRHandler) MergePR(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req MergePRRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}

	pr, err := h.svc.MergePR(r.Context(), req.PullRequestID)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := PRResponseDTO{
		PR: prToDTO(*pr),
	}
	writeJSON(w, stdhttp.StatusOK, resp)
}

// POST /pullRequest/reassign
func (h *PRHandler) ReassignReviewer(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req ReassignRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}

	pr, replacedBy, err := h.svc.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := ReassignResponseDTO{
		PR:         prToDTO(*pr),
		ReplacedBy: replacedBy,
	}
	writeJSON(w, stdhttp.StatusOK, resp)
}
