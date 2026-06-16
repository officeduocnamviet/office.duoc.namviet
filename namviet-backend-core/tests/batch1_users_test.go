package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"github.com/stretchr/testify/assert"
)

func TestUserAPI(t *testing.T) {
	r := SetupTestRouter()

	var createdUserID string
	// Generate a unique email for each test run to avoid unique constraint violations
	uniqueEmail := "autotest_" + t.Name() + "@example.com"

	var roleID, companyID string
	// Fetch a random role and company to satisfy foreign key constraints
	err := supabase.DB.Table("roles").Select("id").Limit(1).Pluck("id", &roleID).Error
	if err != nil || roleID == "" {
		t.Skip("Skipping user test because no roles exist in DB")
	}
	err = supabase.DB.Table("companies").Select("id").Limit(1).Pluck("id", &companyID).Error
	if err != nil || companyID == "" {
		t.Skip("Skipping user test because no companies exist in DB")
	}

	t.Run("Create User", func(t *testing.T) {
		body := `{"email":"` + uniqueEmail + `", "password":"password123", "full_name":"Auto Test User", "phone":"0999999999", "role_id":"` + roleID + `", "company_id":"` + companyID + `"}`
		w := PerformRequest(r, "POST", "/api/users", body)
		if !assert.Equal(t, http.StatusCreated, w.Code, "Expected StatusCreated") {
			t.Logf("Response: %s", w.Body.String())
		}
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data, ok := response["data"].(map[string]interface{})
		if ok {
			createdUserID = data["id"].(string)
		} else {
			if val, exists := response["id"]; exists {
				createdUserID = val.(string)
			}
		}
		assert.NotEmpty(t, createdUserID, "Should have created a user with ID")
	})

	t.Run("Get All Users", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/users", "")
		assert.Equal(t, http.StatusOK, w.Code, "Expected StatusOK")
	})

	t.Run("Update User", func(t *testing.T) {
		if createdUserID == "" {
			t.Skip("Skipping because create failed")
		}
		body := `{"full_name":"Auto Test User - Edited"}`
		w := PerformRequest(r, "PUT", "/api/users/"+createdUserID, body)
		assert.Equal(t, http.StatusOK, w.Code, "Expected StatusOK")
	})

	t.Run("Delete User", func(t *testing.T) {
		if createdUserID == "" {
			t.Skip("Skipping because create failed")
		}
		w := PerformRequest(r, "DELETE", "/api/users/"+createdUserID, "")
		assert.Equal(t, http.StatusNoContent, w.Code, "Expected StatusNoContent")
	})
}
