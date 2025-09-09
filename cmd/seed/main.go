package main

import (
	"Admin-gin/internal/models"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func main() {
	var (
		database = getEnv("BLUEPRINT_DB_DATABASE", "admin-gin")
		password = getEnv("BLUEPRINT_DB_PASSWORD", "abcd")
		username = getEnv("BLUEPRINT_DB_USERNAME", "postgres")
		port     = getEnv("BLUEPRINT_DB_PORT", "5432")
		host     = getEnv("BLUEPRINT_DB_HOST", "localhost")
		schema   = getEnv("BLUEPRINT_DB_SCHEMA", "public")
	)

	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s",
		host, username, password, database, port, schema,
	)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserHasRole{},
		&models.RoleHasPermission{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if err := seedDatabase(db); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	fmt.Println("Database seeded successfully!")
}

func seedDatabase(db *gorm.DB) error {
	permissions := []models.Permission{
		{Name: "user.create"},
		{Name: "user.read"},
		{Name: "user.update"},
		{Name: "user.delete"},
		{Name: "role.create"},
		{Name: "role.read"},
		{Name: "role.update"},
		{Name: "role.delete"},
		{Name: "role.assign"},
		{Name: "permission.create"},
		{Name: "permission.read"},
		{Name: "permission.update"},
		{Name: "permission.delete"},
		{Name: "system.admin"},
		{Name: "system.manage"},
	}

	userPermissions := []models.Permission{
		{Name: "user.read"},
		{Name: "role.read"},
		{Name: "permission.read"},
	}

	for i := range permissions {
		permissions[i].CreatedAt = time.Now()
		var existingPermission models.Permission
		if err := db.Where("name = ?", permissions[i].Name).First(&existingPermission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&permissions[i]).Error; err != nil {
					return fmt.Errorf("failed to create permission %s: %v", permissions[i].Name, err)
				}
				fmt.Printf("Created permission: %s\n", permissions[i].Name)
			} else {
				return fmt.Errorf("error checking permission %s: %v", permissions[i].Name, err)
			}
		} else {
			permissions[i] = existingPermission
			fmt.Printf("Permission already exists: %s\n", permissions[i].Name)
		}
	}

	superAdminRole := models.Role{
		Name:      "super_admin",
		CreatedAt: time.Now(),
	}

	userRole := models.Role{
		Name:      "user",
		CreatedAt: time.Now(),
	}

	var existingRole models.Role
	if err := db.Where("name = ?", superAdminRole.Name).First(&existingRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&superAdminRole).Error; err != nil {
				return fmt.Errorf("failed to create super admin role: %v", err)
			}

			fmt.Println("Created role: super_admin")
		} else {
			return fmt.Errorf("error checking super admin role: %v", err)
		}
	} else {
		superAdminRole = existingRole
		fmt.Println("Role already exists: super_admin")
	}
	if err := db.Where("name = ?", userRole.Name).First(&existingRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&userRole).Error; err != nil {
				return fmt.Errorf("failed to create user role: %v", err)
			}
			fmt.Println("Created role: user")
		} else {
			return fmt.Errorf("error checking user role: %v", err)
		}
	} else {
		userRole = existingRole
		fmt.Println("Role already exists: user")
	}

	for _, permission := range permissions {
		var existingRolePermission models.RoleHasPermission
		if err := db.Where("role_id = ? AND permission_id = ?", superAdminRole.ID, permission.ID).First(&existingRolePermission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				roleHasPermission := models.RoleHasPermission{
					RoleID:       superAdminRole.ID,
					PermissionID: permission.ID,
					CreatedAt:    time.Now(),
				}
				if err := db.Create(&roleHasPermission).Error; err != nil {
					return fmt.Errorf("failed to assign permission %s to super admin role: %v", permission.Name, err)
				}
				fmt.Printf("Assigned permission %s to super_admin role\n", permission.Name)
			} else {
				return fmt.Errorf("error checking role permission assignment: %v", err)
			}
		} else {
			fmt.Printf("Permission %s already assigned to super_admin role\n", permission.Name)
		}
	}

	for _, usrPermission := range userPermissions {
		fmt.Print(usrPermission, " UserPermissions")
		var existingRolePermission models.RoleHasPermission
		if err := db.Where("role_id = ? AND permission_id = ?", userRole.ID, usrPermission.ID).First(&existingRolePermission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				roleHasPermission := models.RoleHasPermission{
					RoleID:       userRole.ID,
					PermissionID: usrPermission.ID,
					CreatedAt:    time.Now(),
				}
				if err := db.Create(&roleHasPermission).Error; err != nil {
					return fmt.Errorf("failed to assign permission %s to user role: %v", usrPermission.Name, err)
				}
				fmt.Printf("Assigned permission %s to user role\n", usrPermission.Name)
			} else {
				return fmt.Errorf("error checking role permission assignment: %v", err)
			}
		} else {
			fmt.Printf("Permission %s already assigned to user role\n", usrPermission.Name)
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("SuperAdmin123!"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	superAdminUser := models.User{
		Name:      "Super Administrator",
		Email:     "superadmin@example.com",
		Password:  string(hashedPassword),
		Status:    "active",
		CreatedAt: time.Now(),
	}

	var existingUser models.User
	if err := db.Where("email = ?", superAdminUser.Email).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&superAdminUser).Error; err != nil {
				return fmt.Errorf("failed to create super admin user: %v", err)
			}
			fmt.Println("Created user: Super Administrator")
		} else {
			return fmt.Errorf("error checking super admin user: %v", err)
		}
	} else {
		superAdminUser = existingUser
		fmt.Println("User already exists: Super Administrator")
	}

	var existingUserRole models.UserHasRole
	if err := db.Where("user_id = ? AND role_id = ?", superAdminUser.ID, superAdminRole.ID).First(&existingUserRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			userHasRole := models.UserHasRole{
				UserID: superAdminUser.ID,
				RoleID: superAdminRole.ID,
			}
			if err := db.Create(&userHasRole).Error; err != nil {
				return fmt.Errorf("failed to assign super admin role to user: %v", err)
			}
			fmt.Println("Assigned super_admin role to Super Administrator user")
		} else {

			return fmt.Errorf("error checking user role assignment: %v", err)
		}
	} else {
		fmt.Println("Super admin role already assigned to Super Administrator user")
	}

	additionalRoles := []models.Role{
		{Name: "admin", CreatedAt: time.Now()},
		{Name: "editor", CreatedAt: time.Now()},
		{Name: "viewer", CreatedAt: time.Now()},
	}

	for i := range additionalRoles {
		var existingAdditionalRole models.Role
		if err := db.Where("name = ?", additionalRoles[i].Name).First(&existingAdditionalRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&additionalRoles[i]).Error; err != nil {
					return fmt.Errorf("failed to create role %s: %v", additionalRoles[i].Name, err)
				}
				fmt.Printf("Created role: %s\n", additionalRoles[i].Name)
			} else {
				return fmt.Errorf("error checking role %s: %v", additionalRoles[i].Name, err)
			}
		} else {
			fmt.Printf("Role already exists: %s\n", additionalRoles[i].Name)
		}
	}

	return nil
}
