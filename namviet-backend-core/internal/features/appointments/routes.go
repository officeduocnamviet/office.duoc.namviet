package appointments

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the appointments routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/appointments")
	{
		group.GET("", GetAllAppointmentsHandler)
		group.GET("/:id", GetAppointmentHandler)
		group.POST("", CreateAppointmentHandler)
		group.PUT("/:id", UpdateAppointmentHandler)
		group.DELETE("/:id", DeleteAppointmentHandler)
	}
}
