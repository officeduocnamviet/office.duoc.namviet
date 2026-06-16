package manufacturers

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the manufacturers routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/manufacturers")
	{
		group.GET("", GetAllManufacturersHandler)
		group.GET("/:id", GetManufacturerHandler)
		group.POST("", CreateManufacturerHandler)
		group.PUT("/:id", UpdateManufacturerHandler)
	}
}
