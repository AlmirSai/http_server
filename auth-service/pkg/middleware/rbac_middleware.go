package middleware

import (
	"net/http"

	"http_server/auth-service/internal/service"
	"http_server/auth-service/pkg/logging"
	"http_server/auth-service/pkg/monitoring"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RBACMiddleware struct {
	authService service.AuthService
	logger      *logging.Logger
	metrics     *monitoring.Metrics
}

func NewRBACMiddleware(authService service.AuthService, logger *logging.Logger, metrics *monitoring.Metrics) *RBACMiddleware {
	return &RBACMiddleware{
		authService: authService,
		logger:      logger,
		metrics:     metrics,
	}
}

func (m *RBACMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceID := r.Context().Value("trace_id").(string)
			logger := m.logger.WithContext(r.Context()).With(zap.String("trace_id", traceID))
			logger.Debug("Checking user role permissions", zap.Strings("required_roles", roles))

			m.metrics.RBACRequests.Inc()

			userID, ok := r.Context().Value("user_id").(string)
			if !ok {
				logger.Error("User ID not found in context")
				m.metrics.RBACFailures.WithLabelValues("missing_user_id").Inc()
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userUUID, err := uuid.Parse(userID)
			if err != nil {
				logger.Error("Invalid user ID format", zap.Error(err), zap.String("user_id", userID))
				m.metrics.RBACFailures.WithLabelValues("invalid_user_id").Inc()
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}

			userRoles, err := m.authService.GetUserRoles(r.Context(), userUUID)
			if err != nil {
				logger.Error("Failed to get user roles",
					zap.Error(err),
					zap.String("user_id", userID))
				m.metrics.RBACFailures.WithLabelValues("role_lookup_failed").Inc()
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			hasRequiredRole := false
			for _, userRole := range userRoles {
				for _, requiredRole := range roles {
					if userRole == requiredRole {
						hasRequiredRole = true
						break
					}
				}
				if hasRequiredRole {
					break
				}
			}

			if !hasRequiredRole {
				logger.Warn("User does not have required role",
					zap.String("user_id", userID),
					zap.Strings("required_roles", roles),
					zap.Strings("user_roles", userRoles))
				m.metrics.RBACFailures.WithLabelValues("insufficient_permissions").Inc()
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			logger.Debug("Role permission check passed",
				zap.String("user_id", userID),
				zap.Strings("required_roles", roles),
				zap.Strings("user_roles", userRoles))

			m.metrics.RBACSuccess.Inc()
			next.ServeHTTP(w, r)
		})
	}
}
