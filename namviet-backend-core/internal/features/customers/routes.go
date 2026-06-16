package customers

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the customers routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/customers")
	{
		group.GET("", GetAllCustomersHandler)
		group.GET("/:id", GetCustomerHandler)
		group.POST("", CreateCustomerHandler)
		group.PUT("/:id", UpdateCustomerHandler)
		group.DELETE("/:id", DeleteCustomerHandler)
	}
}
