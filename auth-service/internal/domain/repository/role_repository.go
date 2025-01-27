package repository

import (
	"context"

	"http_server/auth-service/internal/domain/models"

	"github.com/google/uuid"
)

type RoleRepository interface {
	// Existing methods
	Create(ctx context.Context, role *models.Role) error
	FindByName(ctx context.Context, name string) (*models.Role, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	AssignRoleToUser(ctx context.Context, userRole *models.UserRole) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error)
	RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error

	// Role hierarchy management
	AddParentRole(ctx context.Context, roleID, parentRoleID uuid.UUID) error
	RemoveParentRole(ctx context.Context, roleID, parentRoleID uuid.UUID) error
	GetParentRoles(ctx context.Context, roleID uuid.UUID) ([]models.Role, error)
	GetChildRoles(ctx context.Context, roleID uuid.UUID) ([]models.Role, error)

	// Batch operations
	BatchAssignRolesToUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error
	BatchRemoveRolesFromUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error

	// Role permission management
	AddPermissionToRole(ctx context.Context, roleID uuid.UUID, permission string) error
	RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permission string) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]string, error)
	HasPermission(ctx context.Context, roleID uuid.UUID, permission string) (bool, error)
}
