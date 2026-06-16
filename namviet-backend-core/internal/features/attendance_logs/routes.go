package attendance_logs

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the attendance logs routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/attendance-logs")
	{
		group.GET("", GetAllAttendanceLogsHandler)
		group.GET("/:id", GetAttendanceLogHandler)
		group.POST("", CreateAttendanceLogHandler)
		group.PUT("/:id", UpdateAttendanceLogHandler)
		group.DELETE("/:id", DeleteAttendanceLogHandler)
	}
}
