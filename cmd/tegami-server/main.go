package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"tegami-server/internal/config"
	"tegami-server/internal/platform/app"
	"tegami-server/internal/platform/logger"

	"github.com/rs/zerolog/log"
)

//	@title			Tegami API
//	@version		0.1
//	@description	API documentation for Tegami Server.

//	@BasePath	/api

//	@securityDefinitions.apiKey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter "Bearer <token>" (token from /auth/login).

func main() {
	logger.Setup()
	log.Info().Msg("starting server")

	cfg, err := config.New(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	logger.UpdateLevel(cfg.LogLevel)

	app, err := app.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create app")
	}
	app.Start()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Debug().Msg("waiting for shutdown signal")
	<-ctx.Done()

	app.Shutdown()
}
