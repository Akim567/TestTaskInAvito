package http

import (
	"encoding/json"
	stdhttp "net/http"

	"TestTaskInAvito/internal/service"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// POST /users/setIsActive
func (h *UserHandler) SetIsActive(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req SetIsActiveRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}

	updated, err := h.svc.SetIsActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := UserDTO{
		UserID:   updated.ID,
		Username: updated.Username,
		TeamName: updated.TeamName,
		IsActive: updated.IsActive,
	}

	writeJSON(w, stdhttp.StatusOK, resp)
}

// GET /users/getReview?user_id=123
func (h *UserHandler) GetReview(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, stdhttp.ErrNotSupported) // лёгкий вариант -> можно сделать "BAD_REQUEST"
		return
	}

	prs, err := h.svc.GetUserReviews(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	out := make([]PullRequestShortDTO, 0, len(prs))
	for _, pr := range prs {
		out = append(out, prToShortDTO(pr))
	}

	resp := UserReviewsResponseDTO{
		UserID:       userID,
		PullRequests: out,
	}

	writeJSON(w, stdhttp.StatusOK, resp)
}
