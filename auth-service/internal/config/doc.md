# Configuration Package Documentation

## Overview
The configuration package provides a robust configuration management system for the auth service. It handles loading, parsing, and validating configuration settings from files and environment variables.

## Configuration Structure

### Main Configuration
The main `Config` struct encompasses all configuration settings:

```go
type Config struct {
    Server   ServerConfig
    JWT      JWTConfig
    Database DatabaseConfig
    Security SecurityConfig
    Logging  LoggingConfig
    Metrics  MetricsConfig
}
```

### Server Configuration
Controls HTTP server settings:
- Port: Server listening port
- ReadTimeout/WriteTimeout: Request timeouts
- TLS settings for secure connections

### JWT Configuration
Manages JSON Web Token settings:
- SecretKey: For token signing
- Expiration: Token lifetime
- Refresh token configuration
- Token rotation settings

### Database Configuration
Database connection parameters:
- Host, Port, User, Password, DBName
- Connection pool settings
- SSL mode configuration

### Security Configuration
Security-related settings:
- Rate limiting
- Password policies
- CORS and security headers

### Logging Configuration
Logging system settings:
- Log level
- Output format
- Time format

### Metrics Configuration
Metrics collection settings:
- Enable/disable metrics
- Metrics endpoint
- Service name for metrics

## Usage

### Loading Configuration
```go
config, err := config.LoadConfig("config.yaml")
if err != nil {
    // Handle error
}
```

### Environment Variables
The configuration system supports environment variable overrides:
- Environment variables take precedence over file settings
- Uses dot notation replacement (e.g., `SERVER_PORT` for `server.port`)

### Validation
The configuration package performs validation checks:
- Required fields (port, JWT secret, database settings)
- Format validation
- Logical constraints

## Best Practices
1. Always use the provided `LoadConfig` function
2. Handle validation errors appropriately
3. Use environment variables for sensitive data
4. Keep configuration files secure
5. Regularly review and update configuration settings

## Example Configuration File
```yaml
server:
  port: 8080
  readTimeout: 5s
  writeTimeout: 10s

jwt:
  secretKey: "your-secret-key"
  expiration: 24h
  refreshTokenSecret: "refresh-secret"
  refreshTokenExpiry: 168h

database:
  host: "localhost"
  port: 5432
  user: "dbuser"
  password: "dbpassword"
  dbName: "authdb"

security:
  rateLimit:
    enabled: true
    requestsPerMinute: 60
```