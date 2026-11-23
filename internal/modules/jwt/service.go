package jwt

import (
	"context"
	"errors"
	"fmt"
	"tegami-server/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const TokenTTL = 30 * 24 * time.Hour

type service struct {
	cfg *config.Config
}

func NewService(cfg *config.Config) *service {
	return &service{
		cfg: cfg,
	}
}

type Claims struct {
	jwt.RegisteredClaims
}

func (s *service) GenerateToken(ctx context.Context, userID uuid.UUID) (string, error) {
	now := time.Now()

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *service) ParseToken(ctx context.Context, tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(s.cfg.JWTSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return uuid.Nil, fmt.Errorf("parse token: %w: %w", ErrTokenExpired, err)
		default:
			return uuid.Nil, fmt.Errorf("parse token: %w: %w", ErrJWTInvalid, err)
		}
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return uuid.Nil, ErrInvalidTokenClaims
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse token: %w: %w", ErrJWTInvalid, err)
	}
	if userID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("parse token: %w: empty subject uuid", ErrJWTInvalid)
	}

	return userID, nil
}
