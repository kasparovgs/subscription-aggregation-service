package main

import (
	"log"
	"os"
	"pkg/config"
	appConfig "subscriptions_backend/cmd/app/config"
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

}
