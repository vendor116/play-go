package internal

import (
	"log/slog"
	"os"
)

var logLevel slog.LevelVar

func DefaultJSONLogger(app, version string) {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: &logLevel,
	})

	slog.SetDefault(slog.New(h).With(
		slog.String("app", app),
		slog.String("version", version),
	))
}

func SetLogLevel(level string) error {
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		return err
	}

	return nil
}
