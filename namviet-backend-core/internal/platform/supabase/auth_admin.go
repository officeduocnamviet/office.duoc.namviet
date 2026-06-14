package supabase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// SupabaseAuthUser represents the response from Supabase Auth Admin API
type SupabaseAuthUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// CreateAuthUser calls Supabase Management API to create a real user in auth.users
func CreateAuthUser(email, password, fullName string) (*SupabaseAuthUser, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	serviceRoleKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if supabaseURL == "" || serviceRoleKey == "" {
		return nil, errors.New("missing SUPABASE_URL or SUPABASE_SERVICE_ROLE_KEY")
	}

	endpoint := fmt.Sprintf("%s/auth/v1/admin/users", supabaseURL)

	payload := map[string]interface{}{
		"email":          email,
		"password":       password,
		"email_confirm":  true, // Auto-confirm for employee
		"user_metadata": map[string]string{
			"full_name": fullName,
		},
	}

	bodyBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", serviceRoleKey)
	req.Header.Set("Authorization", "Bearer "+serviceRoleKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create user in Supabase Auth, status code: %d", resp.StatusCode)
	}

	var authUser SupabaseAuthUser
	if err := json.NewDecoder(resp.Body).Decode(&authUser); err != nil {
		return nil, err
	}

	return &authUser, nil
}

// DeleteAuthUser deletes a user from auth.users
func DeleteAuthUser(userID string) error {
	supabaseURL := os.Getenv("SUPABASE_URL")
	serviceRoleKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	endpoint := fmt.Sprintf("%s/auth/v1/admin/users/%s", supabaseURL, userID)

	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("apikey", serviceRoleKey)
	req.Header.Set("Authorization", "Bearer "+serviceRoleKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete user in Supabase Auth, status code: %d", resp.StatusCode)
	}

	return nil
}
