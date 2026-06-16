package payrolls

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the payrolls routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/payrolls")
	{
		group.GET("", GetAllPayrollsHandler)
		group.GET("/:id", GetPayrollHandler)
		group.POST("", CreatePayrollHandler)
		group.PUT("/:id", UpdatePayrollHandler)
		group.DELETE("/:id", DeletePayrollHandler)
	}
}
