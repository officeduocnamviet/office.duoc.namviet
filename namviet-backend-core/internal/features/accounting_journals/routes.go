package accounting_journals

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the accounting journals routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/accounting-journals")
	{
		group.GET("", GetAllAccountingJournalsHandler)
		group.GET("/:id", GetAccountingJournalHandler)
		group.POST("", CreateAccountingJournalHandler)
		group.PUT("/:id", UpdateAccountingJournalHandler)
		group.DELETE("/:id", DeleteAccountingJournalHandler)
	}
}
