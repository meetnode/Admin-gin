package controller

import (
	"Admin-gin/internal/models"
	"Admin-gin/internal/services"
	"Admin-gin/internal/utils"
	"fmt"
	"strconv"

	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginCred struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RolePermissionRequest struct {
	RoleID        uint   `json:"role_id" binding:"required"`
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

type UpdateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"new_password" binding:"required"`
}

// UserListing godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "users"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /users [get]
func UserListing(c *gin.Context) {
	userService := services.NewUserService()
	users, err := userService.GetAllUsers()
	// userId, isUserId := c.Get("userID")
	// fmt.Println("UserID from token:", userId, isUserId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"users": users})
}

// LoginHandler godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body LoginCred true "User credentials"
// @Success 200 {object} map[string]interface{} "token and user data"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var cred LoginCred
	if err := c.ShouldBindJSON(&cred); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userService := services.NewUserService()
	usr, err := userService.UserLogin(cred.Email, cred.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.CreateToken(usr)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	resp := make(map[string]interface{})
	userData := make(map[string]interface{})

	userData["id"] = usr.ID
	userData["name"] = usr.Name
	userData["email"] = usr.Email
	userData["status"] = usr.Status
	userData["created_at"] = usr.CreatedAt

	resp["token"] = token
	resp["user"] = userData

	c.JSON(http.StatusOK, resp)
}

// RegisterHandler godoc
// @Summary User registration
// @Description Register a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /register [post]
func RegisterHandler(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userService := services.NewUserService()
	userData, err := userService.GetUserByEmail(req.Email)
	fmt.Print(userData, "::userData")
	if err != nil {
		c.JSON(500, gin.H{"error": "somethings went wrong"})
		return
	} else if userData.Email == req.Email {
		c.JSON(400, gin.H{"error": "email already exists"})
		return
	}
	user, err := userService.AddUser(&req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	roleService := services.NewRoleService()
	roleService.AssignRoleToUser(&models.UserHasRole{
		UserID: user.ID,
		RoleID: 2,
	})

	c.JSON(200, gin.H{"message": "Registered successfully, check your email to verify"})
}

// VerifyEmail godoc
// @Summary Verify email address
// @Description Verify user email with token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /verify [get]
func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	decryptedEmail, err := utils.Decrypt(token)
	if err != nil {
		c.JSON(500, gin.H{"error": "email is not valid"})
		return
	}
	userService := services.NewUserService()
	user, err := userService.GetUserByEmail(decryptedEmail)
	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}

	user.Status = "active"
	if err = userService.UpdateUser(user); err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(200, gin.H{"message": "Email verified successfully"})
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update user data by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body UpdateUserRequest true "User update data"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		ID:   uint(id),
		Name: req.Name,
	}

	userService := services.NewUserService()
	if err := userService.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// CreatePermission godoc
// @Summary Create a new permission
// @Description Create a new permission
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param permission body models.Permission true "Permission data"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /permissions [post]
func CreatePermission(c *gin.Context) {
	var perm models.Permission
	if err := c.ShouldBindJSON(&perm); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	permissionService := services.NewPermissionService()
	if err := permissionService.AddPermission(&perm); err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(200, gin.H{"message": "Permission created successfully"})
}

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role body models.Role true "Role data"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /roles [post]
func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	roleService := services.NewRoleService()
	if err := roleService.AddRole(&role); err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(200, gin.H{"message": "Role created successfully"})
}

// AssignRoleToUser godoc
// @Summary Assign role to user
// @Description Assign a role to a specific user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param userRole body models.UserHasRole true "User role assignment"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /users/{id}/assign-role [post]
func AssignRoleToUser(c *gin.Context) {
	var userRole models.UserHasRole
	if err := c.ShouldBindJSON(&userRole); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	roleService := services.NewRoleService()
	if err := roleService.AssignRoleToUser(&userRole); err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(200, gin.H{"message": "Role assigned to user successfully"})
}

// AssignPermissionsToRole godoc
// @Summary Assign permissions to role
// @Description Assign multiple permissions to a specific role
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param rolePermission body RolePermissionRequest true "Role permission assignment"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /roles/permissions [post]
func AssignPermissionsToRole(c *gin.Context) {
	var rolePerm RolePermissionRequest
	if err := c.ShouldBindJSON(&rolePerm); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	roleService := services.NewRoleService()
	if err := roleService.AssignPermissionsToRole(rolePerm.RoleID, rolePerm.PermissionIDs); err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(200, gin.H{"message": "Permissions assigned to role successfully"})
}

// GetPermissions godoc
// @Summary Get all permissions
// @Description Get a list of all permissions
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Permission
// @Failure 500 {object} map[string]interface{} "error"
// @Router /permissions [get]
func GetPermissions(c *gin.Context) {
	permissionService := services.NewPermissionService()
	permissions, err := permissionService.GetPermissions()
	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(200, permissions)
}

// DeletePermission godoc
// @Summary Delete permission
// @Description Delete a permission by ID
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Permission ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /permissions/{id} [delete]
func DeletePermission(c *gin.Context) {
	permissionID := c.Param("id")
	id, err := strconv.ParseUint(permissionID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}
	permissionService := services.NewPermissionService()
	if err := permissionService.DeletePermission(uint(id)); err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(200, gin.H{"message": "Permission deleted successfully"})
}

// GetRoles godoc
// @Summary Get all roles
// @Description Get a list of all roles
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Role
// @Failure 500 {object} map[string]interface{} "error"
// @Router /roles [get]
func GetRoles(c *gin.Context) {
	roleService := services.NewRoleService()
	roles, err := roleService.GetRoles()
	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(200, roles)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user information by user ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	userService := services.NewUserService()
	user, err := userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(200, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	userService := services.NewUserService()
	if err := userService.DeleteUser(uint(id)); err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password for a specific user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param passwordData body ChangePasswordRequest true "Password change data"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /users/{id}/password [put]
func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.Param("id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	userService := services.NewUserService()
	if err := userService.ChangePassword(uint(id), req.OldPassword, req.NewPassword); err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(200, gin.H{"message": "Password changed successfully"})
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset email to user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param email body ForgotPasswordRequest true "Email address"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /forgot-password [post]
func ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := utils.Encrypt(req.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := utils.SendPasswordResetEmail(req.Email, token); err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(200, gin.H{"message": "Password reset email sent"})
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset user password with token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param resetData body ResetPasswordRequest true "Reset password data"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /reset-password [post]
func ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	email, err := utils.Decrypt(req.Token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userService := services.NewUserService()
	if err := userService.ResetPassword(email, req.Password); err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(200, gin.H{"message": "Password reset successfully"})
}

// DeleteRole godoc
// @Summary Delete role
// @Description Delete a role by ID
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /roles/{id} [delete]
func DeleteRole(c *gin.Context) {
	roleID := c.Param("id")
	id, err := strconv.ParseUint(roleID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}
	roleService := services.NewRoleService()
	if err := roleService.DeleteRole(uint(id)); err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(200, gin.H{"message": "Role deleted successfully"})
}
