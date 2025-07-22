package permission

import (
	"gorm.io/gorm"
)

type Repository interface {
	CreatePermission(permission *Permission) error
	GetAllPermissions() ([]Permission, error)
	GetPermissionByID(id uint) (*Permission, error)
	UpdatePermission(permission *Permission) error
	DeletePermission(id uint) error
	GetPermissionsByIDs(ids []uint) ([]Permission, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreatePermission(permission *Permission) error {
	return r.db.Create(permission).Error
}

func (r *repository) GetAllPermissions() ([]Permission, error) {
	var permissions []Permission
	err := r.db.Where("status_id = ?", 0).Find(&permissions).Error
	return permissions, err
}

func (r *repository) GetPermissionByID(id uint) (*Permission, error) {
	var permission Permission
	err := r.db.Where("id = ? AND status_id = ?", id, 0).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *repository) UpdatePermission(permission *Permission) error {
	return r.db.Save(permission).Error
}

func (r *repository) DeletePermission(id uint) error {
	return r.db.Model(&Permission{}).Where("id = ?", id).Update("status_id", 1).Error
}

func (r *repository) GetPermissionsByIDs(ids []uint) ([]Permission, error) {
	var permissions []Permission
	err := r.db.Where("id IN ? AND status_id = ?", ids, 0).Find(&permissions).Error
	return permissions, err
}