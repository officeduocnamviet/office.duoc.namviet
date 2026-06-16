package manufacturers

import "time"

func GetAllManufacturersService() ([]Manufacturer, error) {
	return GetAllManufacturers()
}

func GetManufacturerByIDService(id int64) (*Manufacturer, error) {
	return GetManufacturerByID(id)
}

func CreateManufacturerService(req CreateManufacturerRequest) (*Manufacturer, error) {
	status := "active"
	if req.Status != "" {
		status = req.Status
	}

	manufacturer := &Manufacturer{
		Name:    req.Name,
		Country: req.Country,
		Status:  status,
	}

	if err := CreateManufacturer(manufacturer); err != nil {
		return nil, err
	}
	return manufacturer, nil
}

func UpdateManufacturerService(id int64, req UpdateManufacturerRequest) (*Manufacturer, error) {
	manufacturer, err := GetManufacturerByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		manufacturer.Name = *req.Name
	}
	if req.Country != nil {
		manufacturer.Country = req.Country
	}
	if req.Status != nil {
		manufacturer.Status = *req.Status
	}
	
	now := time.Now()
	manufacturer.UpdatedAt = &now

	if err := UpdateManufacturer(manufacturer); err != nil {
		return nil, err
	}
	return manufacturer, nil
}
