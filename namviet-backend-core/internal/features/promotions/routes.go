package promotions

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the promotions routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/promotions")
	{
		group.GET("", GetAllPromotionsHandler)
		group.GET("/:id", GetPromotionHandler)
		group.POST("", CreatePromotionHandler)
		group.PUT("/:id", UpdatePromotionHandler)
		group.DELETE("/:id", DeletePromotionHandler)
	}
}
