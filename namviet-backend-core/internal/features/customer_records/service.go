package customer_records

// Vaccination Records
func GetAllVaccinationRecordsService() ([]CustomerVaccinationRecord, error) {
	return GetAllVaccinationRecords()
}

func GetVaccinationRecordByIDService(id string) (*CustomerVaccinationRecord, error) {
	return GetVaccinationRecordByID(id)
}

func CreateVaccinationRecordService(req CreateVaccinationRecordRequest) (*CustomerVaccinationRecord, error) {
	dose := 1
	if req.DoseNumber != 0 {
		dose = req.DoseNumber
	}

	record := &CustomerVaccinationRecord{
		CustomerID:      req.CustomerID,
		VaccineName:     req.VaccineName,
		DoseNumber:      dose,
		VaccinationDate: req.VaccinationDate,
		NextDueDate:     req.NextDueDate,
		AdministeredBy:  req.AdministeredBy,
		Notes:           req.Notes,
	}

	if err := CreateVaccinationRecord(record); err != nil {
		return nil, err
	}
	return record, nil
}

func UpdateVaccinationRecordService(id string, req UpdateVaccinationRecordRequest) (*CustomerVaccinationRecord, error) {
	record, err := GetVaccinationRecordByID(id)
	if err != nil {
		return nil, err
	}

	if req.VaccineName != nil {
		record.VaccineName = *req.VaccineName
	}
	if req.DoseNumber != nil {
		record.DoseNumber = *req.DoseNumber
	}
	if req.VaccinationDate != nil {
		record.VaccinationDate = *req.VaccinationDate
	}
	if req.NextDueDate != nil {
		record.NextDueDate = req.NextDueDate
	}
	if req.AdministeredBy != nil {
		record.AdministeredBy = req.AdministeredBy
	}
	if req.Notes != nil {
		record.Notes = req.Notes
	}

	if err := UpdateVaccinationRecord(record); err != nil {
		return nil, err
	}
	return record, nil
}

func DeleteVaccinationRecordService(id string) error {
	return DeleteVaccinationRecord(id)
}

// Vouchers
func GetAllCustomerVouchersService() ([]CustomerVoucher, error) {
	return GetAllCustomerVouchers()
}

func GetCustomerVoucherByIDService(id string) (*CustomerVoucher, error) {
	return GetCustomerVoucherByID(id)
}

func CreateCustomerVoucherService(req CreateCustomerVoucherRequest) (*CustomerVoucher, error) {
	voucher := &CustomerVoucher{
		CustomerID:  req.CustomerID,
		PromotionID: req.PromotionID,
		VoucherCode: req.VoucherCode,
		IsUsed:      false,
	}

	if err := CreateCustomerVoucher(voucher); err != nil {
		return nil, err
	}
	return voucher, nil
}

func UpdateCustomerVoucherService(id string, req UpdateCustomerVoucherRequest) (*CustomerVoucher, error) {
	voucher, err := GetCustomerVoucherByID(id)
	if err != nil {
		return nil, err
	}

	if req.IsUsed != nil {
		voucher.IsUsed = *req.IsUsed
	}
	if req.UsedAt != nil {
		voucher.UsedAt = req.UsedAt
	}

	if err := UpdateCustomerVoucher(voucher); err != nil {
		return nil, err
	}
	return voucher, nil
}

func DeleteCustomerVoucherService(id string) error {
	return DeleteCustomerVoucher(id)
}
