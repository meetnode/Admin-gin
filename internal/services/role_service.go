package services

import (
	"Admin-gin/internal/database"
	"Admin-gin/internal/models"
)

type RoleService interface {
	AddRole(role *models.Role) error
	GetRoles() ([]models.Role, error)
	AssignRoleToUser(userRole *models.UserHasRole) error
	AssignPermissionsToRole(roleID uint, permIDs []uint) error
	DeleteRole(id uint) error
}

type roleService struct {
	db database.Service
}

func NewRoleService() RoleService {
	return &roleService{
		db: database.New(),
	}
}

func (s *roleService) AddRole(role *models.Role) error {
	if err := s.db.GetDB().Create(role).Error; err != nil {
		return err
	}
	return nil
}

func (s *roleService) GetRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := s.db.GetDB().Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *roleService) AssignRoleToUser(userRole *models.UserHasRole) error {
	if err := s.db.GetDB().Create(userRole).Error; err != nil {
		return err
	}
	return nil
}

func (s *roleService) AssignPermissionsToRole(roleID uint, permIDs []uint) error {
	for _, pid := range permIDs {
		rp := models.RoleHasPermission{RoleID: roleID, PermissionID: pid}
		if err := s.db.GetDB().Create(&rp).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *roleService) DeleteRole(id uint) error {
	var role models.Role
	if err := s.db.GetDB().First(&role, id).Error; err != nil {
		return err
	}

	if err := s.db.GetDB().Delete(&role).Error; err != nil {
		return err
	}
	return nil
}
