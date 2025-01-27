# Validator Package Documentation

## Overview
The validator package provides a robust set of validation utilities for ensuring data integrity and security in the authentication service. It implements validation rules for common fields like email addresses, passwords, and user names.

## Components

### ValidationError
A custom error type that provides structured validation error information:
```go
type ValidationError struct {
    Field   string
    Message string
}
```

### Email Validation
Validates email addresses using the following rules:
- Non-empty requirement
- RFC-compliant email format validation
- Proper domain structure verification

Example usage:
```go
err := validator.ValidateEmail("user@example.com")
if err != nil {
    // Handle validation error
}
```

### Password Validation
Enforces secure password requirements:
- Minimum length of 8 characters
- Must contain at least:
  - One uppercase letter
  - One lowercase letter
  - One number

Example usage:
```go
err := validator.ValidatePassword("SecurePass123")
if err != nil {
    // Handle validation error
}
```

### Name Validation
Ensures proper formatting of user names:
- Non-empty requirement
- Minimum length of 2 characters
- Trims leading and trailing whitespace

Example usage:
```go
err := validator.ValidateName("John")
if err != nil {
    // Handle validation error
}
```

## Error Handling
All validation functions return a `ValidationError` that includes:
- The field name that failed validation
- A descriptive error message

Example error handling:
```go
if err := validator.ValidateEmail(email); err != nil {
    if validErr, ok := err.(*validator.ValidationError); ok {
        fmt.Printf("Validation failed for %s: %s\n", validErr.Field, validErr.Message)
    }
}
```

## Best Practices
1. Always validate input data before processing
2. Handle validation errors appropriately in your application
3. Use the validation functions in combination for complete form validation
4. Consider the validation rules when designing user interfaces

## Configuration
The validator package includes configurable constants:
- `passwordMinLength`: Minimum required password length (default: 8)
- `emailRegex`: Regular expression pattern for email validation

## Thread Safety
All validation functions are stateless and thread-safe, making them suitable for concurrent use in web applications.