package main

import (
	"tegami-server/internal/app"
	"tegami-server/internal/config"

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
	cfg, err := config.New(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if err := app.Run(cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to start app")
	}
}
