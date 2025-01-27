# Auth Service

## Overview
The Auth Service is responsible for user authentication, authorization, and role management in the microservices architecture. It provides secure JWT token generation and validation, along with comprehensive user role management capabilities.

## Features
- User registration and authentication
- JWT token generation and validation
- Role-based access control (RBAC)
- User role management
- Secure password hashing
- Prometheus metrics integration
- Structured logging with Zap
- Rate limiting and request throttling
- Session management with Redis
- Password reset functionality

## Prerequisites
- Go 1.19 or later
- PostgreSQL 14+
- Redis 6+
- Make (optional)
- Docker (optional)

## Local Development Setup

### 1. Configure the Service
Copy the example configuration and modify it for your environment:
```bash
cp config.example.yaml config.yaml
```

Update the following in `config.yaml`:
- Database connection settings
- Redis connection settings
- JWT secret key and configuration
- Server port and timeouts
- Rate limiting parameters
- Logging configuration

### 2. Set Up the Database
Ensure PostgreSQL is running and create a database:
```bash
psql -U postgres
CREATE DATABASE auth_service;
```

### 3. Run the Service
```bash
# Using Docker
docker build -t auth-service .
docker run -p 8080:8080 auth-service

# Or build and run locally
go build -o auth-service ./cmd/app
./auth-service

# Or using go run
go run ./cmd/app/main.go
```

### 4. Verify the Service
The service should be running on the configured port (default: 8080).
Test the health check endpoint:
```bash
curl http://localhost:8080/health
```

## API Endpoints

### Authentication
#### Register User
```
POST /api/v1/auth/register
Content-Type: application/json

{
    "username": "string",
    "email": "string",
    "password": "string"
}
```

#### Login
```
POST /api/v1/auth/login
Content-Type: application/json

{
    "email": "string",
    "password": "string"
}
```

#### Refresh Token
```
POST /api/v1/auth/refresh
Authorization: Bearer <refresh_token>
```

### Role Management
#### Assign Role
```
POST /api/v1/auth/roles
Content-Type: application/json
Authorization: Bearer <access_token>

{
    "user_id": "string",
    "role": "string"
}
```

#### Remove Role
```
DELETE /api/v1/auth/roles
Content-Type: application/json
Authorization: Bearer <access_token>

{
    "user_id": "string",
    "role": "string"
}
```

#### Get User Roles
```
GET /api/v1/auth/roles/{user_id}
Authorization: Bearer <access_token>
```

## Configuration

### Environment Variables
- `APP_ENV` - Application environment (development, production)
- `CONFIG_PATH` - Path to configuration file (default: config.yaml)
- `JWT_SECRET` - JWT signing secret (overrides config file)
- `DB_URL` - Database connection URL (overrides config file)

### Configuration File (config.yaml)
```yaml
server:
  port: 8080
  timeout: 30s
  rate_limit:
    requests: 100
    duration: 1m

database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: auth_service
  sslmode: disable
  max_connections: 100

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret_key: your_secret_key
  access_token_expiration: 15m
  refresh_token_expiration: 24h

logging:
  level: debug
  output: stdout
  format: json
```

## Security Considerations
- Passwords are hashed using bcrypt
- JWT tokens are signed with HS256 algorithm
- Rate limiting prevents brute force attacks
- Input validation for all API endpoints
- CORS configuration for web clients
- TLS/SSL in production environment

## Testing

### Unit Tests
```bash
go test -v ./...
```

### Integration Tests
```bash
go test -tags=integration -v ./...
```

### Load Testing
```bash
# Using k6 for load testing
k6 run tests/load/auth_test.js
```

## Metrics
Prometheus metrics are available at `/metrics` endpoint, including:
- Request latencies (histogram)
- Error rates by endpoint
- Authentication success/failure counts
- Active sessions count
- Rate limiting metrics

## Logging
Structured logging is implemented using Zap logger. Logs include:
- Request tracing with correlation IDs
- Error details with stack traces
- Authentication events
- Role management operations
- Performance metrics

## Development Guidelines
- Follow Go best practices and project conventions
- Write tests for new functionality
- Update API documentation when adding/modifying endpoints
- Use proper error handling and logging
- Follow security best practices
- Maintain backward compatibility

## Integration with Other Services
- Use shared authentication middleware
- Implement service discovery
- Handle cross-service communication
- Maintain API versioning
- Document service dependencies

## Troubleshooting
- Check logs for detailed error messages
- Verify database connectivity
- Ensure Redis is running for session management
- Validate JWT token configuration
- Monitor rate limiting metrics