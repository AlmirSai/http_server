# Server Package Documentation

## Overview
The server package implements the HTTP server infrastructure for the authentication service. It handles server configuration, routing setup, and server lifecycle management.

## Components

### Server
```go
type Server struct {
    httpServer *http.Server
}
```
The Server struct encapsulates the HTTP server configuration and provides methods for server lifecycle management.

#### Methods
- `NewServer(cfg *config.Config, authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) *Server`
  - Creates a new server instance with configured timeouts and routing
  - Initializes the HTTP server with the provided configuration

- `Run() error`
  - Starts the HTTP server and begins accepting connections
  - Blocks until the server is shut down

- `Shutdown(ctx context.Context) error`
  - Gracefully shuts down the server
  - Waits for active connections to complete

### Router
Implements the HTTP routing using gorilla/mux with the following features:

#### Middleware Configuration
- Logging middleware for request tracking
- CORS middleware with configurable options
  - Allows all origins (*)
  - Supports GET, POST, PUT, DELETE, OPTIONS methods
  - Allows Content-Type and Authorization headers

#### Endpoints
- Health Check: `GET /health`
  - Returns server status

#### API Routes (v1)
Base path: `/api/v1`

Public Routes:
- `POST /auth/register` - User registration
- `POST /auth/login` - User authentication

Protected Routes:
- `GET /auth/validate` - Token validation (requires JWT)

### Server Options
```go
type Options struct {
    Port         int
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
}
```

Default Configuration:
- Port: 8080
- Read Timeout: 15 seconds
- Write Timeout: 15 seconds
- Idle Timeout: 15 seconds

## Best Practices
1. Use middleware for cross-cutting concerns
2. Implement proper timeout handling
3. Configure CORS for security
4. Use subrouters for route grouping
5. Implement health checks
6. Handle graceful shutdown