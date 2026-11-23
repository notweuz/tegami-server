package auth

import (
	"context"
	"fmt"
	userModule "tegami-server/internal/modules/user"
	"tegami-server/internal/transport/http/auth"
	"tegami-server/internal/transport/http/user"

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

type service struct {
	hashService HashService
	jwtService  JWTService
	userService user.Service
}

func NewService(hashService HashService, jwtService JWTService) *service {
	return &service{
		hashService: hashService,
		jwtService:  jwtService,
	}
}

func (s *service) Register(ctx context.Context, req auth.Credentials) (string, error) {
	encryptedPassword, err := s.hashService.HashPassword(ctx, req.Password)
	if err != nil {
		return "", fmt.Errorf("register: %w", err)
	}
	u := &userModule.User{Username: req.Username, Password: encryptedPassword}
	if err := s.userService.Create(ctx, u); err != nil {
		return "", fmt.Errorf("register: %w", err)
	}
	token, err := s.jwtService.GenerateToken(ctx, u.ID)
	if err != nil {
		return "", fmt.Errorf("register: %w", err)
	}
	return token, nil
}

func (s *service) Login(ctx context.Context, req auth.Credentials) (string, error) {
	u, err := s.userService.GetByUsername(ctx, req.Username)
	if err != nil {
		return "", fmt.Errorf("login: %w", err)
	}

	if ok, err := s.hashService.VerifyPassword(ctx, req.Password, u.Password); !ok || err != nil {
		return "", fmt.Errorf("login: %w", err)
	}

	token, err := s.jwtService.GenerateToken(ctx, u.ID)
	if err != nil {
		return "", fmt.Errorf("login: %w", err)
	}
	return token, nil
}
