package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/kasparovgs/subscription-aggregation-service/cmd/app/config"
)

func Init(cfg *config.AppConfig) {
	lvl := slog.LevelInfo
	if cfg.LoggerConfig.Level != "" {
		switch strings.ToLower(cfg.LoggerConfig.Level) {
		case "debug":
			lvl = slog.LevelDebug
		case "warn":
			lvl = slog.LevelWarn
		case "error":
			lvl = slog.LevelError
		}
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(logger)
}
