package http

import (
	"net/http"

	"TestTaskInAvito/internal/service"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	// Реализуем на следующем шаге
	panic("not implemented")
}

func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	// Реализуем на следующем шаге
	panic("not implemented")
}
