package fund_accounts

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the fund accounts routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/fund-accounts")
	{
		group.GET("", GetAllFundAccountsHandler)
		group.GET("/:id", GetFundAccountHandler)
		group.POST("", CreateFundAccountHandler)
		group.PUT("/:id", UpdateFundAccountHandler)
		group.DELETE("/:id", DeleteFundAccountHandler)
	}
}
