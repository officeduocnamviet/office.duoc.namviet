package audit_logs

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the audit logs routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/audit-logs")
	{
		group.GET("", GetAllAuditLogsHandler)
		group.GET("/:id", GetAuditLogHandler)
	}
}
