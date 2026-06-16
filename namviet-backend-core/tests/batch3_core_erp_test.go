package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"github.com/stretchr/testify/assert"
)

func TestCustomersAPI(t *testing.T) {
	r := SetupTestRouter()

	supabase.DB.Exec("SELECT setval('customers_id_seq', (SELECT COALESCE(MAX(id), 1) FROM customers));")

	var createdCustomerID string
	uniqueCode := fmt.Sprintf("CUST_%d", time.Now().UnixNano())
	t.Run("Create Customer", func(t *testing.T) {
		body := `{"name":"Auto Test Customer", "customer_code":"` + uniqueCode + `"}`
		w := PerformRequest(r, "POST", "/api/customers", body)
		
		if !assert.Equal(t, http.StatusCreated, w.Code, "Expected StatusCreated") {
			t.Logf("Response: %s", w.Body.String())
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if data, ok := response["data"].(map[string]interface{}); ok {
			if floatId, ok := data["id"].(float64); ok {
				createdCustomerID = fmt.Sprintf("%.0f", floatId)
			}
		} else if val, exists := response["id"]; exists {
			if floatId, ok := val.(float64); ok {
				createdCustomerID = fmt.Sprintf("%.0f", floatId)
			}
		}
		assert.NotEmpty(t, createdCustomerID)
	})

	t.Run("Get All Customers", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/customers", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Customer", func(t *testing.T) {
		if createdCustomerID == "" {
			t.Skip()
		}
		w := PerformRequest(r, "DELETE", "/api/customers/"+createdCustomerID, "")
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}



func TestOrdersAPI(t *testing.T) {
	r := SetupTestRouter()

	var createdOrderID string
	uniqueCode := fmt.Sprintf("ORD_%d", time.Now().UnixNano())
	
	// Create a category and product for the order items
	var productID int64
	supabase.DB.Raw("INSERT INTO products (name, status) VALUES ('Order Test Product', 'active') RETURNING id;").Scan(&productID)
	
	t.Run("Create Order", func(t *testing.T) {
		body := `{"code":"` + uniqueCode + `", "status":"PENDING", "order_type":"B2C", "items": [{"product_id": ` + fmt.Sprintf("%d", productID) + `, "quantity": 1, "unit_price": 100, "uom": "cai"}]}`
		w := PerformRequest(r, "POST", "/api/orders", body)
		
		if !assert.Equal(t, http.StatusCreated, w.Code, "Expected StatusCreated") {
			t.Logf("Response: %s", w.Body.String())
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if data, ok := response["data"].(map[string]interface{}); ok {
			if idStr, ok := data["id"].(string); ok {
				createdOrderID = idStr
			}
		} else if idStr, ok := response["id"].(string); ok {
			createdOrderID = idStr
		}
		assert.NotEmpty(t, createdOrderID)
	})

	t.Run("Get All Orders", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/orders", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Order", func(t *testing.T) {
		if createdOrderID == "" {
			t.Skip()
		}
		w := PerformRequest(r, "DELETE", "/api/orders/"+createdOrderID, "")
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
