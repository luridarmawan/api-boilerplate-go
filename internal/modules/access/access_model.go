package access

import (
	"time"

	"apiserver/internal/modules/group"
	"apiserver/internal/utils"

	"gorm.io/gorm"
)

type User struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Email       string         `json:"email" gorm:"uniqueIndex;not null"`
	APIKey      string         `json:"-" gorm:"uniqueIndex;not null"`
	GroupID     *uint          `json:"group_id" gorm:"index"`
	Group       *group.Group   `json:"group,omitempty" gorm:"foreignKey:GroupID"`
	ExpiredDate *time.Time     `json:"expired_date" gorm:"index"`
	RateLimit   int            `json:"rate_limit" gorm:"not null;default:120"` // Requests per minute
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	StatusID    *int16         `json:"status_id" gorm:"type:smallint;not null;default:1;index"`
}

func (User) TableName() string {
	return "access"
}

// Implement User interface
func (u *User) GetID() string {
	return u.ID
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetGroup() *group.Group {
	return u.Group
}

// UpdateExpiredDateRequest is the request body for updating API key expiration date
type UpdateExpiredDateRequest struct {
	ExpiredDate *time.Time `json:"expired_date"`
}
// GetRateLimit returns the rate limit for this user
func (u *User) GetRateLimit() int {
	return u.RateLimit
}

// UpdateRateLimitRequest is the request body for updating API key rate limit
type UpdateRateLimitRequest struct {
	RateLimit int `json:"rate_limit" validate:"required,min=1"`
}

// BeforeCreate hook to generate UUIDv7 before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = utils.GenerateUUIDv7()
	}
	return nil
}