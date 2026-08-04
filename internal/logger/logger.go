package logger

import (
	"fmt"
	"log/slog"
	"os"

	"REST_Server/internal/errors"
)

func New(logLevel string) (*slog.Logger, error) {
	var lvl slog.Level

	switch logLevel {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lg := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(lg)
		return lg, fmt.Errorf(
			"%w: unknown log level: %s, set default",
			errors.ErrInvalidLogLevel,
			logLevel,
		)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	return slog.New(handler), nil
}
