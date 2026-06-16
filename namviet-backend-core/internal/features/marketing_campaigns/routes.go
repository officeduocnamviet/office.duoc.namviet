package marketing_campaigns

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the marketing campaigns routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/marketing-campaigns")
	{
		group.GET("", GetAllMarketingCampaignsHandler)
		group.GET("/:id", GetMarketingCampaignHandler)
		group.POST("", CreateMarketingCampaignHandler)
		group.PUT("/:id", UpdateMarketingCampaignHandler)
		group.DELETE("/:id", DeleteMarketingCampaignHandler)
	}
}
