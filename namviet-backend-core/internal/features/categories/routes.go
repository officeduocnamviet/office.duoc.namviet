package categories

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the categories routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/categories")
	{
		group.GET("", GetAllCategoriesHandler)
		group.GET("/:id", GetCategoryHandler)
		group.POST("", CreateCategoryHandler)
		group.PUT("/:id", UpdateCategoryHandler)
		group.DELETE("/:id", DeleteCategoryHandler)
	}
}
