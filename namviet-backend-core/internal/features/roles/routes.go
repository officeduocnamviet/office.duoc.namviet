package roles

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the roles routes
func RegisterRoutes(router *gin.RouterGroup) {
	rolesGroup := router.Group("/roles")
	{
		rolesGroup.GET("", GetAllRolesHandler)
		rolesGroup.GET("/:id", GetRoleHandler)
		rolesGroup.POST("", CreateRoleHandler)
		rolesGroup.PUT("/:id", UpdateRoleHandler)
		rolesGroup.DELETE("/:id", DeleteRoleHandler)
	}
}
