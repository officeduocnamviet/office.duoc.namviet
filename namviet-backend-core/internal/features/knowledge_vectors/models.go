package knowledge_vectors

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// MedicalKnowledgeVector represents the medical_knowledge_vectors table
type MedicalKnowledgeVector struct {
	ID        string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title     string       `gorm:"type:text;not null" json:"title"`
	Content   string       `gorm:"type:text;not null" json:"content"`
	Embedding string       `gorm:"type:vector" json:"embedding,omitempty"` // simplified
	Metadata  *roles.JSONB `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt *time.Time   `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}

// ProductVector represents the product_vectors table
type ProductVector struct {
	ID        string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ProductID *string      `gorm:"type:uuid" json:"product_id,omitempty"`
	Content   string       `gorm:"type:text;not null" json:"content"`
	Embedding string       `gorm:"type:vector" json:"embedding,omitempty"` // simplified
	Metadata  *roles.JSONB `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt *time.Time   `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
