package configuration

import (
	"gorm.io/gorm"
)

type Repository interface {
	CreateConfiguration(configuration *Configuration) error
	GetAllConfigurations() ([]Configuration, error)
	GetConfigurationByID(id string) (*Configuration, error)
	UpdateConfiguration(configuration *Configuration) error
	SoftDeleteConfiguration(id string) error
	RestoreConfiguration(id string) error
	GetDeletedConfigurations() ([]Configuration, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateConfiguration(configuration *Configuration) error {
	return r.db.Create(configuration).Error
}

func (r *repository) GetAllConfigurations() ([]Configuration, error) {
	var configurations []Configuration
	err := r.db.Where("status_id = ?", 0).Find(&configurations).Error
	return configurations, err
}

func (r *repository) GetConfigurationByID(id string) (*Configuration, error) {
	var configuration Configuration
	err := r.db.Where("id = ? AND status_id = ?", id, 0).First(&configuration).Error
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

func (r *repository) UpdateConfiguration(configuration *Configuration) error {
	return r.db.Save(configuration).Error
}

func (r *repository) SoftDeleteConfiguration(id string) error {
	return r.db.Model(&Configuration{}).Where("id = ?", id).Update("status_id", 1).Error
}

func (r *repository) RestoreConfiguration(id string) error {
	return r.db.Model(&Configuration{}).Where("id = ?", id).Update("status_id", 0).Error
}

func (r *repository) GetDeletedConfigurations() ([]Configuration, error) {
	var configurations []Configuration
	err := r.db.Where("status_id = ?", 1).Find(&configurations).Error
	return configurations, err
}
