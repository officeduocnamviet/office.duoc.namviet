package integrations

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the integrations routes
func RegisterRoutes(router *gin.RouterGroup) {
	connGroup := router.Group("/connections")
	{
		connGroup.GET("", GetAllConnectionsHandler)
		connGroup.GET("/:id", GetConnectionHandler)
		connGroup.POST("", CreateConnectionHandler)
		connGroup.PUT("/:id", UpdateConnectionHandler)
		connGroup.DELETE("/:id", DeleteConnectionHandler)
	}

	webhookGroup := router.Group("/webhook-logs")
	{
		webhookGroup.GET("", GetAllWebhookLogsHandler)
		webhookGroup.GET("/:id", GetWebhookLogHandler)
		webhookGroup.POST("", CreateWebhookLogHandler)
	}
}
