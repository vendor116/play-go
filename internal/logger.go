package internal

import (
	"log/slog"
	"os"
)

var logLevel = &slog.LevelVar{}

func DefaultJSONLogger(name, version string) {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	slog.SetDefault(slog.New(h).With(
		slog.String("name", name),
		slog.String("version", version),
	))
}

func SetLogLevel(level string) error {
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		return err
	}

	return nil
}
