package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, users ...*User) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{
		repository: repository,
	}
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("uuid", id.String()).Msg("user service get by id")
		return nil, fmt.Errorf("user service get by id: %w", err)
	}
	log.Info().Str("id", id.String()).Msg("user service get by id found user")
	return user, nil
}

func (s *service) GetByUsername(ctx context.Context, username string) (*User, error) {
	user, err := s.repository.GetByUsername(ctx, username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("user service get by username")
		return nil, fmt.Errorf("user service get by username: %w", err)
	}
	log.Info().Str("id", user.ID.String()).Msg("user service get by username found user")
	return user, nil
}

func (s *service) Create(ctx context.Context, users ...*User) error {
	if err := s.repository.Create(ctx, users...); err != nil {
		log.Error().Err(err).Msg("create users failed")
		return fmt.Errorf("user service create: %w", err)
	}
	log.Info().Int("amount", len(users)).Msg("user service successfully created user(s)")
	return nil
}