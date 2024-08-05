package main

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/saur4ig/file-storage/internal/config"
	"github.com/saur4ig/file-storage/internal/rest"
)

// @title           File storage app
// @version         0.1

// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	zerolog.TimeFieldFormat = time.DateTime

	log.Info().Msg("Starting app")
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("could not load config: %v", err)
	}

	rest.CreateServer(*conf)
}
