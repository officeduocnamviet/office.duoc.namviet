package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"github.com/stretchr/testify/assert"
)

func TestWarehousesAPI(t *testing.T) {
	r := SetupTestRouter()

	supabase.DB.Exec("SELECT setval('warehouses_id_seq', (SELECT COALESCE(MAX(id), 1) FROM warehouses));")

	var companyID string
	err := supabase.DB.Table("companies").Select("id").Limit(1).Pluck("id", &companyID).Error
	if err != nil || companyID == "" {
		t.Skip("Skipping warehouses test because no companies exist")
	}

	var createdWarehouseID string
	uniqueKey := "WH_" + t.Name()
	t.Run("Create Warehouse", func(t *testing.T) {
		body := `{"name":"Auto Test Warehouse", "key":"` + uniqueKey + `", "address":"Hanoi"}`
		w := PerformRequest(r, "POST", "/api/warehouses", body)
		
		if !assert.Equal(t, http.StatusCreated, w.Code, "Expected StatusCreated") {
			t.Logf("Response: %s", w.Body.String())
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if data, ok := response["data"].(map[string]interface{}); ok {
			if floatId, ok := data["id"].(float64); ok {
				createdWarehouseID = fmt.Sprintf("%.0f", floatId)
			}
		} else if val, exists := response["id"]; exists {
			if floatId, ok := val.(float64); ok {
				createdWarehouseID = fmt.Sprintf("%.0f", floatId)
			}
		}
		assert.NotEmpty(t, createdWarehouseID)
	})

	t.Run("Get All Warehouses", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/warehouses", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Warehouse", func(t *testing.T) {
		if createdWarehouseID == "" {
			t.Skip()
		}
		w := PerformRequest(r, "DELETE", "/api/warehouses/"+createdWarehouseID, "")
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestInventoryAPI(t *testing.T) {
	r := SetupTestRouter()

	var companyID string
	supabase.DB.Table("companies").Select("id").Limit(1).Pluck("id", &companyID)
	
	var categoryID string
	supabase.DB.Table("categories").Select("id").Limit(1).Pluck("id", &categoryID)

	var warehouseID, productID string
	
	// Create test warehouse directly
	uniqueWHKey := "WH_INV_" + t.Name()
	supabase.DB.Raw("INSERT INTO warehouses (key, name, unit, type) VALUES (?, ?, 'Hộp', 'retail') RETURNING id;", uniqueWHKey, "Inv Test WH").Scan(&warehouseID)
	
	// Create test product directly
	supabase.DB.Raw("INSERT INTO products (name, category_id, status) VALUES (?, ?, 'active') RETURNING id;", "Inv Test Product", categoryID).Scan(&productID)

	if warehouseID == "" || productID == "" {
		t.Skip("Skipping inventory test because pre-requisites could not be created")
	}

	t.Run("Check Columns", func(t *testing.T) {
		var def string
		supabase.DB.Raw("SELECT pg_get_constraintdef(oid) FROM pg_constraint WHERE conname = 'inventory_transactions_type_check'").Scan(&def)
		t.Logf("CONSTRAINT DEF: %s", def)
	})

	t.Run("Create Inventory Transaction", func(t *testing.T) {
		body := `{"warehouse_id":` + warehouseID + `, "product_id":` + productID + `, "type":"inbound", "quantity":100}`
		w := PerformRequest(r, "POST", "/api/inventory/transactions", body)
		
		if !assert.Equal(t, http.StatusCreated, w.Code, "Expected StatusCreated") {
			t.Logf("Response: %s", w.Body.String())
		}
	})

	t.Run("Get Inventory", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/inventory", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
