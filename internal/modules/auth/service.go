package auth

import "context"

type HashService interface {
	HashPassword(ctx context.Context, password string) (string, error)
	VerifyPassword(ctx context.Context, password, encodedHash string) (bool, error)
}
