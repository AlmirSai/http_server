package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid"`
	Email         string     `json:"email" gorm:"unique;not null"`
	Password      string     `json:"-" gorm:"not null"`
	Name          string     `json:"name" gorm:"not null"`
	Active        bool       `json:"active" gorm:"default:true"`
	EmailVerified bool       `json:"email_verified" gorm:"default:false"`
	LastLoginAt   time.Time  `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Profile information
	PhoneNumber    string     `json:"phone_number,omitempty"`
	ProfilePicture string     `json:"profile_picture,omitempty"`
	Bio            string     `json:"bio,omitempty"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty"`
	Location       string     `json:"location,omitempty"`

	// Status tracking
	LastActivityAt     time.Time  `json:"last_activity_at,omitempty"`
	StatusChangedAt    time.Time  `json:"status_changed_at,omitempty"`
	DeactivatedAt      *time.Time `json:"deactivated_at,omitempty"`
	DeactivationReason string     `json:"deactivation_reason,omitempty"`

	// Security and verification
	TwoFactorEnabled    bool       `json:"two_factor_enabled" gorm:"default:false"`
	VerifiedAt          *time.Time `json:"verified_at,omitempty"`
	FailedLoginAttempts int        `json:"failed_login_attempts,omitempty" gorm:"default:0"`
	LockedUntil         *time.Time `json:"locked_until,omitempty"`
}
