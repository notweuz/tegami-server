package user

import (
	"tegami-server/internal/modules/user"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"uuid"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UserFromModel(model user.User) User {
	return User{
		ID:        model.ID,
		Username:  model.Username,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
