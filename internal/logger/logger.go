package logger

import (
	"context"
	"log/slog"
	"os"
)

type CtxLoggerKey struct {
}

func New(level string) *slog.Logger {
	otps := &slog.HandlerOptions{}

	switch level {
	case "DEBUG":
		otps.Level = slog.LevelDebug
	case "WARN":
		otps.Level = slog.LevelWarn
	case "ERROR":
		otps.Level = slog.LevelError
	default:
		otps.Level = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, otps))

}

func MustLoggerFromCtx(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(CtxLoggerKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
