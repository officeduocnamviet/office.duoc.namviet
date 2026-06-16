package companies

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// Company represents the companies table
type Company struct {
	ID                      string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TaxCode                 string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"tax_code"`
	Name                    string      `gorm:"type:varchar(255);not null" json:"name"`
	ShortName               *string     `gorm:"type:varchar(100)" json:"short_name,omitempty"`
	Address                 string      `gorm:"type:text;not null" json:"address"`
	Phone                   string      `gorm:"type:varchar(20);not null" json:"phone"`
	Email                   *string     `gorm:"type:varchar(100)" json:"email,omitempty"`
	LogoURL                 *string     `gorm:"type:text" json:"logo_url,omitempty"`
	RepresentativeName      *string     `gorm:"type:varchar(100)" json:"representative_name,omitempty"`
	BusinessImageLicenseURL roles.JSONB `gorm:"type:jsonb" json:"business_image_license_url,omitempty"`
	Status                  string      `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt               *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt               *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt               *time.Time  `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}

// Branch represents the branches table
type Branch struct {
	ID         string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CompanyID  string     `gorm:"type:uuid;not null;index" json:"company_id"`
	Code       string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name       string     `gorm:"type:varchar(255);not null" json:"name"`
	Address    string     `gorm:"type:text;not null" json:"address"`
	Phone      *string    `gorm:"type:varchar(20)" json:"phone,omitempty"`
	ManagerID  *string    `gorm:"type:uuid" json:"manager_id,omitempty"`
	Latitude   *float64   `gorm:"type:numeric(10,7)" json:"latitude,omitempty"`
	Longitude  *float64   `gorm:"type:numeric(10,7)" json:"longitude,omitempty"`
	Status     string     `gorm:"type:varchar(20);default:'active';index" json:"status"`
	CreatedAt  *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt  *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt  *time.Time `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
