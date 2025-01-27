package server

import (
	"context"
	"fmt"
	"net/http"

	"http_server/auth-service/internal/config"
	"http_server/auth-service/internal/handler"
	"http_server/auth-service/pkg/middleware"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) *Server {
	router := NewRouter(authHandler, authMiddleware)

	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			Handler:      router,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
