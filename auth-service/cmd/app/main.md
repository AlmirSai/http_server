# Auth Service - Main Application Documentation

This document provides a comprehensive overview of the auth service's main application initialization and setup process.

## Overview

The main application file (`main.go`) serves as the entry point for the authentication service. It handles the initialization and configuration of various components including:

- Configuration management
- Logging system
- Metrics monitoring
- Database connection and migration
- Service layer setup
- HTTP server initialization
- Graceful shutdown handling

## Component Initialization

### Configuration Loading

The application starts by loading configuration from `config.yaml`:

```go
cfg, err := config.LoadConfig("config.yaml")
```

This loads all necessary configuration parameters for the service including database credentials, server settings, and JWT configuration.

### Logger Setup

A structured logger is initialized using the following configuration:

- Environment from `APP_ENV` environment variable
- Log level from configuration
- File-based logging with rotation:
  - Maximum file size: 10MB
  - Maximum backups: 5
  - Maximum age: 30 days
  - Compression enabled

### Metrics Initialization

Prometheus metrics are initialized for monitoring service performance and behavior:

```go
metrics := monitoring.NewMetrics(cfg.Metrics.ServiceName)
```

### Database Connection

The service establishes a PostgreSQL database connection with the following features:

- Connection pool configuration
- Automatic schema migration
- Support for User, Role, and UserRole models

Connection pool settings:
- Maximum open connections
- Maximum idle connections
- Connection maximum lifetime

### Service Layer

The following components are initialized in sequence:

1. Repositories:
   - User Repository
   - Role Repository

2. Auth Service:
   - Handles user authentication
   - Manages JWT tokens
   - Integrates with repositories

3. HTTP Handlers and Middleware:
   - Auth Handler for HTTP endpoints
   - Auth Middleware for request authentication

### Server Initialization

The HTTP server is initialized with:
- Configured port
- Request handlers
- Middleware chain

## Graceful Shutdown

The application implements graceful shutdown handling:

1. Listens for interrupt signals (SIGINT, SIGTERM)
2. Allows ongoing requests to complete (30-second timeout)
3. Closes database connections
4. Logs shutdown completion

## Error Handling

The application implements comprehensive error handling:

- Configuration errors result in immediate shutdown
- Database connection failures are logged and terminate the application
- Server startup failures are logged with stack traces
- Shutdown errors are logged before exit

## Logging Strategy

The application uses structured logging with:

- Different log levels (DEBUG, INFO, ERROR, FATAL)
- Contextual information in log entries
- Both file and console output
- Log rotation for maintenance

## Usage

To start the service:

1. Ensure `config.yaml` is properly configured
2. Set appropriate environment variables
3. Run the application

The service will initialize all components and start serving requests.

## Dependencies

Key external dependencies:

- `go.uber.org/zap`: Structured logging
- `gorm.io/gorm`: ORM for database operations
- `gorm.io/driver/postgres`: PostgreSQL driver

## Best Practices

The implementation follows several best practices:

1. Structured error handling
2. Graceful shutdown support
3. Configuration management
4. Connection pooling
5. Metrics monitoring
6. Structured logging