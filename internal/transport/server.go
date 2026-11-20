package transport

import (
	"context"
	"fmt"
	"net/http"
	"tegami-server/internal/config"
	healthModule "tegami-server/internal/modules/health"
	"tegami-server/internal/platform/database"
	httpRouter "tegami-server/internal/transport/http"
	"tegami-server/internal/transport/http/middleware"
	"tegami-server/internal/validator"
	"time"

	_ "tegami-server/docs"

	"tegami-server/internal/transport/http/health"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router interface {
	Register()
}

type server struct {
	httpRouter Router

	gin *gin.Engine
	srv *http.Server
}

func New(cfg *config.Config) (*server, error) {
	db, err := database.Setup(cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed to start database!")
		return nil, err
	}

	validator.Register()
	gin.SetMode(cfg.GinMode)

	baseDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("router db get: %w", err)
	}

	healthSVC := healthModule.NewService(baseDB)

	r := gin.New()

	r.HandleMethodNotAllowed = true
	r.NoRoute(middleware.NoRoute())
	r.NoMethod(middleware.NoMethod())

	if cfg.ExposeSwagger {
		log.Info().Msg("expose swagger is enabling, creating swagger endpoint")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	api := r.Group("/api/v1")
	hRouter := httpRouter.NewRouter(api,
		health.NewHandler(healthSVC),
	)

	log.Info().Msg("registering routers")
	hRouter.Register()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.AppPort),
		Handler: r.Handler(),
	}

	return &server{
		httpRouter: hRouter,
		srv:        srv,
		gin:        r,
	}, nil
}

func (s *server) Start() {
	go func() {
		log.Info().Msg("starting server")
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("failed to start server")
		}
	}()
}

func (s *server) Shutdown() {
	log.Info().Msg("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to gracefully shutdown server")
		return
	}
	log.Info().Msg("server successful graceful shutdown")
}
