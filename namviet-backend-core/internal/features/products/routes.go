package products

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the products routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/products")
	{
		group.GET("", GetAllProductsHandler)
		group.GET("/:id", GetProductHandler)
		group.POST("", CreateProductHandler)
		group.PUT("/:id", UpdateProductHandler)
	}
}
