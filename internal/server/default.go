package server

import (
	"context"
	"fmt"
	"net/http"
	"tegami-server/internal/config"
	"tegami-server/internal/middleware"
	"tegami-server/internal/platform/database"
	"tegami-server/internal/validator"
	"time"

	_ "tegami-server/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type server struct {
	DB  *gorm.DB
	Gin *gin.Engine
	Cfg *config.Config
	Srv *http.Server
}

func New(cfg *config.Config) (Server, error) {
	db, err := database.Setup(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start database!")
		return nil, err
	}

	validator.Register()
	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(gin.CustomRecovery(middleware.Recovery()))
	r.Use(middleware.Logger())
	r.Use(cors.Default())
	// r.Use(middleware.Error()) // TODO: error handler

	r.HandleMethodNotAllowed = true
	r.NoRoute(middleware.NoRoute())
	r.NoMethod(middleware.NoMethod())

	if cfg.ExposeSwagger {
		log.Info().Msg("expose swagger is enabling, creating swagger endpoint")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.AppPort),
		Handler: r.Handler(),
	}

	return &server{
		DB:  db,
		Gin: r,
		Cfg: cfg,
		Srv: srv,
	}, nil
}

func (a *server) Start() {
	go func() {
		log.Info().Msg("starting server")
		if err := a.Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("failed to start server")
		}
	}()
}

func (a *server) Shutdown() {
	log.Info().Msg("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := a.Srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to graceful shutdown")
	} else {
		log.Info().Msg("server gracefully stopped")
	}
}
