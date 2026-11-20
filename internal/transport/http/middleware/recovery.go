package middleware

import (
	"net/http"
	"tegami-server/internal/transport/http/failure"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Recovery() gin.RecoveryFunc {
	return func(c *gin.Context, recovered any) {
		log.Error().
			Interface("panic", recovered).
			Stack().
			Msg("panic recovered")

		if c.Writer.Written() {
			c.Abort()
			return
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			failure.NewResponse(
				"INTERNAL_SERVER_ERROR",
				"An internal server error has occurred",
				nil,
			),
		)
	}
}
