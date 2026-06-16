package chats

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the chat routes
func RegisterRoutes(router *gin.RouterGroup) {
	sessionGroup := router.Group("/chat-sessions")
	{
		sessionGroup.GET("", GetAllChatSessionsHandler)
		sessionGroup.GET("/:id", GetChatSessionHandler)
		sessionGroup.POST("", CreateChatSessionHandler)
		sessionGroup.PUT("/:id", UpdateChatSessionHandler)
		sessionGroup.DELETE("/:id", DeleteChatSessionHandler)
		sessionGroup.GET("/:id/messages", GetMessagesBySessionIDHandler)
	}

	messageGroup := router.Group("/chat-messages")
	{
		messageGroup.POST("", CreateChatMessageHandler)
	}
}
