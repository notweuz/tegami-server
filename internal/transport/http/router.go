package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type RouterRegistrar interface {
	RegisterHandlers(r *gin.RouterGroup)
}

type router struct {
	group    *gin.RouterGroup
	handlers []RouterRegistrar
}

func NewRouter(r *gin.RouterGroup, handlers ...RouterRegistrar) *router {
	return &router{
		group:    r,
		handlers: handlers,
	}
}

func (r *router) Register() {
	log.Info().Int("amount", len(r.handlers)).Msg("registering http handlers")
	for _, h := range r.handlers {
		h.RegisterHandlers(r.group)
	}
}
