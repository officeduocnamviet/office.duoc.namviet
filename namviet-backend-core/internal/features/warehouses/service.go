package warehouses

func GetAllWarehousesService() ([]Warehouse, error) {
	return GetAllWarehouses()
}

func GetWarehouseByIDService(id int64) (*Warehouse, error) {
	return GetWarehouseByID(id)
}

func CreateWarehouseService(req CreateWarehouseRequest) (*Warehouse, error) {
	whType := "retail"
	if req.Type != "" {
		whType = req.Type
	}

	unit := "Hộp"
	if req.Unit != "" {
		unit = req.Unit
	}

	warehouse := &Warehouse{
		Key:       req.Key,
		Name:      req.Name,
		Unit:      unit,
		Address:   req.Address,
		Type:      whType,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Code:      req.Code,
		Manager:   req.Manager,
		Phone:     req.Phone,
	}

	if err := CreateWarehouse(warehouse); err != nil {
		return nil, err
	}
	return warehouse, nil
}

func UpdateWarehouseService(id int64, req UpdateWarehouseRequest) (*Warehouse, error) {
	warehouse, err := GetWarehouseByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		warehouse.Name = *req.Name
	}
	if req.Unit != nil {
		warehouse.Unit = *req.Unit
	}
	if req.Type != nil {
		warehouse.Type = *req.Type
	}
	if req.Address != nil {
		warehouse.Address = req.Address
	}
	if req.Manager != nil {
		warehouse.Manager = req.Manager
	}
	if req.Code != nil {
		warehouse.Code = req.Code
	}
	if req.Phone != nil {
		warehouse.Phone = req.Phone
	}
	if req.Latitude != nil {
		warehouse.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		warehouse.Longitude = req.Longitude
	}

	if err := UpdateWarehouse(warehouse); err != nil {
		return nil, err
	}
	return warehouse, nil
}

func DeleteWarehouseService(id int64) error {
	return DeleteWarehouse(id)
}
