package permission

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	Resource    string         `json:"resource" gorm:"not null"` // e.g., "examples", "users", "reports"
	Action      string         `json:"action" gorm:"not null"`   // e.g., "create", "read", "update", "delete"
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	StatusID    *int16         `json:"status_id" gorm:"type:smallint;not null;default:1;index"`
}

type CreatePermissionRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Resource    string `json:"resource" validate:"required"`
	Action      string `json:"action" validate:"required"`
}

func (Permission) TableName() string {
	return "permissions"
}