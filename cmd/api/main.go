package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/WillRoque/ai-video-descriptor/internal/ai"
	"github.com/WillRoque/ai-video-descriptor/internal/api"
	"github.com/WillRoque/ai-video-descriptor/internal/logger"
	yt "github.com/WillRoque/ai-video-descriptor/internal/youtube"
	"github.com/kkdai/youtube/v2"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

func main() {
	cfg := readConfig()
	aiVision, ytDownloader := initServices(cfg)

	r := api.NewRouter(aiVision, ytDownloader)
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	closed := make(chan struct{})
	go handleInterruption(s, closed)

	log.Info().Msgf("Starting server %v", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server startup failure")
	}

	<-closed
	log.Info().Msgf("Server shutdown successfully")
}

func handleInterruption(s *http.Server, closed chan struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	log.Info().Msgf("Shutting down server %v", s.Addr)

	if err := s.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("Server shutdown failure")
	}

	// If there were anything else that needs to be gracefully closed,
	// like database connections, file handles, etc., it would be added here.

	close(closed)
}

func initServices(cfg Config) (*ai.Vision, *yt.Downloader) {
	logger.New(cfg.Server.Debug)
	aiVision := ai.NewVision(openai.NewClient(cfg.OpenAI.ApiKey))
	ytDownloader := yt.NewDownloader(&youtube.Client{})
	return &aiVision, &ytDownloader
}
