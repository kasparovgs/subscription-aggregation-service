package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kasparovgs/subscription-aggregation-service/api/http"

	"github.com/kasparovgs/subscription-aggregation-service/repository/postgres_storage"

	"github.com/kasparovgs/subscription-aggregation-service/usecases/service"

	appConfig "github.com/kasparovgs/subscription-aggregation-service/cmd/app/config"

	"github.com/kasparovgs/subscription-aggregation-service/pkg/config"
	pkgHttp "github.com/kasparovgs/subscription-aggregation-service/pkg/http"
	"github.com/kasparovgs/subscription-aggregation-service/pkg/logger"

	_ "github.com/kasparovgs/subscription-aggregation-service/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

// @title My API
// @version 1.0
// @description REST-service for aggregating data about users online subscriptions.

// @host localhost:8080
// @BasePath /
func main() {

	appFlags := appConfig.ParseFlags()
	var cfg appConfig.AppConfig
	config.MustLoad(appFlags.ConfigPath, &cfg)

	logger.Init(&cfg)

	slog.Info("starting service", "name", cfg.Name, "version", cfg.Version)

	slog.Info("config loaded", "config_path", appFlags.ConfigPath)

	connStr := os.Getenv("DB_CONN_STR")
	if connStr == "" {
		slog.Error("DB_CONN_STR environment variable is required")
		os.Exit(1)
	}

	subscriptionRepo, err := postgres_storage.NewSubscriptionDB(connStr)
	if err != nil {
		slog.Error("no connection with postgres", "error", err)
		os.Exit(1)
	}
	defer func() {
		slog.Info("closing database connection")
		if err := subscriptionRepo.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}()

	slog.Info("connected to postgres")

	subscriptionService := service.NewSubscription(subscriptionRepo)
	subscriptionHandlers := http.NewSubscriptionHandler(subscriptionService)

	r := chi.NewRouter()
	r.Use(pkgHttp.LoggingMiddleware)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	subscriptionHandlers.WithSubscriptionHandlers(r)

	server := pkgHttp.CreateServer(r, cfg.Address)
	go func() {
		slog.Info("starting HTTP server", "address", cfg.Address)
		if err := server.ListenAndServe(); err != nil {
			slog.Error("failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	slog.Info("waiting for shutdown signal...")
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	} else {
		slog.Info("server exited gracefully")
	}
}
