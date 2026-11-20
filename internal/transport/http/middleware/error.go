package middleware

import (
	"errors"
	"net/http"
	"tegami-server/internal/errs"
	"tegami-server/internal/transport/http/failure"
	"tegami-server/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var statusMap = map[errs.Kind]int{
	errs.Internal:         http.StatusInternalServerError,
	errs.Unauthenticated:  http.StatusUnauthorized,
	errs.MethodNotAllowed: http.StatusMethodNotAllowed,
	errs.BadRequest:       http.StatusBadRequest,
	errs.Forbidden:        http.StatusForbidden,
	errs.NotFound:         http.StatusNotFound,
	errs.Conflict:         http.StatusConflict,
	errs.BadGateway:       http.StatusBadGateway,
}

func parseStatus(kind errs.Kind) (int, bool) {
	status, ok := statusMap[kind]
	return status, ok
}

var internalServerError = errs.New(errs.Internal, "INTERNAL_SERVER_ERROR", "An internal server error has occurred", errs.ErrInternalServer)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err
		var appErr *errs.AppError
		if !errors.As(err, &appErr) {
			appErr = internalServerError
		}

		status, ok := parseStatus(appErr.Kind)
		if !ok {
			appErr = internalServerError
			status = http.StatusInternalServerError
		}

		if status >= http.StatusInternalServerError {
			log.Error().Err(err).
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Msg("request failed")
		} else if status >= http.StatusBadRequest {
			log.Warn().Err(err).
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Msg("request failed")
		}

		resp := failure.NewResponse(appErr.Code, appErr.Message, util.MapToList(appErr.Fields, failure.NewFieldError))

		c.JSON(status, resp)
	}
}
