package health

import (
	"context"
	"net/http"
	"tegami-server/internal/modules/health"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetInfo(ctx context.Context) *health.Health
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) GetInfo(c *gin.Context) {
	health := h.svc.GetInfo(c.Request.Context())

	dto := NewInfoResponse(health.Status, health.Commit, health.Uptime)
	c.JSON(http.StatusOK, dto)
}
