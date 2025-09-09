package services

import (
	"Admin-gin/internal/database"
	"Admin-gin/internal/models"
)

type PermissionService interface {
	AddPermission(permission *models.Permission) error
	GetPermissions() ([]models.Permission, error)
	DeletePermission(id uint) error
}

type permissionService struct {
	db database.Service
}

func NewPermissionService() PermissionService {
	return &permissionService{
		db: database.New(),
	}
}

func (s *permissionService) AddPermission(permission *models.Permission) error {
	if err := s.db.GetDB().Create(permission).Error; err != nil {
		return err
	}
	return nil
}

func (s *permissionService) GetPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := s.db.GetDB().Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *permissionService) DeletePermission(id uint) error {
	if err := s.db.GetDB().Delete(&models.Permission{}, id).Error; err != nil {
		return err
	}
	return nil
}
