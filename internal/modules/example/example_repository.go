package example

import (
	"gorm.io/gorm"
)

type Repository interface {
	CreateExample(example *Example) error
	GetAllExamples() ([]Example, error)
	GetExampleByID(id string) (*Example, error)
	UpdateExample(example *Example) error
	SoftDeleteExample(id string) error
	RestoreExample(id string) error
	GetDeletedExamples() ([]Example, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateExample(example *Example) error {
	return r.db.Create(example).Error
}

func (r *repository) GetAllExamples() ([]Example, error) {
	var examples []Example
	err := r.db.Where("status_id = ?", 0).Find(&examples).Error
	return examples, err
}

func (r *repository) GetExampleByID(id string) (*Example, error) {
	var example Example
	err := r.db.Where("id = ? AND status_id = ?", id, 0).First(&example).Error
	if err != nil {
		return nil, err
	}
	return &example, nil
}

func (r *repository) UpdateExample(example *Example) error {
	return r.db.Save(example).Error
}

func (r *repository) SoftDeleteExample(id string) error {
	return r.db.Model(&Example{}).Where("id = ?", id).Update("status_id", 1).Error
}

func (r *repository) RestoreExample(id string) error {
	return r.db.Model(&Example{}).Where("id = ?", id).Update("status_id", 0).Error
}

func (r *repository) GetDeletedExamples() ([]Example, error) {
	var examples []Example
	err := r.db.Where("status_id = ?", 1).Find(&examples).Error
	return examples, err
}