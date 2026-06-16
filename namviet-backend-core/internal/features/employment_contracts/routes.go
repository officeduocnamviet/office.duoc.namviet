package employment_contracts

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the employment contracts routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/employment-contracts")
	{
		group.GET("", GetAllEmploymentContractsHandler)
		group.GET("/:id", GetEmploymentContractHandler)
		group.POST("", CreateEmploymentContractHandler)
		group.PUT("/:id", UpdateEmploymentContractHandler)
		group.DELETE("/:id", DeleteEmploymentContractHandler)
	}
}
