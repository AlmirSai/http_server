package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	RoleID    uuid.UUID  `json:"role_id" gorm:"type:uuid;not null"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" gorm:"index"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	CreatedBy uuid.UUID  `json:"created_by" gorm:"type:uuid"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty" gorm:"type:uuid"`
	User      User       `json:"user" gorm:"foreignKey:UserID"`
	Role      Role       `json:"role" gorm:"foreignKey:RoleID"`
}
