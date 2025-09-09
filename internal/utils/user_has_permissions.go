package utils

import (
	"Admin-gin/internal/models"

	"gorm.io/gorm"
)

func GetUserPermissions(db *gorm.DB, userID uint) ([]models.Permission, error) {
	var user models.User
	err := db.Preload("Roles.Permissions").First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	permMap := make(map[uint]models.Permission)
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			permMap[perm.ID] = perm
		}
	}

	perms := make([]models.Permission, 0, len(permMap))
	for _, p := range permMap {
		perms = append(perms, p)
	}

	return perms, nil
}
