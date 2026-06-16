package companies

import "time"

// Companies
func GetAllCompaniesService() ([]Company, error) {
	return GetAllCompanies()
}

func GetCompanyByIDService(id string) (*Company, error) {
	return GetCompanyByID(id)
}

func CreateCompanyService(req CreateCompanyRequest) (*Company, error) {
	company := &Company{
		TaxCode:            req.TaxCode,
		Name:               req.Name,
		ShortName:          req.ShortName,
		Address:            req.Address,
		Phone:              req.Phone,
		Email:              req.Email,
		LogoURL:            req.LogoURL,
		RepresentativeName: req.RepresentativeName,
		Status:             "active",
	}

	if req.BusinessImageLicenseURL != nil {
		company.BusinessImageLicenseURL = *req.BusinessImageLicenseURL
	}

	if err := CreateCompany(company); err != nil {
		return nil, err
	}
	return company, nil
}

func UpdateCompanyService(id string, req UpdateCompanyRequest) (*Company, error) {
	company, err := GetCompanyByID(id)
	if err != nil {
		return nil, err
	}

	if req.TaxCode != nil {
		company.TaxCode = *req.TaxCode
	}
	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.ShortName != nil {
		company.ShortName = req.ShortName
	}
	if req.Address != nil {
		company.Address = *req.Address
	}
	if req.Phone != nil {
		company.Phone = *req.Phone
	}
	if req.Email != nil {
		company.Email = req.Email
	}
	if req.LogoURL != nil {
		company.LogoURL = req.LogoURL
	}
	if req.RepresentativeName != nil {
		company.RepresentativeName = req.RepresentativeName
	}
	if req.BusinessImageLicenseURL != nil {
		company.BusinessImageLicenseURL = *req.BusinessImageLicenseURL
	}
	if req.Status != nil {
		company.Status = *req.Status
	}

	now := time.Now()
	company.UpdatedAt = &now

	if err := UpdateCompany(company); err != nil {
		return nil, err
	}
	return company, nil
}

func DeleteCompanyService(id string) error {
	return DeleteCompany(id)
}

// Branches
func GetAllBranchesService() ([]Branch, error) {
	return GetAllBranches()
}

func GetBranchByIDService(id string) (*Branch, error) {
	return GetBranchByID(id)
}

func CreateBranchService(req CreateBranchRequest) (*Branch, error) {
	branch := &Branch{
		CompanyID: req.CompanyID,
		Code:      req.Code,
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		ManagerID: req.ManagerID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Status:    "active",
	}

	if err := CreateBranch(branch); err != nil {
		return nil, err
	}
	return branch, nil
}

func UpdateBranchService(id string, req UpdateBranchRequest) (*Branch, error) {
	branch, err := GetBranchByID(id)
	if err != nil {
		return nil, err
	}

	if req.CompanyID != nil {
		branch.CompanyID = *req.CompanyID
	}
	if req.Code != nil {
		branch.Code = *req.Code
	}
	if req.Name != nil {
		branch.Name = *req.Name
	}
	if req.Address != nil {
		branch.Address = *req.Address
	}
	if req.Phone != nil {
		branch.Phone = req.Phone
	}
	if req.ManagerID != nil {
		branch.ManagerID = req.ManagerID
	}
	if req.Latitude != nil {
		branch.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		branch.Longitude = req.Longitude
	}
	if req.Status != nil {
		branch.Status = *req.Status
	}

	now := time.Now()
	branch.UpdatedAt = &now

	if err := UpdateBranch(branch); err != nil {
		return nil, err
	}
	return branch, nil
}

func DeleteBranchService(id string) error {
	return DeleteBranch(id)
}
