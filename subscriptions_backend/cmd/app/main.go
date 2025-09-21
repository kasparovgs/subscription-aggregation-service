package main

import (
	"log/slog"
	"os"
	"pkg/config"
	pkgHttp "pkg/http"
	"subscriptions_backend/api/http"
	appConfig "subscriptions_backend/cmd/app/config"
	"subscriptions_backend/repository/postgres_storage"
	"subscriptions_backend/usecases/service"

	_ "subscriptions_backend/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

// @title My API
// @version 1.0
// @description This is a sample server.

// @host localhost:8080
// @BasePath /
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	//ctx := context.Background()

	appFlags := appConfig.ParseFlags()
	var cfg appConfig.AppConfig
	config.MustLoad(appFlags.ConfigPath, &cfg)

	logger.Info("config loaded", "config_path", appFlags.ConfigPath)

	connStr := os.Getenv("DB_CONN_STR")
	if connStr == "" {
		logger.Error("DB_CONN_STR environment variable is required")
		os.Exit(1)
	}

	subscriptionRepo, err := postgres_storage.NewSubscriptionDB(connStr)
	if err != nil {
		logger.Error("no connection with postgres", "error", err)
		os.Exit(1)
	}
	logger.Info("connected to postgres")

	subscriptionService := service.NewSubscription(subscriptionRepo)
	subscriptionHandlers := http.NewSubscriptionHandler(subscriptionService)

	r := chi.NewRouter()
	r.Use(pkgHttp.LoggingMiddleware)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	subscriptionHandlers.WithSubscriptionHandlers(r)

	logger.Info("starting HTTP server", "address", cfg.Address)
	if err := pkgHttp.CreateAndRunServer(r, cfg.Address); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
