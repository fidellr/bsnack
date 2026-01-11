package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	once   sync.Once
	logger *slog.Logger
)

func Init(isProduction bool) {
	once.Do(func() {
		var handler slog.Handler

		if isProduction {
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
		} else {
			handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			})
		}

		logger = slog.New(handler)
		slog.SetDefault(logger)
	})
}

func Info(msg string, args ...any)  { slog.Info(msg, args...) }
func Error(msg string, args ...any) { slog.Error(msg, args...) }
func Debug(msg string, args ...any) { slog.Debug(msg, args...) }
func Warn(msg string, args ...any)  { slog.Warn(msg, args...) }
