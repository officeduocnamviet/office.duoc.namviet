package chart_of_accounts

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the chart of accounts routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/chart-of-accounts")
	{
		group.GET("", GetAllChartOfAccountsHandler)
		group.GET("/:id", GetChartOfAccountHandler)
		group.POST("", CreateChartOfAccountHandler)
		group.PUT("/:id", UpdateChartOfAccountHandler)
		group.DELETE("/:id", DeleteChartOfAccountHandler)
	}
}
