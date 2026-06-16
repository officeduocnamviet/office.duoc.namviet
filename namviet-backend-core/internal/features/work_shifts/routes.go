package work_shifts

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the work shifts routes
func RegisterRoutes(router *gin.RouterGroup) {
	shiftsGroup := router.Group("/work-shifts")
	{
		shiftsGroup.GET("", GetAllWorkShiftsHandler)
		shiftsGroup.GET("/:id", GetWorkShiftHandler)
		shiftsGroup.POST("", CreateWorkShiftHandler)
		shiftsGroup.PUT("/:id", UpdateWorkShiftHandler)
		shiftsGroup.DELETE("/:id", DeleteWorkShiftHandler)
	}

	assignGroup := router.Group("/shift-assignments")
	{
		assignGroup.GET("", GetAllShiftAssignmentsHandler)
		assignGroup.GET("/:id", GetShiftAssignmentHandler)
		assignGroup.POST("", CreateShiftAssignmentHandler)
		assignGroup.PUT("/:id", UpdateShiftAssignmentHandler)
		assignGroup.DELETE("/:id", DeleteShiftAssignmentHandler)
	}

	handoverGroup := router.Group("/shift-handovers")
	{
		handoverGroup.GET("", GetAllShiftHandoversHandler)
		handoverGroup.GET("/:id", GetShiftHandoverHandler)
		handoverGroup.POST("", CreateShiftHandoverHandler)
		handoverGroup.PUT("/:id", UpdateShiftHandoverHandler)
		handoverGroup.DELETE("/:id", DeleteShiftHandoverHandler)
	}
}
