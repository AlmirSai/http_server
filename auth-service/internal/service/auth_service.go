package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"http_server/auth-service/internal/domain/models"
	"http_server/auth-service/internal/domain/repository"
	"http_server/auth-service/pkg/logging"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

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

type AuthService interface {
	Register(ctx context.Context, email, password, name string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	AssignRole(ctx context.Context, userID uuid.UUID, roleName string) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error)
	RemoveRole(ctx context.Context, userID uuid.UUID, roleName string) error
}

type authService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	jwtKey   []byte
	logger   *logging.Logger
}

func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, jwtKey []byte, logger *logging.Logger) AuthService {
	return &authService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		jwtKey:   jwtKey,
		logger:   logger,
	}
}

func (s *authService) Register(ctx context.Context, email, password, name string) (*models.User, error) {
	logger := s.logger.WithContext(ctx)
	logger.Info("Starting user registration", zap.String("email", email))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", err, zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			logger.Warn("Attempted to register existing user", zap.String("email", email))
			return nil, ErrUserExists
		}
		logger.Error("Failed to create user", err, zap.String("email", email))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign default user role
	defaultRole, err := s.roleRepo.FindByName(ctx, models.RoleUser)
	if err != nil {
		logger.Error("Failed to find default role", err)
		return nil, fmt.Errorf("failed to find default role: %w", err)
	}

	userRole := &models.UserRole{
		UserID: user.ID,
		RoleID: defaultRole.ID,
	}

	if err := s.roleRepo.AssignRoleToUser(ctx, userRole); err != nil {
		logger.Error("Failed to assign default role", err)
		return nil, fmt.Errorf("failed to assign default role: %w", err)
	}

	logger.Info("User registered successfully", zap.String("user_id", user.ID.String()), zap.String("email", user.Email))
	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	logger := s.logger.WithContext(ctx)
	logger.Info("Attempting user login", zap.String("email", email))

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.Warn("User not found during login", zap.String("email", email))
			return "", ErrInvalidCredentials
		}
		logger.Error("Failed to find user", err, zap.String("email", email))
		return "", fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logger.Warn("Invalid password attempt", zap.String("email", email))
		return "", ErrInvalidCredentials
	}

	// Get user roles for token
	userRoles, err := s.GetUserRoles(ctx, user.ID)
	if err != nil {
		logger.Error("Failed to get user roles", err)
		return "", fmt.Errorf("failed to get user roles: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"roles":   userRoles,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		logger.Error("Failed to generate JWT token", err, zap.String("user_id", user.ID.String()))
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	logger.Info("User logged in successfully", zap.String("user_id", user.ID.String()), zap.String("email", user.Email))
	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	logger := s.logger
	logger.Debug("Validating JWT token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Warn("Invalid token signing method", zap.String("method", token.Method.Alg()))
			return nil, ErrInvalidToken
		}
		return s.jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			logger.Warn("Token has expired")
			return nil, ErrTokenExpired
		}
		logger.Warn("Token validation failed", zap.Error(err))
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				logger.Warn("Token has expired")
				return nil, ErrTokenExpired
			}
		} else {
			logger.Warn("Invalid expiration claim in token")
			return nil, ErrInvalidToken
		}

		if _, ok := claims["user_id"].(string); !ok {
			logger.Warn("Invalid user_id claim in token")
			return nil, ErrInvalidToken
		}
	}

	logger.Debug("Token validated successfully")
	return token, nil
}

func (s *authService) AssignRole(ctx context.Context, userID uuid.UUID, roleName string) error {
	logger := s.logger.WithContext(ctx)
	logger.Info("Assigning role to user", zap.String("user_id", userID.String()), zap.String("role", roleName))

	// Validate user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.Error("User not found", err, zap.String("user_id", userID.String()))
			return ErrUserNotFound
		}
		logger.Error("Failed to find user", err)
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Validate user is active
	if !user.Active {
		logger.Warn("Attempted to assign role to inactive user", zap.String("user_id", userID.String()))
		return fmt.Errorf("cannot assign role to inactive user")
	}

	role, err := s.roleRepo.FindByName(ctx, roleName)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.Error("Role not found", err, zap.String("role", roleName))
			return ErrRoleNotFound
		}
		logger.Error("Failed to find role", err, zap.String("role", roleName))
		return fmt.Errorf("failed to find role: %w", err)
	}

	userRole := &models.UserRole{
		UserID: userID,
		RoleID: role.ID,
	}

	if err := s.roleRepo.AssignRoleToUser(ctx, userRole); err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			logger.Warn("Role already assigned to user",
				zap.String("user_id", userID.String()),
				zap.String("role", roleName))
			return ErrRoleAlreadyAssigned
		}
		logger.Error("Failed to assign role", err)
		return fmt.Errorf("failed to assign role: %w", err)
	}

	logger.Info("Role assigned successfully", zap.String("user_id", userID.String()), zap.String("role", roleName))
	return nil
}

func (s *authService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	logger := s.logger.WithContext(ctx)
	logger.Debug("Getting user roles", zap.String("user_id", userID.String()))

	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user roles", err)
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	return roleNames, nil
}

func (s *authService) RemoveRole(ctx context.Context, userID uuid.UUID, roleName string) error {
	logger := s.logger.WithContext(ctx)
	logger.Info("Removing role from user", zap.String("user_id", userID.String()), zap.String("role", roleName))

	role, err := s.roleRepo.FindByName(ctx, roleName)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.Error("Role not found", err, zap.String("role", roleName))
			return ErrRoleNotFound
		}
		logger.Error("Failed to find role", err, zap.String("role", roleName))
		return fmt.Errorf("failed to find role: %w", err)
	}

	if err := s.roleRepo.RemoveRoleFromUser(ctx, userID, role.ID); err != nil {
		logger.Error("Failed to remove role", err)
		return fmt.Errorf("failed to remove role: %w", err)
	}

	logger.Info("Role removed successfully", zap.String("user_id", userID.String()), zap.String("role", roleName))
	return nil
}
