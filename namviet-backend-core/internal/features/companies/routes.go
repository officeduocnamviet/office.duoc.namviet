package companies

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the companies and branches routes
func RegisterRoutes(router *gin.RouterGroup) {
	companiesGroup := router.Group("/companies")
	{
		companiesGroup.GET("", GetAllCompaniesHandler)
		companiesGroup.GET("/:id", GetCompanyHandler)
		companiesGroup.POST("", CreateCompanyHandler)
		companiesGroup.PUT("/:id", UpdateCompanyHandler)
		companiesGroup.DELETE("/:id", DeleteCompanyHandler)
	}

	branchesGroup := router.Group("/branches")
	{
		branchesGroup.GET("", GetAllBranchesHandler)
		branchesGroup.GET("/:id", GetBranchHandler)
		branchesGroup.POST("", CreateBranchHandler)
		branchesGroup.PUT("/:id", UpdateBranchHandler)
		branchesGroup.DELETE("/:id", DeleteBranchHandler)
	}
}
