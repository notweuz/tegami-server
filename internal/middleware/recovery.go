package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Recovery() gin.RecoveryFunc {
	return func(c *gin.Context, recovered any) {
		cause := fmt.Errorf("panic: %v", recovered)
		log.Error().Err(cause).Msg("panic recovered")

		c.Abort()
		_ = c.Error(cause)
	}
}