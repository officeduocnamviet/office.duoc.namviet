package ai_agent_memories

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the ai agent memories routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai-agent-memories")
	{
		group.GET("", GetAllAIAgentMemoriesHandler)
		group.GET("/:id", GetAIAgentMemoryHandler)
		group.POST("", CreateAIAgentMemoryHandler)
		group.PUT("/:id", UpdateAIAgentMemoryHandler)
		group.DELETE("/:id", DeleteAIAgentMemoryHandler)
	}
}
