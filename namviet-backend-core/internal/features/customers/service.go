package customers

import "time"

func GetAllCustomersService() ([]Customer, error) {
	return GetAllCustomers()
}

func GetCustomerByIDService(id int64) (*Customer, error) {
	return GetCustomerByID(id)
}

func CreateCustomerService(req CreateCustomerRequest) (*Customer, error) {
	customerType := "B2C"
	if req.CustomerType != "" {
		customerType = req.CustomerType
	}

	customer := &Customer{
		CustomerCode: req.CustomerCode,
		Name:         req.Name,
		CustomerType: customerType,
		Phone:        req.Phone,
		Email:        req.Email,
		Address:      req.Address,
		DOB:          req.DOB,
		Gender:       req.Gender,
		CCCD:         req.CCCD,
		B2BMetadata:  req.B2BMetadata,
		Status:       "active",
	}

	if err := CreateCustomer(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func UpdateCustomerService(id int64, req UpdateCustomerRequest) (*Customer, error) {
	customer, err := GetCustomerByID(id)
	if err != nil {
		return nil, err
	}

	if req.CustomerCode != nil {
		customer.CustomerCode = req.CustomerCode
	}
	if req.Name != nil {
		customer.Name = *req.Name
	}
	if req.CustomerType != nil {
		customer.CustomerType = *req.CustomerType
	}
	if req.Phone != nil {
		customer.Phone = req.Phone
	}
	if req.Email != nil {
		customer.Email = req.Email
	}
	if req.Address != nil {
		customer.Address = req.Address
	}
	if req.Status != nil {
		customer.Status = *req.Status
	}
	if req.DOB != nil {
		customer.DOB = req.DOB
	}
	if req.Gender != nil {
		customer.Gender = req.Gender
	}
	if req.CCCD != nil {
		customer.CCCD = req.CCCD
	}
	if req.B2BMetadata != nil {
		customer.B2BMetadata = req.B2BMetadata
	}
	
	now := time.Now()
	customer.UpdatedAt = &now

	if err := UpdateCustomer(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func DeleteCustomerService(id int64) error {
	return DeleteCustomer(id)
}
