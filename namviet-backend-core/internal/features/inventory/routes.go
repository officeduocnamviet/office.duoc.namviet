package inventory

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the inventory routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/inventory")
	{
		group.GET("", GetInventoryHandler)
		group.POST("/transactions", CreateTransactionHandler)
	}
}
