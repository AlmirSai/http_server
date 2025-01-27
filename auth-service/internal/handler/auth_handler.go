package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"http_server/auth-service/internal/service"
	"http_server/auth-service/internal/validator"
	"http_server/auth-service/pkg/logging"
	"http_server/auth-service/pkg/monitoring"

	"go.uber.org/zap"
)

type AuthHandler struct {
	authService service.AuthService
	logger      *logging.Logger
	metrics     *monitoring.Metrics
}

func NewAuthHandler(authService service.AuthService, logger *logging.Logger, metrics *monitoring.Metrics) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
		metrics:     metrics,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.WithContext(r.Context())
	logger.Info("Handling register request")
	h.metrics.RegisterRequests.Inc()

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode request payload", err)
		h.metrics.RegisterFailures.WithLabelValues("invalid_payload").Inc()
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		logger.Error("Invalid email format", err, zap.String("email", req.Email))
		h.metrics.RegisterFailures.WithLabelValues("invalid_email").Inc()
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.ValidatePassword(req.Password); err != nil {
		logger.Error("Invalid password format", err)
		h.metrics.RegisterFailures.WithLabelValues("invalid_password").Inc()
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.ValidateName(req.Name); err != nil {
		logger.Error("Invalid name format", err, zap.String("name", req.Name))
		h.metrics.RegisterFailures.WithLabelValues("invalid_name").Inc()
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		if err == service.ErrUserExists {
			logger.Warn("Attempted to register existing user", zap.String("email", req.Email))
			h.metrics.RegisterFailures.WithLabelValues("user_exists").Inc()
			respondWithError(w, http.StatusConflict, "User already exists")
			return
		}
		logger.Error("Failed to register user", err, zap.String("email", req.Email))
		h.metrics.RegisterFailures.WithLabelValues("internal_error").Inc()
		respondWithError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	logger.Info("User registered successfully", zap.String("user_id", user.ID.String()), zap.String("email", user.Email))
	h.metrics.RegisterSuccess.Inc()
	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.WithContext(r.Context())
	logger.Info("Handling login request")
	h.metrics.LoginRequests.Inc()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode request payload", err)
		h.metrics.LoginFailures.WithLabelValues("invalid_payload").Inc()
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		logger.Error("Invalid email format", err, zap.String("email", req.Email))
		h.metrics.LoginFailures.WithLabelValues("invalid_email").Inc()
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.ValidatePassword(req.Password); err != nil {
		logger.Error("Invalid password format", err)
		h.metrics.LoginFailures.WithLabelValues("invalid_password").Inc()
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			logger.Warn("Invalid login credentials", zap.String("email", req.Email))
			h.metrics.LoginFailures.WithLabelValues("invalid_credentials").Inc()
			respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		logger.Error("Failed to login", err, zap.String("email", req.Email))
		h.metrics.LoginFailures.WithLabelValues("internal_error").Inc()
		respondWithError(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	logger.Info("User logged in successfully", zap.String("email", req.Email))
	h.metrics.LoginSuccess.Inc()
	respondWithJSON(w, http.StatusOK, AuthResponse{Token: token})
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.WithContext(r.Context())
	logger.Info("Handling token validation request")
	h.metrics.TokenValidationRequests.Inc()

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		err := errors.New("user_id not found in context")
		logger.Error("Failed to get user_id from context", err)
		h.metrics.TokenValidationFailures.WithLabelValues("invalid_context").Inc()
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	email, ok := r.Context().Value("email").(string)
	if !ok {
		err := errors.New("email not found in context")
		logger.Error("Failed to get email from context", err)
		h.metrics.TokenValidationFailures.WithLabelValues("invalid_context").Inc()
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	roles, ok := r.Context().Value("roles").([]interface{})
	if !ok {
		err := errors.New("roles not found in context")
		logger.Error("Failed to get roles from context", err)
		h.metrics.TokenValidationFailures.WithLabelValues("invalid_context").Inc()
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := map[string]interface{}{
		"user_id": userID,
		"email":   email,
		"roles":   roles,
	}

	logger.Info("Token validated successfully",
		zap.String("user_id", userID),
		zap.String("email", email),
		zap.Any("roles", roles))

	h.metrics.TokenValidationSuccess.Inc()
	respondWithJSON(w, http.StatusOK, response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, AuthResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
