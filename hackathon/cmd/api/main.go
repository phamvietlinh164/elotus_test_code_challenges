package main

import (
	"hackathon/internal/app"
	"hackathon/internal/config"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	server := app.NewServer(cfg)
	server.Start()
}
