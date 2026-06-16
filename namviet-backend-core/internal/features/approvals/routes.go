package approvals

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the approvals routes
func RegisterRoutes(router *gin.RouterGroup) {
	reqGroup := router.Group("/approval-requests")
	{
		reqGroup.GET("", GetAllApprovalRequestsHandler)
		reqGroup.GET("/:id", GetApprovalRequestHandler)
		reqGroup.POST("", CreateApprovalRequestHandler)
		reqGroup.PUT("/:id", UpdateApprovalRequestHandler)
		reqGroup.DELETE("/:id", DeleteApprovalRequestHandler)
		reqGroup.GET("/:id/steps", GetStepsByRequestIDHandler)
	}

	stepGroup := router.Group("/approval-steps")
	{
		stepGroup.POST("", CreateApprovalStepHandler)
		stepGroup.PUT("/:id", UpdateApprovalStepHandler)
	}
}
