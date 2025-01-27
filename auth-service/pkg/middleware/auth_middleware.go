package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"http_server/auth-service/internal/service"
	"http_server/auth-service/pkg/logging"
	"http_server/auth-service/pkg/monitoring"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// authContextKey is a custom type for auth-related context keys to avoid collisions
type authContextKey string

// Auth-related context keys
const (
	TraceIDKey = authContextKey("trace_id")
	UserIDKey  = authContextKey("user_id")
	EmailKey   = authContextKey("email")
	RolesKey   = authContextKey("roles")
)

type AuthMiddleware struct {
	authService service.AuthService
	logger      *logging.Logger
	metrics     *monitoring.Metrics
}

func NewAuthMiddleware(authService service.AuthService, logger *logging.Logger, metrics *monitoring.Metrics) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
		metrics:     metrics,
	}
}

func (m *AuthMiddleware) ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), TraceIDKey, traceID)
		logger := m.logger.WithContext(ctx).With(
			zap.String("trace_id", traceID),
			zap.String("remote_ip", r.RemoteAddr),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)

		logger.Debug("Starting JWT token validation")
		m.metrics.AuthRequests.Inc()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Warn("Missing authorization header")
			m.metrics.AuthFailures.WithLabelValues("missing_header").Inc()
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			logger.Warn("Invalid authorization format",
				zap.String("auth_header", authHeader))
			m.metrics.AuthFailures.WithLabelValues("invalid_format").Inc()
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		var token *jwt.Token
		var err error

		timer := prometheus.NewTimer(m.metrics.TokenValidationDuration)
		defer timer.ObserveDuration()

		token, err = m.authService.ValidateToken(bearerToken[1])

		if err != nil {
			logger.Warn("Token validation failed",
				zap.Error(err),
				zap.String("token_length", fmt.Sprintf("%d", len(bearerToken[1]))),
				zap.String("validation_error", err.Error()))
			m.metrics.AuthFailures.WithLabelValues("invalid_token").Inc()

			var statusCode int
			var message string

			switch err {
			case service.ErrTokenExpired:
				statusCode = http.StatusUnauthorized
				message = "Token has expired"
			case service.ErrInvalidToken:
				statusCode = http.StatusUnauthorized
				message = "Invalid token"
			default:
				statusCode = http.StatusInternalServerError
				message = "Internal server error"
			}

			http.Error(w, message, statusCode)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Error("Failed to parse token claims")
			m.metrics.AuthFailures.WithLabelValues("invalid_claims").Inc()
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			logger.Error("Invalid user_id claim type")
			m.metrics.AuthFailures.WithLabelValues("invalid_user_id").Inc()
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			logger.Error("Invalid email claim type")
			m.metrics.AuthFailures.WithLabelValues("invalid_email").Inc()
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		roles, ok := claims["roles"].([]interface{})
		if !ok {
			logger.Error("Invalid roles claim type")
			m.metrics.AuthFailures.WithLabelValues("invalid_roles").Inc()
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, UserIDKey, userID)
		ctx = context.WithValue(ctx, EmailKey, email)
		ctx = context.WithValue(ctx, RolesKey, roles)

		logger.Debug("JWT token validated successfully",
			zap.String("user_id", userID),
			zap.String("email", email),
			zap.Any("roles", roles))

		m.metrics.AuthSuccess.Inc()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
