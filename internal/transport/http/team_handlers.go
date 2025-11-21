package http

import (
	"encoding/json"
	"net/http"

	"TestTaskInAvito/internal/service"
)

type TeamHandler struct {
	svc service.TeamService
}

func NewTeamHandler(svc service.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

// POST /team/add
func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var dto TeamDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, err)
		return
	}

	team := teamFromDTO(dto)

	created, err := h.svc.CreateTeam(r.Context(), team)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := TeamResponseDTO{
		Team: teamToDTO(*created),
	}
	writeJSON(w, http.StatusCreated, resp)
}

// GET /team/get?team_name=...
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		// в спецификации такого кейса нет, но логично вернуть 400
		writeJSON(w, http.StatusBadRequest, ErrorResponseDTO{
			Error: ErrorBodyDTO{
				Code:    "BAD_REQUEST",
				Message: "team_name is required",
			},
		})
		return
	}

	team, err := h.svc.GetTeam(r.Context(), teamName)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := TeamResponseDTO{
		Team: teamToDTO(*team),
	}
	writeJSON(w, http.StatusOK, resp)
}
