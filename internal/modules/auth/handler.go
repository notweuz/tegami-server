package auth

import (
	"context"
	"tegami-server/internal/transport/http/auth"
)

type Service interface {
	Register(ctx context.Context, req auth.Credentials) (string, error)
	Login(ctx context.Context, req auth.Credentials) (string, error)
}