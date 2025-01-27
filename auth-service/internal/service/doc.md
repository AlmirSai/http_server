# Service Package Documentation

## Overview
The service package implements the core business logic for the authentication service. It handles user authentication, token management, and role-based access control.

## Components

### AuthService Interface
```go
type AuthService interface {
    Register(ctx context.Context, email, password, name string) (*models.User, error)
    Login(ctx context.Context, email, password string) (string, error)
    ValidateToken(token string) (*jwt.Token, error)
    AssignRole(ctx context.Context, userID uuid.UUID, roleName string) error
    GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error)
    RemoveRole(ctx context.Context, userID uuid.UUID, roleName string) error
}
```

## Core Operations

### User Registration
- Validates user input (email, password, name)
- Hashes password using bcrypt
- Creates user record in database
- Assigns default user role
- Returns created user or appropriate error

### User Authentication
- Validates credentials
- Compares password hash
- Retrieves user roles
- Generates JWT token with claims:
  - User ID
  - Email
  - Roles
  - Expiration time

### Token Management
- JWT-based authentication
- Token validation and verification
- Expiration handling
- Claim verification

### Role Management
- Role assignment to users
- Role removal from users
- User role retrieval
- Default role handling

## Error Types
```go
var (
    ErrInvalidCredentials  = errors.New("invalid credentials")
    ErrUserExists          = errors.New("user already exists")
    ErrTokenExpired        = errors.New("token has expired")
    ErrInvalidToken        = errors.New("invalid token")
    ErrRoleNotFound        = errors.New("role not found")
    ErrRoleAlreadyAssigned = errors.New("role already assigned to user")
    ErrUserNotFound        = errors.New("user not found")
    ErrInvalidPassword     = errors.New("invalid password format")
    ErrInvalidEmail        = errors.New("invalid email format")
)
```

## Dependencies
- UserRepository: User data management
- RoleRepository: Role data management
- JWT: Token generation and validation
- Bcrypt: Password hashing
- Logger: Operation logging

## Security Features
1. Password hashing with bcrypt
2. JWT-based authentication
3. Role-based access control
4. Input validation
5. Secure error handling

## Best Practices
1. Context-based operations
2. Comprehensive error handling
3. Secure password management
4. Proper logging
5. Clean interface design
6. Repository pattern usage