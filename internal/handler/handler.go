package handler

import "golang/stockLkBack/internal/service"

const (
	userRoleKey = "role"
)

type Handler struct {
	Services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{Services: services}
}
