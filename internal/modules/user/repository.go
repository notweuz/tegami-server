package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get user by id %d: %w", id, err)
		}
		return nil, fmt.Errorf("get user by id %d: %w", id, err)
	}
	return &user, nil
}

func (r *repository) Create(ctx context.Context, users ...*User) error {
	if len(users) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).Create(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("create user: %w", ErrUserAlreadyExist)
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}
