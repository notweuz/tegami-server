package user

import "github.com/google/uuid"

type UserID struct {
	ID uuid.UUID `uri:"id" binding:"required,uuid"`
}
