# Handlers Package Documentation

## Overview
The handlers package implements HTTP request handlers for the authentication service. It manages incoming HTTP requests, validates input data, interacts with the domain layer, and formats responses.

## Handler Components

### AuthHandler
Manages authentication-related endpoints:
```go
type AuthHandler struct {
    authService domain.AuthService
    validator   validator.Validator
    logger      logging.Logger
}
```

## Endpoints

### Register User
```go
// POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request)
```
**Request Body:**
```json
{
    "email": "user@example.com",
    "password": "SecurePass123!"
}
```
**Response (201 Created):**
```json
{
    "id": "user-uuid",
    "email": "user@example.com",
    "created_at": "2023-01-01T12:00:00Z"
}
```

### Login
```go
// POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request)
```
**Request Body:**
```json
{
    "email": "user@example.com",
    "password": "SecurePass123!"
}
```
**Response (200 OK):**
```json
{
    "access_token": "jwt-token",
    "refresh_token": "refresh-token",
    "expires_in": 3600
}
```

### Refresh Token
```go
// POST /auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request)
```
**Request Body:**
```json
{
    "refresh_token": "refresh-token"
}
```
**Response (200 OK):**
```json
{
    "access_token": "new-jwt-token",
    "refresh_token": "new-refresh-token",
    "expires_in": 3600
}
```

### Logout
```go
// POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request)
```
**Headers:**
- Authorization: Bearer {access_token}

**Response (204 No Content)**

## Middleware

### Authentication Middleware
```go
func AuthMiddleware(next http.Handler) http.Handler
```
- Validates JWT tokens
- Extracts user information
- Injects user context

### Rate Limiting Middleware
```go
func RateLimitMiddleware(next http.Handler) http.Handler
```
- Implements token bucket algorithm
- Configurable rate limits
- IP-based rate limiting

### Request Validation Middleware
```go
func ValidationMiddleware(next http.Handler) http.Handler
```
- Validates request bodies
- Checks required fields
- Enforces data constraints

## Error Handling

### HTTP Error Responses
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    int    `json:"code"`
}
```

### Common Error Status Codes
- 400 Bad Request: Invalid input data
- 401 Unauthorized: Invalid credentials
- 403 Forbidden: Insufficient permissions
- 404 Not Found: Resource not found
- 429 Too Many Requests: Rate limit exceeded
- 500 Internal Server Error: Server-side error

## Request Processing Flow
1. Request received by handler
2. Middleware chain execution
3. Request validation
4. Business logic execution
5. Response formatting
6. Error handling if needed

## Best Practices
1. Use middleware for cross-cutting concerns
2. Implement proper input validation
3. Return appropriate HTTP status codes
4. Include detailed error messages
5. Log all significant operations
6. Use proper HTTP methods
7. Implement proper CORS handling
8. Rate limit sensitive endpoints