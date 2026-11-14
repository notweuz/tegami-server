package app

import (
	"context"
	"fmt"
	"net/http"
	"tegami-server/internal/config"
	"tegami-server/internal/platform/database"
	"tegami-server/internal/platform/validation"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type app struct {
	DB  *gorm.DB
	Gin *gin.Engine
	Cfg *config.Config
	Srv *http.Server
}

func New(cfg *config.Config) (App, error) {
	db, err := database.Setup(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start database!")
		return nil, err
	}

	validation.Register()
	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(cors.Default())
	r.HandleMethodNotAllowed = true

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.AppPort),
		Handler: r.Handler(),

	}

	return &app{
		DB:  db,
		Gin: r,
		Cfg: cfg,
		Srv: srv,
	}, nil
}

func (a *app) Start() {
	go func() {
		log.Info().Msg("starting server")
		if err := a.Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("failed to start server")
		}
	}()
}

func (a *app) Shutdown() {
	log.Info().Msg("shutting down application")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := a.Srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to graceful shutdown")
	} else {
		log.Info().Msg("server gracefully stopped")
	}
}