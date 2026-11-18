package middleware

import (
	"tegami-server/internal/errs"

	"github.com/gin-gonic/gin"
)

func NoMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = c.Error(errs.New(errs.MethodNotAllowed, "METHOD_NOT_ALLOWED", "Requested method is not allowed", nil))
	}
}
