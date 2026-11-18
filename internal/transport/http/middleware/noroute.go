package middleware

import (
	"tegami-server/internal/errs"

	"github.com/gin-gonic/gin"
)

func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = c.Error(errs.New(errs.NotFound, "ROUTE_NOT_FOUND", "route not found", nil))
	}
}
