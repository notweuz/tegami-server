package auth

import (
	"context"
	"errors"
	"fmt"
	userModule "tegami-server/internal/modules/user"
	"tegami-server/internal/transport/http/auth"
	"tegami-server/internal/transport/http/user"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
	dummyHash   string
}

func NewService(hashService HashService, jwtService JWTService, userService user.Service) *service {
	dummyHash, err := hashService.HashPassword(context.Background(), "dummy-password")
	if err != nil {
		log.Error().Err(err).Msg("couldn't create dummy hash")
	}

	return &service{
		hashService: hashService,
		jwtService:  jwtService,
		userService: userService,
		dummyHash:   dummyHash,
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

	if err != nil && !errors.Is(err, userModule.ErrUserNotFound) {
		return "", fmt.Errorf("login: %w", err)
	}

	hash := s.dummyHash
	userExists := err == nil && u != nil

	if userExists {
		hash = u.Password
	}

	valid, err := s.hashService.VerifyPassword(ctx, req.Password, hash)
	if err != nil {
		return "", fmt.Errorf("login: %w", err)
	}

	if !userExists || !valid {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(ctx, u.ID)
	if err != nil {
		return "", fmt.Errorf("login: %w", err)
	}

	return token, nil
}
