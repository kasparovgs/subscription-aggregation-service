package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

func MustLoad(cfgPath string, cfg any) {
	if cfgPath == "" {
		log.Fatal("Config path is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist by this path: %s", cfgPath)
	}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatalf("error reading config: %s", err)
	}
}
