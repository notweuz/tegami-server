package http

import (
	"tegami-server/internal/config"
	"tegami-server/internal/platform/database"
	"tegami-server/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Router interface {
	Start()
	Shutdown()
}

type server struct {
	http Router
}

func New(cfg *config.Config) (*server, error) {
	db, err := database.Setup(cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed to start database!")
		return nil, err
	}

	validator.Register()
	gin.SetMode(cfg.GinMode)

	router, err := NewRouter(db, cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed to start router")
		return nil, err
	}

	router.Start()

	return &server{
		http: router,
	}, nil
}

func (a *server) Shutdown() {
	log.Info().Msg("started shutting down server")
	a.http.Shutdown()
}
