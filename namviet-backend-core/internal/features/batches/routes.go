package batches

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the batches routes
func RegisterRoutes(router *gin.RouterGroup) {
	// Nested under products for GET
	productsGroup := router.Group("/products")
	{
		productsGroup.GET("/:id/batches", GetBatchesByProductIDHandler)
	}

	// Direct access
	batchesGroup := router.Group("/batches")
	{
		batchesGroup.GET("/:id", GetBatchHandler)
		batchesGroup.POST("", CreateBatchHandler)
		batchesGroup.PUT("/:id", UpdateBatchHandler)
	}
}
