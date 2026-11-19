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

// GetInfo godoc
//
//	@Summary		Get app health
//	@Description	Returns the current status of the app.
//	@Tags			status
//	@Produce		json
//	@Success		200	{object}	InfoResponse
//	@Router			/health [get]
func (h *Handler) GetInfo(c *gin.Context) {
	health := h.svc.GetInfo(c.Request.Context())

	dto := NewInfoResponse(health.Status, health.Commit, health.Uptime)
	c.JSON(http.StatusOK, dto)
}

func (h *Handler) RegisterHandlers(r *gin.RouterGroup) {
	r.GET("/health", h.GetInfo)
}