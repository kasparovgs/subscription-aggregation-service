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
	Address string `env:"APP_PORT"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type AppInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type AppConfig struct {
	AppInfo `yaml:"app"`
	HTTPConfig
	LoggerConfig `yaml:"logger"`
}
