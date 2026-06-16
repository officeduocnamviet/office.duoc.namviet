package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"github.com/stretchr/testify/assert"
)

func TestCategoriesAPI(t *testing.T) {
	r := SetupTestRouter()

	// Fix sequence issue if testing on existing DB
	supabase.DB.Exec("SELECT setval('categories_id_seq', (SELECT COALESCE(MAX(id), 1) FROM categories));")

	var createdCategoryID string
	t.Run("Create Category", func(t *testing.T) {
		body := `{"name":"Auto Test Category", "slug":"auto-test-category", "description":"Test description"}`
		w := PerformRequest(r, "POST", "/api/categories", body)
		
		if !assert.Equal(t, http.StatusCreated, w.Code, "Expected StatusCreated") {
			t.Logf("Response: %s", w.Body.String())
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if data, ok := response["data"].(map[string]interface{}); ok {
			if floatId, ok := data["id"].(float64); ok {
				createdCategoryID = fmt.Sprintf("%.0f", floatId)
			}
		} else if val, exists := response["id"]; exists {
			if floatId, ok := val.(float64); ok {
				createdCategoryID = fmt.Sprintf("%.0f", floatId)
			}
		}
		assert.NotEmpty(t, createdCategoryID)
	})

	t.Run("Get All Categories", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/categories", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update Category", func(t *testing.T) {
		if createdCategoryID == "" {
			t.Skip()
		}
		body := `{"name":"Auto Test Category Updated"}`
		w := PerformRequest(r, "PUT", "/api/categories/"+createdCategoryID, body)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Category", func(t *testing.T) {
		if createdCategoryID == "" {
			t.Skip()
		}
		w := PerformRequest(r, "DELETE", "/api/categories/"+createdCategoryID, "")
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestProductsAPI(t *testing.T) {
	r := SetupTestRouter()

	var categoryID string
	err := supabase.DB.Table("categories").Select("id").Limit(1).Pluck("id", &categoryID).Error
	if err != nil || categoryID == "" {
		t.Skip("Skipping products test because no categories exist")
	}

	// Fix sequence issue if testing on existing DB
	supabase.DB.Exec("SELECT setval('products_id_seq', (SELECT COALESCE(MAX(id), 1) FROM products));")

	var createdProductID string
	uniqueBarcode := "BC_" + t.Name()

	t.Run("Create Product", func(t *testing.T) {
		body := `{"code":"` + uniqueBarcode + `", "name":"Auto Test Product", "base_price":150000, "category_id":` + categoryID + `}`
		w := PerformRequest(r, "POST", "/api/products", body)
		
		if !assert.Equal(t, http.StatusCreated, w.Code, "Expected StatusCreated") {
			t.Logf("Response: %s", w.Body.String())
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if data, ok := response["data"].(map[string]interface{}); ok {
			if floatId, ok := data["id"].(float64); ok {
				createdProductID = fmt.Sprintf("%.0f", floatId)
			}
		} else if val, exists := response["id"]; exists {
			if floatId, ok := val.(float64); ok {
				createdProductID = fmt.Sprintf("%.0f", floatId)
			}
		}
		assert.NotEmpty(t, createdProductID)
	})

	t.Run("Get All Products", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/products", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Product", func(t *testing.T) {
		if createdProductID == "" {
			t.Skip()
		}
		w := PerformRequest(r, "DELETE", "/api/products/"+createdProductID, "")
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
