package server

import (
	"net/http"

	controller "Admin-gin/internal/controllers"
	middleware "Admin-gin/internal/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	{
		api := r.Group("/api")
		{
			api.POST("/login", controller.LoginHandler)
			api.POST("/register", controller.RegisterHandler)
			api.GET("/verify", controller.VerifyEmail)
			api.POST("/forgot-password", controller.ForgotPassword)
			api.POST("/reset-password", controller.ResetPassword)
		}
		{
			auth := api.Group("/")
			auth.Use(middleware.AuthMiddleware())
			{
				//Users
				userRoute := auth.Group("/users")

				userRoute.GET("/",
					middleware.HasPermission(s.db, "user.read"),
					controller.UserListing)

				userRoute.GET("/:id",
					middleware.HasPermission(s.db, "user.read"),
					controller.GetUserByID)

				userRoute.POST("/:id/assign-role",
					middleware.HasPermission(s.db, "role.assign"),
					controller.AssignRoleToUser)

				userRoute.PUT("/:id",
					middleware.HasPermission(s.db, "user.update"),
					controller.UpdateUser)

				userRoute.DELETE("/:id",
					middleware.HasPermission(s.db, "user.delete"),
					controller.DeleteUser)

				userRoute.PUT("/:id/password",
					middleware.HasPermission(s.db, "user.update"),
					controller.ChangePassword)
			}
			{
				//Permissions
				permissionRoute := auth.Group("/permissions")

				permissionRoute.GET("/",
					middleware.HasPermission(s.db, "permission.read"),
					controller.GetPermissions)

				permissionRoute.POST("/",
					middleware.HasPermission(s.db, "permission.create"),
					controller.CreatePermission)

				permissionRoute.DELETE("/:id",
					middleware.HasPermission(s.db, "permission.delete"),
					controller.DeletePermission)
			}
			{
				//Roles
				roleRoute := auth.Group("/roles")

				roleRoute.GET("/",
					middleware.HasPermission(s.db, "role.read"),
					controller.GetRoles)

				roleRoute.POST("/",
					middleware.HasPermission(s.db, "role.create"),
					controller.CreateRole)

				roleRoute.POST("/permissions",
					middleware.HasPermission(s.db, "role.update"),
					controller.AssignPermissionsToRole)
			}
		}
	}

	r.GET("/health", s.healthHandler)

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
