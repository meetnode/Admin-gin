package models

import (
	"time"

	"gorm.io/gorm"
)

type RoleHasPermission struct {
	ID           uint `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID       uint `gorm:"not null" json:"role_id"`
	PermissionID uint `gorm:"not null" json:"permission_id"`

	Role       Role       `gorm:"foreignKey:RoleID" json:"role"`
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission"`

	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
