package models

import (
	"github.com/google/uuid"
)

// RoleHierarchy represents the hierarchical relationship between roles
type RoleHierarchy struct {
	RoleID       uuid.UUID `gorm:"primaryKey;type:uuid"`
	ParentRoleID uuid.UUID `gorm:"primaryKey;type:uuid"`
	Role         Role      `gorm:"foreignKey:RoleID"`
	ParentRole   Role      `gorm:"foreignKey:ParentRoleID"`
}

// RolePermission represents the permissions assigned to a role
type RolePermission struct {
	RoleID     uuid.UUID `gorm:"primaryKey;type:uuid"`
	Permission string    `gorm:"primaryKey;size:255"`
	Role       Role      `gorm:"foreignKey:RoleID"`
}
