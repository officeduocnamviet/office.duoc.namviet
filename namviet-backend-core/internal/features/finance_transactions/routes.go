package finance_transactions

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the finance transactions routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/finance-transactions")
	{
		group.GET("", GetAllFinanceTransactionsHandler)
		group.GET("/:id", GetFinanceTransactionHandler)
		group.POST("", CreateFinanceTransactionHandler)
		group.PUT("/:id", UpdateFinanceTransactionHandler)
		group.DELETE("/:id", DeleteFinanceTransactionHandler)
	}
}
