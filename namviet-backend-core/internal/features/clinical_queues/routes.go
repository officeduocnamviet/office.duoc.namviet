package clinical_queues

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the clinical queues routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/clinical-queues")
	{
		group.GET("", GetAllClinicalQueuesHandler)
		group.GET("/:id", GetClinicalQueueHandler)
		group.POST("", CreateClinicalQueueHandler)
		group.PUT("/:id", UpdateClinicalQueueHandler)
		group.DELETE("/:id", DeleteClinicalQueueHandler)
	}
}
