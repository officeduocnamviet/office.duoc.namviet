package knowledge_vectors

import "github.com/namviet/backend-core/internal/features/roles"

// Medical Knowledge Vector DTOs
type CreateMedicalKnowledgeVectorRequest struct {
	Title     string       `json:"title" binding:"required"`
	Content   string       `json:"content" binding:"required"`
	Embedding string       `json:"embedding"`
	Metadata  *roles.JSONB `json:"metadata"`
}

type UpdateMedicalKnowledgeVectorRequest struct {
	Title     *string      `json:"title"`
	Content   *string      `json:"content"`
	Embedding *string      `json:"embedding"`
	Metadata  *roles.JSONB `json:"metadata"`
}

// Product Vector DTOs
type CreateProductVectorRequest struct {
	ProductID *string      `json:"product_id"`
	Content   string       `json:"content" binding:"required"`
	Embedding string       `json:"embedding"`
	Metadata  *roles.JSONB `json:"metadata"`
}

type UpdateProductVectorRequest struct {
	ProductID *string      `json:"product_id"`
	Content   *string      `json:"content"`
	Embedding *string      `json:"embedding"`
	Metadata  *roles.JSONB `json:"metadata"`
}
