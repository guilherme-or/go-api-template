package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

// HTTP Server implementation with graceful shutdown and safe start
type Server struct {
	server  http.Server
	handler http.Handler
	port    int
	wg      sync.WaitGroup
}

func NewServer(port int, handler http.Handler) *Server {
	return &Server{
		handler: handler,
		port:    port,
	}
}

func (s *Server) Start() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.server = http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", s.port),
		Handler:      s.handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 0 * time.Second,
	}

	s.wg.Add(1)

	s.wg.Go(func() {
		slog.Info(fmt.Sprintf("Server started at http://%s", s.server.Addr), "address", s.server.Addr, "port", s.port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}

		slog.Info("Server stopped")
		s.wg.Done()
	})
}

func (s *Server) Stop() error {
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	slog.Info("Gracefully shutting down server", "timeout", timeout)

	if err := s.server.Shutdown(ctx); err != nil {
		if err := s.server.Close(); err != nil {
			slog.Error("Server shutdown error", "error", err)
			return err
		}
		return err
	}

	s.wg.Wait()
	return nil
}
