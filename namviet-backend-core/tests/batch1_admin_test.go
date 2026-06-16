package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Company struct {
	ID                  string     `json:"id"`
	TaxCode             string     `json:"tax_code"`
	Name                string     `json:"name"`
	Address             string     `json:"address"`
	Phone               string     `json:"phone"`
	Email               string     `json:"email"`
	RepresentativeName  string     `json:"representative_name"`
	Status              string     `json:"status"`
	CreatedAt           *time.Time `json:"created_at,omitempty"`
}

func TestCompanyAPI(t *testing.T) {
	r := SetupTestRouter()

	var createdCompanyID string

	t.Run("Create Company", func(t *testing.T) {
		body := `{"tax_code":"1234567890", "name":"Công ty Cổ phần Mẫu AI", "address":"Hà Nội", "phone":"0987654321", "email":"ai@test.com", "representative_name":"Nguyễn Văn AI"}`
		w := PerformRequest(r, "POST", "/api/companies", body)
		
		assert.Equal(t, http.StatusCreated, w.Code, "Expected StatusCreated")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data, ok := response["data"].(map[string]interface{})
		if ok {
			createdCompanyID = data["id"].(string)
		} else {
			// If it doesn't wrap in "data"
			if val, exists := response["id"]; exists {
				createdCompanyID = val.(string)
			}
		}
		assert.NotEmpty(t, createdCompanyID, "Should have created a company with ID")
	})

	t.Run("Get All Companies", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/companies", "")
		assert.Equal(t, http.StatusOK, w.Code, "Expected StatusOK")
	})

	t.Run("Get Single Company", func(t *testing.T) {
		if createdCompanyID == "" {
			t.Skip("Skipping because create failed")
		}
		w := PerformRequest(r, "GET", "/api/companies/"+createdCompanyID, "")
		assert.Equal(t, http.StatusOK, w.Code, "Expected StatusOK")
	})

	t.Run("Update Company", func(t *testing.T) {
		if createdCompanyID == "" {
			t.Skip("Skipping because create failed")
		}
		body := `{"name":"Công ty Cổ phần Mẫu AI - Đã Sửa"}`
		w := PerformRequest(r, "PUT", "/api/companies/"+createdCompanyID, body)
		assert.Equal(t, http.StatusOK, w.Code, "Expected StatusOK")
	})

	t.Run("Delete Company", func(t *testing.T) {
		if createdCompanyID == "" {
			t.Skip("Skipping because create failed")
		}
		w := PerformRequest(r, "DELETE", "/api/companies/"+createdCompanyID, "")
		assert.Equal(t, http.StatusNoContent, w.Code, "Expected StatusNoContent")
	})
}
