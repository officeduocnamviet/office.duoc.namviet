package customers

type CreateCustomerRequest struct {
	CustomerCode *string     `json:"customer_code"`
	Name         string      `json:"name" binding:"required"`
	CustomerType string      `json:"customer_type"` // B2B or B2C
	Phone        *string     `json:"phone"`
	Email        *string     `json:"email"`
	Address      *string     `json:"address"`
	DOB          *string     `json:"dob"`
	Gender       *string     `json:"gender"`
	CCCD         *string     `json:"cccd"`
	B2BMetadata  JSONMap     `json:"b2b_metadata"`
}

type UpdateCustomerRequest struct {
	CustomerCode *string     `json:"customer_code"`
	Name         *string     `json:"name"`
	CustomerType *string     `json:"customer_type"`
	Phone        *string     `json:"phone"`
	Email        *string     `json:"email"`
	Address      *string     `json:"address"`
	Status       *string     `json:"status"`
	DOB          *string     `json:"dob"`
	Gender       *string     `json:"gender"`
	CCCD         *string     `json:"cccd"`
	B2BMetadata  JSONMap     `json:"b2b_metadata"`
}
