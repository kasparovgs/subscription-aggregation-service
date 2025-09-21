package main

import (
	"log"
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
	appFlags := appConfig.ParseFlags()
	var cfg appConfig.AppConfig
	config.MustLoad(appFlags.ConfigPath, &cfg)

	connStr := os.Getenv("DB_CONN_STR")
	if connStr == "" {
		log.Fatal("DB_CONN_STR environment variable is required")
	}

	subscriptionRepo, err := postgres_storage.NewSubscriptionDB(connStr)
	if err != nil {
		log.Fatalf("no connection with postgres: %v", err)
	}
	subscriptionService := service.NewSubscription(subscriptionRepo)
	subscriptionHandlers := http.NewSubscriptionHandler(subscriptionService)

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	subscriptionHandlers.WithSubscriptionHandlers(r)

	log.Printf("Starting server on %s", cfg.Address)
	if err := pkgHttp.CreateAndRunServer(r, cfg.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
