package knowledge_vectors

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the knowledge vectors routes
func RegisterRoutes(router *gin.RouterGroup) {
	medicalGroup := router.Group("/medical-knowledge-vectors")
	{
		medicalGroup.GET("", GetAllMedicalKnowledgeVectorsHandler)
		medicalGroup.GET("/:id", GetMedicalKnowledgeVectorHandler)
		medicalGroup.POST("", CreateMedicalKnowledgeVectorHandler)
		medicalGroup.PUT("/:id", UpdateMedicalKnowledgeVectorHandler)
		medicalGroup.DELETE("/:id", DeleteMedicalKnowledgeVectorHandler)
	}

	productGroup := router.Group("/product-vectors")
	{
		productGroup.GET("", GetAllProductVectorsHandler)
		productGroup.GET("/:id", GetProductVectorHandler)
		productGroup.POST("", CreateProductVectorHandler)
		productGroup.PUT("/:id", UpdateProductVectorHandler)
		productGroup.DELETE("/:id", DeleteProductVectorHandler)
	}
}
