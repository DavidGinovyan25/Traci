package httptransport

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"traci/backend/internal/config"
)

type Server struct {
	server          *http.Server
	logger          *slog.Logger
	shutdownTimeout time.Duration
}

func NewServer(cfg config.Config, handler http.Handler, logger *slog.Logger) *Server {
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	return &Server{server: server, logger: logger, shutdownTimeout: cfg.ShutdownTimeout}
}

func (s *Server) Run(ctx context.Context) error {
	server, logger := s.server, s.logger

	done := make(chan error, 1)
	go func() {
		logger.Info("server starting", "addr", server.Addr)
		done <- server.ListenAndServe()
	}()
	select {
	case err := <-done:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		logger.Info("server shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return err
		}
		err := <-done
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
