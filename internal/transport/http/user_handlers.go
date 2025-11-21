package http

import (
	"encoding/json"

	"TestTaskInAvito/internal/service"
	stdhttp "net/http"
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

	resp := userToDTO(*updated)
	writeJSON(w, stdhttp.StatusOK, resp)
}

// GET /users/getReview?user_id=123
func (h *UserHandler) GetReview(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		// Явный 400, чтобы не городить странные ошибки
		writeJSON(w, stdhttp.StatusBadRequest, ErrorResponseDTO{
			Error: ErrorBodyDTO{
				Code:    "BAD_REQUEST",
				Message: "user_id is required",
			},
		})
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
