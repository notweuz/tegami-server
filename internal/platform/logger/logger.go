package logger

import (
	"os"
	"tegami-server/internal/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(cfg *config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	updateLevel(cfg.LogLevel)
}

func updateLevel(new string) {
	lvl, err := zerolog.ParseLevel(new)
	if err != nil {
		log.Error().Str("lvl", new).Msg("Failed to parse log level")
		log.Warn().Msg("Fallback to default one")
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
	log.Info().Msgf("Setting log level to %s", lvl.String())
}
