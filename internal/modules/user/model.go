package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username string
	Password string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "users"
}
