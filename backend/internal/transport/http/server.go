package httptransport

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"
	"traci/backend/internal/config"
)

type Server struct {
	server          *http.Server
	logger          *slog.Logger
	shutdownTimeout time.Duration
}

type ServerError struct {
	Host string
	Port int
	Err  error
}

func (e *ServerError) Error() string {
	return e.Err.Error()
}

func (e *ServerError) Unwrap() error {
	return e.Err
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
	host, port, err := net.SplitHostPort(server.Addr)
	if err != nil {
		return fmt.Errorf("invalid server address: %w", err)
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("invalid server port: %w", err)
	}
	serverError := func(operation string, err error) error {
		return fmt.Errorf("%s: %w", operation, &ServerError{
			Host: host,
			Port: portNumber,
			Err:  err,
		})
	}
	done := make(chan error, 1)
	go func() {
		logger.Info(
			"server starting",
			"server.address", host,
			"server.port", portNumber,
		)
		done <- server.ListenAndServe()
	}()
	select {
	case err := <-done:
		if errors.Is(err, http.ErrServerClosed) {
			logger.Info(
				"server stopped",
				"server.address", host,
				"server.port", portNumber,
			)
			return nil
		}
		return serverError("server failed", err)
	case <-ctx.Done():
		logger.Info(
			"server shutting down",
			"server.address", host,
			"server.port", portNumber,
		)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			<-done
			return serverError("server shutdown failed", err)
		}
		err := <-done
		if errors.Is(err, http.ErrServerClosed) {
			logger.Info(
				"server stopped",
				"server.address", host,
				"server.port", portNumber,
			)
			return nil
		}
		return serverError("server failed", err)
	}
}
