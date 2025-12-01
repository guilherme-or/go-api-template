package logger

import (
	"log/slog"
)

func Load() {
	slog.SetDefault(slog.New(handler()))
}
