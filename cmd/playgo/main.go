package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/vendor116/playgo/internal"
	"github.com/vendor116/playgo/internal/api"
	"github.com/vendor116/playgo/internal/config"
	"github.com/vendor116/playgo/internal/generated"
	"github.com/vendor116/playgo/internal/http"
)

var (
	app     = "play-go"
	version = "dev"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "path to config file")
	flag.Parse()

	internal.DefaultJSONLogger(app, version)

	cfg, err := config.Load[config.App](cfgPath)
	if err != nil {
		slog.Default().Error("failed to load config", "error", err)
		return
	}

	if err = internal.SetLogLevel(cfg.LogLevel); err != nil {
		slog.Default().Warn("failed to set log level", "error", err)
	}

	apiServer := api.NewServer()

	ctx, cancel := context.WithCancelCause(context.Background())

	var wg sync.WaitGroup
	wg.Go(func() {
		err = http.StartAPIServer(
			ctx,
			generated.HandlerFromMux(apiServer, api.GetRouter(apiServer)),
			cfg.APIPServer.Host,
			cfg.APIPServer.Port,
		)
		if err != nil {
			cancel(fmt.Errorf("failed to start API server: %w", err))
		}
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		slog.Default().Warn("received shutdown signal", "signal", sig.String())
		cancel(context.Canceled)
	case <-ctx.Done():
		if err = context.Cause(ctx); !errors.Is(err, context.Canceled) {
			slog.Default().Error("application completed", "error", err)
		}
	}

	wg.Wait()
}
