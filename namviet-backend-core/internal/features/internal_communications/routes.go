package internal_communications

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the internal communications routes
func RegisterRoutes(router *gin.RouterGroup) {
	channelGroup := router.Group("/internal-channels")
	{
		channelGroup.GET("", GetAllInternalChannelsHandler)
		channelGroup.GET("/:id", GetInternalChannelHandler)
		channelGroup.POST("", CreateInternalChannelHandler)
		channelGroup.PUT("/:id", UpdateInternalChannelHandler)
		channelGroup.DELETE("/:id", DeleteInternalChannelHandler)
		
		// Messages for a channel
		channelGroup.GET("/:id/messages", GetMessagesByChannelIDHandler)
	}

	msgGroup := router.Group("/internal-messages")
	{
		msgGroup.POST("", CreateInternalMessageHandler)
		msgGroup.DELETE("/:id", DeleteInternalMessageHandler)
	}
}
