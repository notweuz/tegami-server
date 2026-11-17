package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tegami-server/internal/config"
	"tegami-server/internal/platform/logger"
	"tegami-server/internal/transport/http"

	"github.com/rs/zerolog/log"
)

type Server interface {
	Shutdown()
}

func Run(cfg *config.Config) error {
	logger.Setup(cfg)
	log.Info().Msg("starting server")

	srv, err := http.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Debug().Msg("waiting for shutdown signal")
	<-ctx.Done()

	srv.Shutdown()
	return nil
}