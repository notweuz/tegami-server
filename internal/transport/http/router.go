package http

import (
	"context"
	"fmt"
	"net/http"
	"tegami-server/internal/config"
	healthModule "tegami-server/internal/modules/health"
	"tegami-server/internal/transport/http/health"
	"tegami-server/internal/transport/http/middleware"
	"time"

	_ "tegami-server/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type router struct {
	DB  *gorm.DB
	Gin *gin.Engine
	Cfg *config.Config
	Srv *http.Server
}

func NewRouter(db *gorm.DB, cfg *config.Config) (*router, error) {
	r := gin.New()
	r.Use(
		middleware.Logger(),
		gin.CustomRecovery(middleware.Recovery()),
		cors.Default(),
		middleware.Error(),
	)

	r.HandleMethodNotAllowed = true
	r.NoRoute(middleware.NoRoute())
	r.NoMethod(middleware.NoMethod())

	if cfg.ExposeSwagger {
		log.Info().Msg("expose swagger is enabling, creating swagger endpoint")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	baseDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("router db get: %w", err)
	}

	healthSVC := healthModule.NewService(baseDB)
	healthHandler := health.NewHandler(healthSVC)

	// TODO: registrar interface
	api := r.Group("/api")
	{
		api.GET("/health", healthHandler.GetInfo)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.AppPort),
		Handler: r.Handler(),
	}

	return &router{
		DB:  db,
		Gin: r,
		Cfg: cfg,
		Srv: srv,
	}, nil
}

func (a *router) Start() {
	go func() {
		log.Info().Msg("starting http router")
		if err := a.Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("failed to start http server")
		}
	}()
}

func (a *router) Shutdown() {
	log.Info().Msg("shutting down http router")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := a.Srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to gracefully shutdown http router")
		return
	}
	log.Info().Msg("http router successful graceful shutdown")
}
