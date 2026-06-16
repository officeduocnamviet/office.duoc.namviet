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

func TestFundAccountsAPI(t *testing.T) {
	r := SetupTestRouter()
	
	// Create table sequence and reset if necessary, assuming auto increment
	supabase.DB.Exec("SELECT setval('fund_accounts_id_seq', (SELECT COALESCE(MAX(id), 1) FROM fund_accounts));")

	var createdFundAccountID float64
	uniqueCode := fmt.Sprintf("FUND_%d", time.Now().UnixNano())
	t.Run("Create Fund Account", func(t *testing.T) {
		body := `{"name":"Auto Test ` + uniqueCode + `", "type":"cash", "initial_balance": 1000}`
		w := PerformRequest(r, "POST", "/api/fund-accounts", body)
		
		if !assert.Equal(t, http.StatusCreated, w.Code, "Expected StatusCreated") {
			t.Logf("Response: %s", w.Body.String())
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if data, ok := response["data"].(map[string]interface{}); ok {
			if idFloat, ok := data["id"].(float64); ok {
				createdFundAccountID = idFloat
			}
		} else if idFloat, ok := response["id"].(float64); ok {
			createdFundAccountID = idFloat
		}
		assert.NotZero(t, createdFundAccountID)
	})

	t.Run("Get All Fund Accounts", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/fund-accounts", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Fund Account", func(t *testing.T) {
		if createdFundAccountID == 0 {
			t.Skip()
		}
		w := PerformRequest(r, "DELETE", fmt.Sprintf("/api/fund-accounts/%v", createdFundAccountID), "")
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestFinanceTransactionsAPI(t *testing.T) {
	r := SetupTestRouter()
	
	supabase.DB.Exec("SELECT setval('finance_transactions_id_seq', (SELECT COALESCE(MAX(id), 1) FROM finance_transactions));")

	// Create a fund account for the transaction
	var fundAccountID int64
	err := supabase.DB.Raw("INSERT INTO fund_accounts (name, type, status) VALUES ('TXN Fund', 'cash', 'active') RETURNING id;").Scan(&fundAccountID).Error
	if err != nil {
		t.Logf("Error creating fund account: %v", err)
	}

	var createdTxnID float64
	uniqueCode := fmt.Sprintf("TXN_%d", time.Now().UnixNano())
	
	t.Run("Create Finance Transaction", func(t *testing.T) {
		body := `{"code":"` + uniqueCode + `", "flow":"in", "amount": 500, "fund_account_id": ` + fmt.Sprintf("%d", fundAccountID) + `}`
		w := PerformRequest(r, "POST", "/api/finance-transactions", body)
		
		if !assert.Equal(t, http.StatusCreated, w.Code, "Expected StatusCreated") {
			t.Logf("Response: %s", w.Body.String())
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if data, ok := response["data"].(map[string]interface{}); ok {
			if idFloat, ok := data["id"].(float64); ok {
				createdTxnID = idFloat
			}
		} else if idFloat, ok := response["id"].(float64); ok {
			createdTxnID = idFloat
		}
		assert.NotZero(t, createdTxnID)
	})

	t.Run("Get All Finance Transactions", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/finance-transactions", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Finance Transaction", func(t *testing.T) {
		if createdTxnID == 0 {
			t.Skip()
		}
		w := PerformRequest(r, "DELETE", fmt.Sprintf("/api/finance-transactions/%v", createdTxnID), "")
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestAccountingJournalsAPI(t *testing.T) {
	r := SetupTestRouter()

	// Insert chart of accounts for the test to pass the foreign key constraint
	supabase.DB.Exec("INSERT INTO chart_of_accounts (account_code, name, type, balance_type) VALUES ('111', 'Tiền mặt', 'asset', 'BOTH') ON CONFLICT (account_code) DO NOTHING;")
	supabase.DB.Exec("INSERT INTO chart_of_accounts (account_code, name, type, balance_type) VALUES ('112', 'Tiền gửi ngân hàng', 'asset', 'BOTH') ON CONFLICT (account_code) DO NOTHING;")

	var createdJournalID string
	t.Run("Create Accounting Journal", func(t *testing.T) {
		body := `{"entry_date":"` + time.Now().Format(time.RFC3339) + `", "doc_type":"payment", "account_debit":"111", "account_credit":"112", "amount": 1000}`
		w := PerformRequest(r, "POST", "/api/accounting-journals", body)
		
		if !assert.Equal(t, http.StatusCreated, w.Code, "Expected StatusCreated") {
			t.Logf("Response: %s", w.Body.String())
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if data, ok := response["data"].(map[string]interface{}); ok {
			if idStr, ok := data["id"].(string); ok {
				createdJournalID = idStr
			}
		} else if idStr, ok := response["id"].(string); ok {
			createdJournalID = idStr
		}
		assert.NotEmpty(t, createdJournalID)
	})

	t.Run("Get All Accounting Journals", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/accounting-journals", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Accounting Journal", func(t *testing.T) {
		if createdJournalID == "" {
			t.Skip()
		}
		w := PerformRequest(r, "DELETE", "/api/accounting-journals/"+createdJournalID, "")
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
