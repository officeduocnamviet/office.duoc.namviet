package employees

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the employees routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/employees")
	{
		group.GET("", GetAllEmployeesHandler)
		group.GET("/:id", GetEmployeeHandler)
		group.POST("", CreateEmployeeHandler)
		group.PUT("/:id", UpdateEmployeeHandler)
		group.DELETE("/:id", DeleteEmployeeHandler)
	}
}
