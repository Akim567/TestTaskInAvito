package http

import (
	"net/http"

	"TestTaskInAvito/internal/service"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	teamSvc service.TeamService,
	userSvc service.UserService,
	prSvc service.PRService,
) http.Handler {
	r := chi.NewRouter()

	teamHandler := NewTeamHandler(teamSvc)
	userHandler := NewUserHandler(userSvc)
	prHandler := NewPRHandler(prSvc)

	// --- Teams ---
	r.Post("/team/add", teamHandler.CreateTeam)
	r.Get("/team/get", teamHandler.GetTeam)

	// --- Users ---
	r.Post("/users/setIsActive", userHandler.SetIsActive)
	r.Get("/users/getReview", userHandler.GetReview)

	// --- Pull Requests ---
	r.Post("/pullRequest/create", prHandler.CreatePR)
	r.Post("/pullRequest/merge", prHandler.MergePR)
	r.Post("/pullRequest/reassign", prHandler.ReassignReviewer)

	return r
}
