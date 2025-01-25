package server

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: fmt.Sprintf(":%s", port),
		},
	}
}

func (s *Server) Run(certFile, keyFile string) error {
	log.Printf("Starting HTTPS server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServeTLS(certFile, keyFile)
}
