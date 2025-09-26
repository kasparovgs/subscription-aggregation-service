package config

import (
	"log"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/kasparovgs/subscription-aggregation-service/cmd/app/config"
)

func MustLoad(cfgPath string, cfg *config.AppConfig) {
	if cfgPath == "" {
		log.Fatal("Config path is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist by this path: %s", cfgPath)
	}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	if !strings.HasPrefix(cfg.HTTPConfig.Address, ":") {
		cfg.HTTPConfig.Address = ":" + cfg.HTTPConfig.Address
	}
}
