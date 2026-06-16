package orders

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the orders routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/orders")
	{
		group.GET("", GetAllOrdersHandler)
		group.GET("/:id", GetOrderHandler)
		group.POST("", CreateOrderHandler)
		group.PUT("/:id", UpdateOrderHandler)
		group.DELETE("/:id", DeleteOrderHandler)
	}
}
