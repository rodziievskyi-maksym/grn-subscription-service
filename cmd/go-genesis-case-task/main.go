package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/config"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/delivery/handler"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/cache"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/github"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/repository"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/server"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/usecase"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/worker"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/pkg/databases/postgres"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/pkg/databases/redis"
)

// @title           GitHub Release Notifier API
// @version         1.0
// @description     API Server for monitoring GitHub releases and notifying users.
// @BasePath        /
func main() {
	if err := run(); err != nil {
		slog.Error("Application stopped with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	initCtx, initCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer initCancel()

	validation := validator.New()

	if err := config.NewConfig(validation); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := postgres.NewPostgreClient(initCtx, config.Cfg().PostgresDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer client.Close()

	redisClient, err := redis.NewRedisClient(initCtx, config.GetRedisAddress(), config.Cfg().RedisPass)
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	defer redisClient.Close()

	tagCache := cache.NewTagCache(redisClient.Client, config.Cfg().RedisCacheTTL)

	baseGhClient := github.NewClient(config.Cfg().GitHubToken)
	cachedGhClient := github.NewCachedGitHubProvider(baseGhClient, tagCache)

	subRepo := repository.NewSubscriptionRepository(client)
	subUseCase := usecase.NewSubscriptionUseCase(subRepo, cachedGhClient)
	subHandler := handler.NewSubscriptionHandler(subUseCase, validation)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	scanner, err := worker.NewScanner(subRepo, cachedGhClient)
	if err != nil {
		return fmt.Errorf("failed to create scanner: %w", err)
	}

	go func() {
		if err := scanner.Run(appCtx); err != nil {
			slog.Error("Scanner stopped with error", "error", err)
		}
	}()

	srv := server.NewHTTPServer(subHandler)

	go func() {
		slog.Info("Starting HTTP server", "port", config.Cfg().Port)

		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server Run error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	slog.Info("Shutting down server...")

	// provide signal to scanner
	appCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// gracefully shutdown the server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Application successfully stopped")

	return nil
}
