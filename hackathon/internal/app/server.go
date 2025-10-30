package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hackathon/internal/config"
	"hackathon/internal/infrastructure/database"
	"hackathon/internal/infrastructure/logger"
	"hackathon/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Server struct {
	cfg    *config.Config
	engine *gin.Engine
	http   *http.Server
}

func NewServer() *Server {
	logger.Initialize(config.Cfg)

	db := database.Init(config.Cfg.Database)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middlewares.CORSMiddleware())
	engine.Use(middlewares.PanicRecoveryMiddleware())
	engine.Use(middlewares.RequestIDMiddleware())
	engine.Use(middlewares.LoggerMiddleware())

	registerRoutes(engine, db)

	srv := &http.Server{
		Addr:    config.Cfg.App.Port,
		Handler: engine,
	}

	return &Server{cfg: config.Cfg, engine: engine, http: srv}
}

func (s *Server) Start() {
	go func() {
		log.Info().Str("port", s.cfg.App.Port).Msg("server started")
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.http.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
	} else {
		log.Info().Msg("server exited")
	}
}
