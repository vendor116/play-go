package http

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	readHeaderTimeout = 2 * time.Second
	handleTimeout     = 2 * time.Second
	shutdownTimeout   = 2 * time.Second
)

func StartAPIServer(ctx context.Context, handler http.Handler, host, port string) error {
	server := &http.Server{
		Addr:              net.JoinHostPort(host, port),
		Handler:           http.TimeoutHandler(handler, handleTimeout, http.ErrHandlerTimeout.Error()),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	logger := slog.Default().With("addr", server.Addr)

	g, gCtx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		logger.Info("starting API server")

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		logger.Info("API server shutdown gracefully")
		return nil
	})

	g.Go(func() error {
		select {
		case <-gCtx.Done():
			return gCtx.Err()
		case <-ctx.Done():
			logger.Warn("shutting down API server")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				logger.ErrorContext(shutdownCtx, "failed to shutdown API server", "error", err)
			}

			return nil
		}
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
