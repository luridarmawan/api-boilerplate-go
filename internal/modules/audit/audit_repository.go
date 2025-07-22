package audit

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateAuditLog(log *AuditLog) error
	GetAuditLogs(filter AuditLogFilter) ([]AuditLogResponse, int64, error)
	GetAuditLogByID(id uint) (*AuditLog, error)
	DeleteOldLogs(days int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateAuditLog(log *AuditLog) error {
	return r.db.Create(log).Error
}

func (r *repository) GetAuditLogs(filter AuditLogFilter) ([]AuditLogResponse, int64, error) {
	var logs []AuditLogResponse
	var total int64

	query := r.db.Model(&AuditLog{})

	// Apply base filter for active records
	query = query.Where("status_id = ?", 0)

	// Apply filters
	if filter.UserEmail != "" {
		query = query.Where("user_email ILIKE ?", "%"+filter.UserEmail+"%")
	}
	if filter.Method != "" {
		query = query.Where("method = ?", filter.Method)
	}
	if filter.Path != "" {
		query = query.Where("path ILIKE ?", "%"+filter.Path+"%")
	}
	if filter.StatusCode != 0 {
		query = query.Where("status_code = ?", filter.StatusCode)
	}
	if filter.DateFrom != "" {
		if dateFrom, err := time.Parse("2006-01-02", filter.DateFrom); err == nil {
			query = query.Where("created_at >= ?", dateFrom)
		}
	}
	if filter.DateTo != "" {
		if dateTo, err := time.Parse("2006-01-02", filter.DateTo); err == nil {
			query = query.Where("created_at <= ?", dateTo.Add(24*time.Hour))
		}
	}

	// Count total records
	query.Count(&total)

	// Apply pagination
	if filter.Limit == 0 {
		filter.Limit = 50 // default limit
	}
	if filter.Limit > 1000 {
		filter.Limit = 1000 // max limit
	}

	err := query.Select("id, user_email, method, path, status_code, response_time, ip_address, created_at").
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&logs).Error

	return logs, total, err
}

func (r *repository) GetAuditLogByID(id uint) (*AuditLog, error) {
	var log AuditLog
	err := r.db.Where("id = ? AND status_id = ?", id, 0).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *repository) DeleteOldLogs(days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return r.db.Model(&AuditLog{}).Where("created_at < ? AND status_id = ?", cutoffDate, 0).Update("status_id", 1).Error
}