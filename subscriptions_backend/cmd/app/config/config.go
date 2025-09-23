package config

import "flag"

type AppFlags struct {
	ConfigPath string
}

func ParseFlags() AppFlags {
	configPath := flag.String("config", "", "Path to config")
	flag.Parse()
	return AppFlags{
		ConfigPath: *configPath,
	}
}

type HTTPConfig struct {
	Address string `yaml:"address"`
}

type AppConfig struct {
	HTTPConfig `yaml:"http"`
}
