package main

import (
	"tegami-server/internal/config"
	"tegami-server/internal/platform/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	logger.Setup()
	log.Info().Msg("starting server")

	cfg, err := config.New(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	logger.UpdateLevel(cfg.LogLevel)
}