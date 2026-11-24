package auth

import (
	"context"
)

type Service interface {
	Register(ctx context.Context, req Credentials) (string, error)
	Login(ctx context.Context, req Credentials) (string, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}