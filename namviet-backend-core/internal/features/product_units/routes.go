package product_units

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the product_units routes
func RegisterRoutes(router *gin.RouterGroup) {
	// Nested under products for POST
	productsGroup := router.Group("/products")
	{
		productsGroup.POST("/:product_id/units", CreateProductUnitHandler)
	}

	// Direct access for PUT and DELETE
	unitsGroup := router.Group("/products/units")
	{
		unitsGroup.PUT("/:unit_id", UpdateProductUnitHandler)
		unitsGroup.DELETE("/:unit_id", DeleteProductUnitHandler)
	}
}
