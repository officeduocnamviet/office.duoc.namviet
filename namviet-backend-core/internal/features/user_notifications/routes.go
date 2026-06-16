package user_notifications

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the user notifications and social routes
func RegisterRoutes(router *gin.RouterGroup) {
	fcmGroup := router.Group("/fcm-tokens")
	{
		fcmGroup.GET("", GetAllFCMTokensHandler)
		fcmGroup.GET("/:id", GetFCMTokenHandler)
		fcmGroup.POST("", CreateFCMTokenHandler)
		fcmGroup.PUT("/:id", UpdateFCMTokenHandler)
		fcmGroup.DELETE("/:id", DeleteFCMTokenHandler)
	}

	socialGroup := router.Group("/social-mappings")
	{
		socialGroup.GET("", GetAllSocialMappingsHandler)
		socialGroup.GET("/:id", GetSocialMappingHandler)
		socialGroup.POST("", CreateSocialMappingHandler)
		socialGroup.DELETE("/:id", DeleteSocialMappingHandler)
	}
}
