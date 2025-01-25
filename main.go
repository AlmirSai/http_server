package main

import (
	"log"
	"net/http"

	"github.com/almirsai/http_server/api/handlers"
	"github.com/almirsai/http_server/config"
	"github.com/almirsai/http_server/internal/server"
)

func main() {
	cfg := config.NewConfig()

	http.HandleFunc("/", handlers.HomeHandler)

	srv := server.NewServer(cfg.Port)
	if err := srv.Run(cfg.CertFile, cfg.KeyFile); err != nil {
		log.Fatal("Server error: ", err)
	}
}
