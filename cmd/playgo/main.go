package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/vendor116/playgo/internal"
	"github.com/vendor116/playgo/internal/api"
	"github.com/vendor116/playgo/internal/config"
	"github.com/vendor116/playgo/internal/generated"
	"github.com/vendor116/playgo/internal/http"
)

var (
	name    = "playgo"
	version = "dev"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "path to config file")
	flag.Parse()

	internal.DefaultJSONLogger(name, version)

	cfg, err := config.Load[config.App](cfgPath)
	if err != nil {
		slog.Default().Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if err = internal.SetLogLevel(cfg.LogLevel); err != nil {
		slog.Default().Warn("failed to set log level", "error", err)
	}

	apiServer := api.NewServer()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err = http.RunAPIServer(
		ctx,
		generated.HandlerFromMux(apiServer, api.GetRouter(apiServer)),
		cfg.APIPServer.Host,
		cfg.APIPServer.Port,
	)
	if err != nil {
		slog.Default().Error("failed to start API server", "error", err)
	}
}
