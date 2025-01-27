# Domain Package Documentation

## Overview
The domain package contains the core business logic and entities for the authentication service. It defines the fundamental types, interfaces, and business rules that govern the authentication system.

## Core Entities

### User
Represents a user in the system:
```go
type User struct {
    ID        string
    Email     string
    Password  string    // Hashed password
    Role      UserRole
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Session
Represents an active user session:
```go
type Session struct {
    ID           string
    UserID       string
    RefreshToken string
    ExpiresAt    time.Time
    CreatedAt    time.Time
}
```

### UserRole
Defines user authorization levels:
```go
type UserRole string

const (
    RoleUser  UserRole = "user"
    RoleAdmin UserRole = "admin"
)
```

## Interfaces

### UserRepository
Defines data access operations for users:
```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}
```

### SessionRepository
Manages user session data:
```go
type SessionRepository interface {
    Create(ctx context.Context, session *Session) error
    Get(ctx context.Context, id string) (*Session, error)
    Delete(ctx context.Context, id string) error
    DeleteAllForUser(ctx context.Context, userID string) error
}
```

### AuthService
Defines core authentication operations:
```go
type AuthService interface {
    Register(ctx context.Context, email, password string) (*User, error)
    Login(ctx context.Context, email, password string) (*Session, error)
    Logout(ctx context.Context, sessionID string) error
    RefreshToken(ctx context.Context, refreshToken string) (*Session, error)
    ValidateToken(ctx context.Context, token string) (*User, error)
}
```

## Business Rules

### Password Requirements
- Minimum length: 8 characters
- Must contain at least one uppercase letter
- Must contain at least one lowercase letter
- Must contain at least one number
- Must contain at least one special character

### Session Management
- Access tokens expire after 1 hour
- Refresh tokens expire after 7 days
- Users can have multiple active sessions
- Sessions are invalidated on password change

### Rate Limiting
- Maximum 60 login attempts per minute per IP
- Maximum 3 failed login attempts per account before temporary lockout
- Account lockout duration: 15 minutes

### Security Policies
- Passwords are hashed using bcrypt
- Tokens are signed using HMAC-SHA256
- All sensitive operations are logged
- Input validation for all user data

## Error Handling

### Domain Errors
Standardized error types for domain operations:
```go
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrSessionExpired    = errors.New("session expired")
    ErrInvalidToken      = errors.New("invalid token")
    ErrUserLocked        = errors.New("user account locked")
)
```

## Best Practices
1. Always validate input data at the domain level
2. Use context for cancellation and timeouts
3. Keep business logic independent of infrastructure
4. Handle errors explicitly and consistently
5. Use strong types instead of primitive types
6. Implement proper logging for debugging
7. Follow SOLID principles in interface design