package training_courses

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the training courses routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/training-courses")
	{
		group.GET("", GetAllTrainingCoursesHandler)
		group.GET("/:id", GetTrainingCourseHandler)
		group.POST("", CreateTrainingCourseHandler)
		group.PUT("/:id", UpdateTrainingCourseHandler)
		group.DELETE("/:id", DeleteTrainingCourseHandler)
	}
}
