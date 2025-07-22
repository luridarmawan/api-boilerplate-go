package audit

import (
	"time"
)

type AuditLog struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	UserID         *string   `json:"user_id" gorm:"type:uuid;index"`
	UserEmail      string    `json:"user_email" gorm:"index"`
	APIKey         string    `json:"api_key" gorm:"index"`
	Method         string    `json:"method" gorm:"not null"`
	Path           string    `json:"path" gorm:"not null;index"`
	StatusCode     int       `json:"status_code" gorm:"not null;index"`
	RequestHeaders string    `json:"request_headers" gorm:"type:text"`
	RequestBody    string    `json:"request_body" gorm:"type:text"`
	ResponseBody   string    `json:"response_body" gorm:"type:text"`
	ResponseTime   int64     `json:"response_time"` // in milliseconds
	IPAddress      string    `json:"ip_address" gorm:"index"`
	UserAgent      string    `json:"user_agent"`
	CreatedAt      time.Time `json:"created_at" gorm:"index"`
	UpdatedAt      time.Time `json:"-"`
	StatusID       *int16    `json:"status_id" gorm:"type:smallint;not null;default:1;index"`
}

type AuditLogResponse struct {
	ID           uint      `json:"id"`
	UserEmail    string    `json:"user_email"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int64     `json:"response_time"`
	IPAddress    string    `json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuditLogFilter struct {
	UserEmail  string `json:"user_email"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	StatusCode int    `json:"status_code"`
	DateFrom   string `json:"date_from"` // YYYY-MM-DD format
	DateTo     string `json:"date_to"`   // YYYY-MM-DD format
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}