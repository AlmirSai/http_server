package server

import (
	"log"
	"net/http"

	"http_server/auth-service/internal/handler"
	"http_server/auth-service/pkg/middleware"
	handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRouter(authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) *mux.Router {
	r := mux.NewRouter()

	// Add logging middleware
	r.Use(func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(log.Writer(), next)
	})

	// Add CORS middleware
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.ExposedHeaders([]string{"Content-Length"}),
		handlers.AllowCredentials(),
	)
	r.Use(corsMiddleware)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"OK"}`)) // Return JSON response
	}).Methods("GET")

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Public routes
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Protected routes
	protected := api.PathPrefix("/auth").Subrouter()
	protected.Use(authMiddleware.ValidateJWT)
	protected.HandleFunc("/validate", authHandler.ValidateToken).Methods("GET")

	return r
}
