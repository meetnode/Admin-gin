package models

type UserHasRole struct {
	id     uint `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID uint `gorm:"not null" json:"role_id"`
	UserID uint `gorm:"not null" json:"user_id"`

	Role Role `gorm:"foreignKey:RoleID" json:"role"`
	User User `gorm:"foreignKey:UserID" json:"user"`
}
