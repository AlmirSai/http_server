package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID              `json:"id" gorm:"primaryKey;type:uuid"`
	Name        string                 `json:"name" gorm:"unique;not null"`
	Description string                 `json:"description"`
	ParentID    *uuid.UUID             `json:"parent_id,omitempty" gorm:"type:uuid"`
	Metadata    map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	Priority    int                    `json:"priority" gorm:"default:0"`
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
	Parent      *Role                  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []*Role                `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Users       []UserRole             `json:"users,omitempty" gorm:"foreignKey:RoleID"`
	Permissions []RolePermission       `json:"permissions,omitempty" gorm:"foreignKey:RoleID"`
}

// Predefined roles
const (
	RoleAdmin        = "admin"
	RoleUser         = "user"
	RoleGuest        = "guest"
	RoleDeveloper    = "developer"
	RoleCompany      = "company"
	RoleOfficialNews = "official_news"
	RoleImportant    = "important_person"
)

// RolePermissions defines the available permissions for each role
var RolePermissions = map[string][]string{
	RoleAdmin: {"*"}, // Admin has all permissions
	RoleUser: {
		"read:profile",
		"update:profile",
		"read:posts",
		"create:posts",
	},
	RoleGuest: {
		"read:posts",
	},
	RoleDeveloper: {
		"read:profile",
		"update:profile",
		"read:posts",
		"create:posts",
		"update:posts",
		"delete:posts",
		"access:api",
	},
	RoleCompany: {
		"read:profile",
		"update:profile",
		"read:posts",
		"create:posts",
		"update:posts",
		"create:official_posts",
	},
	RoleOfficialNews: {
		"read:profile",
		"update:profile",
		"read:posts",
		"create:posts",
		"update:posts",
		"create:news",
		"update:news",
	},
	RoleImportant: {
		"read:profile",
		"update:profile",
		"read:posts",
		"create:posts",
		"update:posts",
		"create:verified_posts",
	},
}

// GetAllPermissions returns all permissions for this role including inherited ones
func (r *Role) GetAllPermissions() []string {
	perms := make(map[string]bool)

	// Add own permissions
	if rolePerms, exists := RolePermissions[r.Name]; exists {
		for _, p := range rolePerms {
			perms[p] = true
		}
	}

	// Add parent permissions
	if r.Parent != nil {
		parentPerms := r.Parent.GetAllPermissions()
		for _, p := range parentPerms {
			perms[p] = true
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(perms))
	for p := range perms {
		result = append(result, p)
	}

	return result
}

// HasPermission checks if a role has a specific permission
func (r *Role) HasPermission(permission string) bool {
	// Check own permissions first
	perms, exists := RolePermissions[r.Name]
	if exists {
		for _, p := range perms {
			if p == "*" || p == permission {
				return true
			}
		}
	}

	// Check parent permissions
	if r.Parent != nil {
		return r.Parent.HasPermission(permission)
	}

	return false
}
