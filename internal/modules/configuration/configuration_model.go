package configuration

import (
	"time"

	"apiserver/internal/utils"
	"gorm.io/gorm"
)

type Configuration struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey"`
	Key         string         `json:"key" gorm:"not null;uniqueIndex"`
	Value       string         `json:"value" gorm:"type:text"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	StatusID    *int16         `json:"status_id" gorm:"type:smallint;not null;default:1;index"`
}

type CreateConfigurationRequest struct {
	Key         string `json:"key" validate:"required"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type UpdateConfigurationRequest struct {
	Value       string `json:"value"`
	Description string `json:"description"`
}

func (Configuration) TableName() string {
	return "configurations"
}

// BeforeCreate hook to generate UUIDv7 before creating a new configuration
func (e *Configuration) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = utils.GenerateUUIDv7()
	}
	return nil
}
