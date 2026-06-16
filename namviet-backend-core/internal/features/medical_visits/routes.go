package medical_visits

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the medical visits routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/medical-visits")
	{
		group.GET("", GetAllMedicalVisitsHandler)
		group.GET("/:id", GetMedicalVisitHandler)
		group.POST("", CreateMedicalVisitHandler)
		group.PUT("/:id", UpdateMedicalVisitHandler)
		group.DELETE("/:id", DeleteMedicalVisitHandler)
	}
}
