package example

import (
	"time"

	"apiserver/internal/utils"
	"gorm.io/gorm"
)

type Example struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	StatusID    *int16         `json:"status_id" gorm:"type:smallint;not null;default:1;index"`
}

type CreateExampleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func (Example) TableName() string {
	return "examples"
}

// BeforeCreate hook to generate UUIDv7 before creating a new example
func (e *Example) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = utils.GenerateUUIDv7()
	}
	return nil
}