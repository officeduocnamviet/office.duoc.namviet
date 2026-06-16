package shipping_partners

import "time"

func GetAllShippingPartnersService() ([]ShippingPartner, error) {
	return GetAllShippingPartners()
}

func GetShippingPartnerByIDService(id string) (*ShippingPartner, error) {
	return GetShippingPartnerByID(id)
}

func CreateShippingPartnerService(req CreateShippingPartnerRequest) (*ShippingPartner, error) {
	partner := &ShippingPartner{
		Code:                req.Code,
		Name:                req.Name,
		PartnerType:         req.PartnerType,
		TrackingURLTemplate: req.TrackingURLTemplate,
		Status:              "active",
	}

	if req.APIConfig != nil {
		partner.APIConfig = *req.APIConfig
	}

	if err := CreateShippingPartner(partner); err != nil {
		return nil, err
	}
	return partner, nil
}

func UpdateShippingPartnerService(id string, req UpdateShippingPartnerRequest) (*ShippingPartner, error) {
	partner, err := GetShippingPartnerByID(id)
	if err != nil {
		return nil, err
	}

	if req.Code != nil {
		partner.Code = *req.Code
	}
	if req.Name != nil {
		partner.Name = *req.Name
	}
	if req.PartnerType != nil {
		partner.PartnerType = *req.PartnerType
	}
	if req.APIConfig != nil {
		partner.APIConfig = *req.APIConfig
	}
	if req.TrackingURLTemplate != nil {
		partner.TrackingURLTemplate = req.TrackingURLTemplate
	}
	if req.Status != nil {
		partner.Status = *req.Status
	}

	now := time.Now()
	partner.UpdatedAt = &now

	if err := UpdateShippingPartner(partner); err != nil {
		return nil, err
	}
	return partner, nil
}

func DeleteShippingPartnerService(id string) error {
	return DeleteShippingPartner(id)
}
