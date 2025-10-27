package logger

import (
	"os"

	"hackathon/internal/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Initialize(cfg *config.Config) {
	level, err := zerolog.ParseLevel(cfg.App.Log)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
