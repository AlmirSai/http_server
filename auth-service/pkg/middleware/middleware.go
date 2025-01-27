package middleware

import (
	"http_server/auth-service/internal/service"
	"http_server/auth-service/pkg/logging"
	"http_server/auth-service/pkg/monitoring"
)

// Config holds all middleware configuration
type Config struct {
	RateLimit struct {
		Enabled           bool
		RequestsPerMinute int
		Burst             int
	}
	Security struct {
		Enabled bool
	}
	Auth struct {
		Enabled bool
	}
	RBAC struct {
		Enabled bool
	}
}

// Middleware holds all middleware instances and their dependencies
type Middleware struct {
	auth      *AuthMiddleware
	rbac      *RBACMiddleware
	rateLimit *RateLimiter
	security  *SecurityMiddleware
	logging   *logging.Logger
	metrics   *monitoring.Metrics
}

// NewMiddleware creates a new Middleware instance with all components
func NewMiddleware(config *Config, authService service.AuthService, logger *logging.Logger, metrics *monitoring.Metrics) *Middleware {
	m := &Middleware{
		logging: logger,
		metrics: metrics,
	}

	// Initialize enabled middleware components
	if config.Auth.Enabled {
		m.auth = NewAuthMiddleware(authService, logger, metrics)
	}

	if config.RBAC.Enabled {
		m.rbac = NewRBACMiddleware(authService, logger, metrics)
	}

	if config.RateLimit.Enabled {
		m.rateLimit = NewRateLimiter(config.RateLimit.RequestsPerMinute, config.RateLimit.Burst)
	}

	if config.Security.Enabled {
		m.security = NewSecurityMiddleware(config.RateLimit.RequestsPerMinute, config.RateLimit.Burst)
	}

	return m
}

// GetAuthMiddleware returns the auth middleware instance
func (m *Middleware) GetAuthMiddleware() *AuthMiddleware {
	return m.auth
}

// GetRBACMiddleware returns the RBAC middleware instance
func (m *Middleware) GetRBACMiddleware() *RBACMiddleware {
	return m.rbac
}

// GetRateLimiter returns the rate limiter middleware instance
func (m *Middleware) GetRateLimiter() *RateLimiter {
	return m.rateLimit
}

// GetSecurityMiddleware returns the security middleware instance
func (m *Middleware) GetSecurityMiddleware() *SecurityMiddleware {
	return m.security
}
