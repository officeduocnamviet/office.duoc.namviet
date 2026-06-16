package warehouses

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the warehouses routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/warehouses")
	{
		group.GET("", GetAllWarehousesHandler)
		group.GET("/:id", GetWarehouseHandler)
		group.POST("", CreateWarehouseHandler)
		group.PUT("/:id", UpdateWarehouseHandler)
		group.DELETE("/:id", DeleteWarehouseHandler)
	}
}
