package agent_workflows

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the agent workflows routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/agent-workflows")
	{
		group.GET("", GetAllAgentWorkflowsHandler)
		group.GET("/:id", GetAgentWorkflowHandler)
		group.POST("", CreateAgentWorkflowHandler)
		group.PUT("/:id", UpdateAgentWorkflowHandler)
		group.DELETE("/:id", DeleteAgentWorkflowHandler)
	}
}
