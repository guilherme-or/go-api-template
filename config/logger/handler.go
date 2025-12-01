package logger

import (
	"log/slog"
	"os"

	"github.com/guilherme-or/go-api-template/config/env"
	phuslog "github.com/phuslu/log"
)

func logLvl() slog.Level {
	if env.G.Debug {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func handler() slog.Handler {
	return phuslog.SlogNewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: env.G.Debug,
		Level:     logLvl(),
	})
}
