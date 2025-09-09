package middleware

import (
	"Admin-gin/internal/database"
	"Admin-gin/internal/utils"

	"github.com/gin-gonic/gin"
)

// HasPermission middleware checks if the user has the required permission
func HasPermission(db database.Service, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Get user's permissions
		permissions, err := utils.GetUserPermissions(db.GetDB(), uint(userID.(float64)))
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to get user permissions"})
			c.Abort()
			return
		}

		// Check if user has the required permission
		hasPermission := false
		for _, p := range permissions {
			if p.Name == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(403, gin.H{"error": "forbidden: insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
