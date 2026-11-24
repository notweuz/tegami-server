package auth

import (
	"context"
	"errors"
	"net/http"
	"tegami-server/internal/errs"
	"tegami-server/internal/modules/user"

	"github.com/gin-gonic/gin"
)

var (
	ErrBadRequest = errs.New(errs.BadRequest, "BAD_REQUEST", "Bad request", nil)
)

type Service interface {
	Register(ctx context.Context, req Credentials) (string, error)
	Login(ctx context.Context, req Credentials) (string, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req Credentials
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(ErrBadRequest) // TODO: proper validation error binding
		return
	}
	token, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(parseError(err))
		return
	}
	c.JSON(http.StatusCreated, NewResponse(token))
}

func (h *Handler) Login(c *gin.Context) {
	var req Credentials
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(ErrBadRequest) // TODO: proper validation error binding
		return
	}
	token, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(parseError(err))
		return
	}
	c.JSON(http.StatusOK, NewResponse(token))
}

func (h *Handler) RegisterHandlers(r *gin.RouterGroup) {
	a := r.Group("/auth")
	{
		a.POST("/register", h.Register)
		a.POST("/login", h.Login)
	}
}

// TODO: add auth errors to parser
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
