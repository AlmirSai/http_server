package repository

import (
	"context"
	"errors"

	"http_server/auth-service/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *models.Role) error {
	result := r.db.WithContext(ctx).Create(role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrDuplicateKey
		}
		return result.Error
	}
	return nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &role, nil
}

func (r *roleRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	result := r.db.WithContext(ctx).First(&role, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &role, nil
}

func (r *roleRepository) AssignRoleToUser(ctx context.Context, userRole *models.UserRole) error {
	result := r.db.WithContext(ctx).Create(userRole)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrDuplicateKey
		}
		return result.Error
	}
	return nil
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error) {
	var roles []models.Role
	result := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (r *roleRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&models.UserRole{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *roleRepository) AddParentRole(ctx context.Context, roleID, parentRoleID uuid.UUID) error {
	result := r.db.WithContext(ctx).Create(&models.RoleHierarchy{
		RoleID:       roleID,
		ParentRoleID: parentRoleID,
	})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrDuplicateKey
		}
		return result.Error
	}
	return nil
}

func (r *roleRepository) RemoveParentRole(ctx context.Context, roleID, parentRoleID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("role_id = ? AND parent_role_id = ?", roleID, parentRoleID).
		Delete(&models.RoleHierarchy{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *roleRepository) GetParentRoles(ctx context.Context, roleID uuid.UUID) ([]models.Role, error) {
	var roles []models.Role
	result := r.db.WithContext(ctx).
		Joins("JOIN role_hierarchies ON roles.id = role_hierarchies.parent_role_id").
		Where("role_hierarchies.role_id = ?", roleID).
		Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (r *roleRepository) GetChildRoles(ctx context.Context, roleID uuid.UUID) ([]models.Role, error) {
	var roles []models.Role
	result := r.db.WithContext(ctx).
		Joins("JOIN role_hierarchies ON roles.id = role_hierarchies.role_id").
		Where("role_hierarchies.parent_role_id = ?", roleID).
		Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (r *roleRepository) BatchAssignRolesToUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error {
	if len(roleIDs) == 0 {
		return nil
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	userRoles := make([]models.UserRole, len(roleIDs))
	for i, roleID := range roleIDs {
		// Verify role exists
		var exists bool
		if err := tx.Model(&models.Role{}).Select("count(*) > 0").Where("id = ?", roleID).Scan(&exists).Error; err != nil {
			tx.Rollback()
			return err
		}
		if !exists {
			tx.Rollback()
			return ErrNotFound
		}

		userRoles[i] = models.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
	}

	result := tx.Create(&userRoles)
	if result.Error != nil {
		tx.Rollback()
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrDuplicateKey
		}
		return result.Error
	}

	return tx.Commit().Error
}

func (r *roleRepository) BatchRemoveRolesFromUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error {
	if len(roleIDs) == 0 {
		return nil
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Verify user exists
	var exists bool
	if err := tx.Model(&models.User{}).Select("count(*) > 0").Where("id = ?", userID).Scan(&exists).Error; err != nil {
		tx.Rollback()
		return err
	}
	if !exists {
		tx.Rollback()
		return ErrNotFound
	}

	result := tx.Where("user_id = ? AND role_id IN ?", userID, roleIDs).Delete(&models.UserRole{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

func (r *roleRepository) AddPermissionToRole(ctx context.Context, roleID uuid.UUID, permission string) error {
	result := r.db.WithContext(ctx).Create(&models.RolePermission{
		RoleID:     roleID,
		Permission: permission,
	})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrDuplicateKey
		}
		return result.Error
	}
	return nil
}

func (r *roleRepository) RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permission string) error {
	result := r.db.WithContext(ctx).
		Where("role_id = ? AND permission = ?", roleID, permission).
		Delete(&models.RolePermission{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *roleRepository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]string, error) {
	var permissions []string
	result := r.db.WithContext(ctx).
		Model(&models.RolePermission{}).
		Where("role_id = ?", roleID).
		Pluck("permission", &permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}

func (r *roleRepository) HasPermission(ctx context.Context, roleID uuid.UUID, permission string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&models.RolePermission{}).
		Where("role_id = ? AND permission = ?", roleID, permission).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
