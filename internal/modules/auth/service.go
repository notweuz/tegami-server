package auth

import (
	"context"

	"github.com/google/uuid"
)

type HashService interface {
	HashPassword(ctx context.Context, password string) (string, error)
	VerifyPassword(ctx context.Context, password, encodedHash string) (bool, error)
}

type JWTService interface {
	GenerateToken(ctx context.Context, userID uuid.UUID) (string, error)
	ParseToken(ctx context.Context, tokenString string) (uuid.UUID, error)
}
