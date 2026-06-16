package system_configs

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the system configs routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/system-configs")
	{
		group.GET("", GetAllSystemConfigsHandler)
		group.GET("/:key", GetSystemConfigHandler)
		group.POST("", CreateSystemConfigHandler)
		group.PUT("/:key", UpdateSystemConfigHandler)
		group.DELETE("/:key", DeleteSystemConfigHandler)
	}
}
