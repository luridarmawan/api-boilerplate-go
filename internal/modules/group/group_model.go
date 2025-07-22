package group

import (
	"time"

	"apiserver/internal/modules/permission"

	"gorm.io/gorm"
)

type Group struct {
	ID          uint                    `json:"id" gorm:"primaryKey"`
	Name        string                  `json:"name" gorm:"uniqueIndex;not null"`
	Description string                  `json:"description"`
	Permissions []permission.Permission `json:"permissions" gorm:"many2many:group_permissions;"`
	CreatedAt   time.Time               `json:"-"`
	UpdatedAt   time.Time               `json:"-"`
	DeletedAt   gorm.DeletedAt          `json:"-" gorm:"index"`
	StatusID    *int16                  `json:"status_id" gorm:"type:smallint;not null;default:1;index"`
}

type CreateGroupRequest struct {
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description"`
	PermissionIDs []uint `json:"permission_ids"`
}

type UpdateGroupPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids"`
}

func (Group) TableName() string {
	return "groups"
}