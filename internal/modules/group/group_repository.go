package group

import (
	"apiserver/internal/modules/permission"

	"gorm.io/gorm"
)

type Repository interface {
	CreateGroup(group *Group) error
	GetAllGroups() ([]Group, error)
	GetGroupByID(id uint) (*Group, error)
	GetGroupWithPermissions(id uint) (*Group, error)
	UpdateGroup(group *Group) error
	DeleteGroup(id uint) error
	UpdateGroupPermissions(groupID uint, permissionIDs []uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateGroup(group *Group) error {
	return r.db.Create(group).Error
}

func (r *repository) GetAllGroups() ([]Group, error) {
	var groups []Group
	err := r.db.Preload("Permissions", "status_id = ?", 0).Where("status_id = ?", 0).Find(&groups).Error
	return groups, err
}

func (r *repository) GetGroupByID(id uint) (*Group, error) {
	var group Group
	err := r.db.Where("id = ? AND status_id = ?", id, 0).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *repository) GetGroupWithPermissions(id uint) (*Group, error) {
	var group Group
	err := r.db.Preload("Permissions", "status_id = ?", 0).Where("id = ? AND status_id = ?", id, 0).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *repository) UpdateGroup(group *Group) error {
	return r.db.Save(group).Error
}

func (r *repository) DeleteGroup(id uint) error {
	return r.db.Model(&Group{}).Where("id = ?", id).Update("status_id", 1).Error
}

func (r *repository) UpdateGroupPermissions(groupID uint, permissionIDs []uint) error {
	var group Group
	if err := r.db.Where("id = ? AND status_id = ?", groupID, 0).First(&group).Error; err != nil {
		return err
	}

	// Clear existing permissions
	if err := r.db.Model(&group).Association("Permissions").Clear(); err != nil {
		return err
	}

	// Add new permissions if any
	if len(permissionIDs) > 0 {
		var permissions []permission.Permission
		if err := r.db.Where("id IN ? AND status_id = ?", permissionIDs, 0).Find(&permissions).Error; err != nil {
			return err
		}
		return r.db.Model(&group).Association("Permissions").Append(permissions)
	}

	return nil
}