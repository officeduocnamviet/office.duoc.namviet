package users

import (
	"github.com/gin-gonic/gin"
	auth_middleware "github.com/namviet/backend-core/internal/middleware"
)

// RegisterRoutes registers endpoints for the users module
func RegisterRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", LoginHandler)
	}

	// Apply Auth Middleware to all users endpoints
	usersGroup := router.Group("/users")
	usersGroup.Use(auth_middleware.RequireAuth())
	{
		// Current user operations
		usersGroup.POST("/me/fcm-token", RegisterFCMTokenHandler)

		// users.read permission required to view users
		usersGroup.GET("", auth_middleware.RequirePermission("users.read"), GetUsersHandler)
		
		// users.write permission required to manage users
		usersGroup.POST("", auth_middleware.RequirePermission("users.write"), CreateUserHandler)
		usersGroup.PUT("/:id", auth_middleware.RequirePermission("users.write"), UpdateUserHandler)
		usersGroup.DELETE("/:id", auth_middleware.RequirePermission("users.write"), DeleteUserHandler)
	}
}
