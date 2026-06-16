package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleAPI(t *testing.T) {
	r := SetupTestRouter()

	var createdRoleID string

	t.Run("Create Role", func(t *testing.T) {
		body := `{"name":"Auto Test Role", "description":"Role for automated testing", "permissions":["users.read","users.write"]}`
		w := PerformRequest(r, "POST", "/api/roles", body)
		
		assert.Equal(t, http.StatusCreated, w.Code, "Expected StatusCreated")
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data, ok := response["data"].(map[string]interface{})
		if ok {
			createdRoleID = data["id"].(string)
		} else {
			if val, exists := response["id"]; exists {
				createdRoleID = val.(string)
			}
		}
		assert.NotEmpty(t, createdRoleID, "Should have created a role with ID")
	})

	t.Run("Get All Roles", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/roles", "")
		assert.Equal(t, http.StatusOK, w.Code, "Expected StatusOK")
	})

	t.Run("Get Single Role", func(t *testing.T) {
		if createdRoleID == "" {
			t.Skip("Skipping because create failed")
		}
		w := PerformRequest(r, "GET", "/api/roles/"+createdRoleID, "")
		assert.Equal(t, http.StatusOK, w.Code, "Expected StatusOK")
	})

	t.Run("Update Role", func(t *testing.T) {
		if createdRoleID == "" {
			t.Skip("Skipping because create failed")
		}
		body := `{"name":"Auto Test Role - Edited", "permissions":["users.read"]}`
		w := PerformRequest(r, "PUT", "/api/roles/"+createdRoleID, body)
		assert.Equal(t, http.StatusOK, w.Code, "Expected StatusOK")
	})

	t.Run("Delete Role", func(t *testing.T) {
		if createdRoleID == "" {
			t.Skip("Skipping because create failed")
		}
		w := PerformRequest(r, "DELETE", "/api/roles/"+createdRoleID, "")
		assert.Equal(t, http.StatusNoContent, w.Code, "Expected StatusNoContent")
	})
}
