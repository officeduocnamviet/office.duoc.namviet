package time_attendance

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the time attendance routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/time-attendance")
	{
		group.GET("", GetAllTimeAttendancesHandler)
		group.GET("/:id", GetTimeAttendanceHandler)
		group.POST("", CreateTimeAttendanceHandler)
		group.PUT("/:id", UpdateTimeAttendanceHandler)
		group.DELETE("/:id", DeleteTimeAttendanceHandler)
	}
}
