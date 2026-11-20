package user

import (
	"context"
	"errors"
	"net/http"
	"tegami-server/internal/errs"
	"tegami-server/internal/modules/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrBadRequest = errs.New(errs.BadRequest, "BAD_REQUEST", "Bad request", nil)
)

type Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
	Create(ctx context.Context, users ...*user.User) error
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// GetByID
//
//	@Summary		Get the user by ID
//	@Description	Returns the user by their ID.
//	@Tags			user
//	@Produce		json
//	@Param			id	path		string			true	"User ID"
//	@Success		200	{object}	User			"OK"
//	@Failure		400	{object}	failure.Response	"No ID provided"
//	@Failure		404	{object}	failure.Response	"User not found"
//	@Failure		500	{object}	failure.Response	"Internal Server Error"
//	@Router			/users/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	var path UserID
	if err := c.ShouldBindUri(&path); err != nil {
		_ = c.Error(ErrBadRequest) // TODO: proper validation error binding
		return
	}

	user, err := h.svc.GetByID(c.Request.Context(), path.ID)
	if err != nil {
		_ = c.Error(parseError(err))
	}

	c.JSON(http.StatusOK, UserFromModel(*user))
}

// TODO: secure it under jwt auth
func (h *Handler) RegisterHandlers(r *gin.RouterGroup) {
	u := r.Group("/users")
	{
		u.GET("/:id", h.GetByID)
	}
}

func parseError(err error) *errs.AppError {
	switch {
	case errors.Is(err, user.ErrUserNotFound):
		return errs.New(errs.NotFound, "USER_NOT_FOUND", "User with this id not found", err)
	case errors.Is(err, user.ErrUserAlreadyExist):
		return errs.New(errs.Conflict, "USER_ALREADY_EXIST", "User with same username already exist", err)
	default:
		return errs.New(errs.Internal, "INTERNAL_ERROR", "Internal server error", err)
	}
}
