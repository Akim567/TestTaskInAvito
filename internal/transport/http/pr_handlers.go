package http

import (
	"net/http"

	"TestTaskInAvito/internal/service"
)

type PRHandler struct {
	svc service.PRService
}

func NewPRHandler(svc service.PRService) *PRHandler {
	return &PRHandler{svc: svc}
}

func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	// Реализуем на следующем шаге
	panic("not implemented")
}

func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	// Реализуем на следующем шаге
	panic("not implemented")
}

func (h *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	// Реализуем на следующем шаге
	panic("not implemented")
}
