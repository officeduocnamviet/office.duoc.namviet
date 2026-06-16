package companies

type CreateCompanyRequest struct {
	TaxCode                 string       `json:"tax_code" binding:"required"`
	Name                    string       `json:"name" binding:"required"`
	ShortName               *string      `json:"short_name"`
	Address                 string       `json:"address" binding:"required"`
	Phone                   string       `json:"phone" binding:"required"`
	Email                   *string      `json:"email"`
	LogoURL                 *string      `json:"logo_url"`
	RepresentativeName      *string      `json:"representative_name"`
	BusinessLicenseURL *[]string `json:"business_license_url"`
}

type UpdateCompanyRequest struct {
	TaxCode                 *string      `json:"tax_code"`
	Name                    *string      `json:"name"`
	ShortName               *string      `json:"short_name"`
	Address                 *string      `json:"address"`
	Phone                   *string      `json:"phone"`
	Email                   *string      `json:"email"`
	LogoURL                 *string      `json:"logo_url"`
	RepresentativeName      *string      `json:"representative_name"`
	BusinessLicenseURL *[]string `json:"business_license_url"`
	Status                  *string      `json:"status"`
}

type CreateBranchRequest struct {
	CompanyID string   `json:"company_id" binding:"required"`
	Code      string   `json:"code" binding:"required"`
	Name      string   `json:"name" binding:"required"`
	Address   string   `json:"address" binding:"required"`
	Phone     *string  `json:"phone"`
	ManagerID *string  `json:"manager_id"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type UpdateBranchRequest struct {
	CompanyID *string  `json:"company_id"`
	Code      *string  `json:"code"`
	Name      *string  `json:"name"`
	Address   *string  `json:"address"`
	Phone     *string  `json:"phone"`
	ManagerID *string  `json:"manager_id"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Status    *string  `json:"status"`
}
