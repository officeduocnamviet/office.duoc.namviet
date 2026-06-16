package customer_records

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the customer records routes
func RegisterRoutes(router *gin.RouterGroup) {
	vaccineGroup := router.Group("/vaccination-records")
	{
		vaccineGroup.GET("", GetAllVaccinationRecordsHandler)
		vaccineGroup.GET("/:id", GetVaccinationRecordHandler)
		vaccineGroup.POST("", CreateVaccinationRecordHandler)
		vaccineGroup.PUT("/:id", UpdateVaccinationRecordHandler)
		vaccineGroup.DELETE("/:id", DeleteVaccinationRecordHandler)
	}

	voucherGroup := router.Group("/customer-vouchers")
	{
		voucherGroup.GET("", GetAllCustomerVouchersHandler)
		voucherGroup.GET("/:id", GetCustomerVoucherHandler)
		voucherGroup.POST("", CreateCustomerVoucherHandler)
		voucherGroup.PUT("/:id", UpdateCustomerVoucherHandler)
		voucherGroup.DELETE("/:id", DeleteCustomerVoucherHandler)
	}
}
