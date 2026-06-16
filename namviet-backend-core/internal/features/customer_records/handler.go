package customer_records

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- Vaccination Records API ---

// GetAllVaccinationRecordsHandler
// @Summary Get All Vaccination Records
// @Description Retrieve a list of all vaccination records
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} CustomerVaccinationRecord
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vaccination-records [get]
func GetAllVaccinationRecordsHandler(c *gin.Context) {
	items, err := GetAllVaccinationRecordsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetVaccinationRecordHandler
// @Summary Get Vaccination Record by ID
// @Description Retrieve a specific vaccination record
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Record ID"
// @Success 200 {object} CustomerVaccinationRecord
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /vaccination-records/{id} [get]
func GetVaccinationRecordHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetVaccinationRecordByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateVaccinationRecordHandler
// @Summary Create a new Vaccination Record
// @Description Create a new vaccination record
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateVaccinationRecordRequest true "Record Details"
// @Success 201 {object} CustomerVaccinationRecord
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vaccination-records [post]
func CreateVaccinationRecordHandler(c *gin.Context) {
	var req CreateVaccinationRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateVaccinationRecordService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateVaccinationRecordHandler
// @Summary Update Vaccination Record
// @Description Update vaccination record details
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Record ID"
// @Param request body UpdateVaccinationRecordRequest true "Update Details"
// @Success 200 {object} CustomerVaccinationRecord
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vaccination-records/{id} [put]
func UpdateVaccinationRecordHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateVaccinationRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateVaccinationRecordService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteVaccinationRecordHandler
// @Summary Delete Vaccination Record
// @Description Delete a vaccination record
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Record ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vaccination-records/{id} [delete]
func DeleteVaccinationRecordHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteVaccinationRecordService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// --- Customer Vouchers API ---

// GetAllCustomerVouchersHandler
// @Summary Get All Customer Vouchers
// @Description Retrieve a list of all customer vouchers
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} CustomerVoucher
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customer-vouchers [get]
func GetAllCustomerVouchersHandler(c *gin.Context) {
	items, err := GetAllCustomerVouchersService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetCustomerVoucherHandler
// @Summary Get Customer Voucher by ID
// @Description Retrieve a specific customer voucher
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Voucher ID"
// @Success 200 {object} CustomerVoucher
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /customer-vouchers/{id} [get]
func GetCustomerVoucherHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetCustomerVoucherByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateCustomerVoucherHandler
// @Summary Create a new Customer Voucher
// @Description Create a new customer voucher
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCustomerVoucherRequest true "Voucher Details"
// @Success 201 {object} CustomerVoucher
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customer-vouchers [post]
func CreateCustomerVoucherHandler(c *gin.Context) {
	var req CreateCustomerVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateCustomerVoucherService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateCustomerVoucherHandler
// @Summary Update Customer Voucher
// @Description Update customer voucher details
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Voucher ID"
// @Param request body UpdateCustomerVoucherRequest true "Update Details"
// @Success 200 {object} CustomerVoucher
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customer-vouchers/{id} [put]
func UpdateCustomerVoucherHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateCustomerVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateCustomerVoucherService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteCustomerVoucherHandler
// @Summary Delete Customer Voucher
// @Description Delete a customer voucher
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Voucher ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customer-vouchers/{id} [delete]
func DeleteCustomerVoucherHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteCustomerVoucherService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
