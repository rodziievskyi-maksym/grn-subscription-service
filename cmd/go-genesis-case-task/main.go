package main

import (
	"context"
	"errors"
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
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/github"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/repository"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/server"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/usecase"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/worker"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/pkg/databases/postgres"
)

func main() {
	initCtx, initCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer initCancel()

	validation := validator.New()

	if err := config.NewConfig(validation); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := postgres.NewPostgreClient(initCtx, config.Cfg().PostgresDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Close()

	ghClient := github.NewClient(config.Cfg().GitHubToken)
	subRepo := repository.NewSubscriptionRepository(client)
	subUseCase := usecase.NewSubscriptionUseCase(subRepo, ghClient)
	subHandler := handler.NewSubscriptionHandler(subUseCase, validation)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	scanner, err := worker.NewScanner(subRepo, ghClient)
	if err != nil {
		log.Fatalf("Failed to create scanner: %v", err)
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

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Application successfully stopped")
}
