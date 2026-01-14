package tests

import (
	"testing"
)

// TestPaymentCreate_ValidMethod tests payment creation with valid method
func TestPaymentCreate_ValidMethod(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)

	query := `
		mutation {
			paymentCreate(input: {
				clientId: $clientId
				amount: 100.0
				method: "mobile-money"
				description: "Test payment"
			}) {
				id
				clientId
				amount
				method
				status
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["paymentCreate"].(map[string]interface{})
	if data["id"] == nil {
		t.Error("Payment ID should not be nil")
	}
	if data["status"].(string) != "completed" {
		t.Error("Payment status should be 'completed'")
	}
}

// TestPaymentCreate_CaisseExit tests caisse exit for payment
func TestPaymentCreate_CaisseExit(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// First add some money to caisse
	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")
	CreateTestSale(t, tc, clientID, productID, 200.0, "paid")

	// Get initial caisse balance
	query := `
		query {
			caisse {
				balance
				totalSorties
			}
		}
	`
	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)
	caisseData := resp.Data["caisse"].(map[string]interface{})
	initialBalance := caisseData["balance"].(float64)

	// Create payment
	paymentQuery := `
		mutation {
			paymentCreate(input: {
				clientId: $clientId
				amount: 50.0
				method: "cash"
				description: "Client payment"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}
	ExecuteGraphQL(t, tc, paymentQuery, variables, tc.AdminToken)

	// Check caisse balance decreased
	resp2 := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp2)
	caisseData2 := resp2.Data["caisse"].(map[string]interface{})
	newBalance := caisseData2["balance"].(float64)

	if newBalance >= initialBalance {
		t.Error("Caisse balance should decrease after payment")
	}
}

// TestPaymentCreate_InvalidMethod tests payment creation with invalid method
func TestPaymentCreate_InvalidMethod(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)

	query := `
		mutation {
			paymentCreate(input: {
				clientId: $clientId
				amount: 100.0
				method: "invalid-method"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertHasErrors(t, resp)
}

// TestPayments_ListWithFilters tests listing payments with filters
func TestPayments_ListWithFilters(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)

	// Create payment
	query := `
		mutation {
			paymentCreate(input: {
				clientId: $clientId
				amount: 100.0
				method: "mobile-money"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}
	ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)

	// List payments
	listQuery := `
		query {
			payments(filter: {
				status: "completed"
			}) {
				id
				amount
				method
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, listQuery, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	payments := resp.Data["payments"].([]interface{})
	if len(payments) == 0 {
		t.Error("Should find at least one payment")
	}
}

// TestPayment_GetByID tests getting a payment by ID
func TestPayment_GetByID(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)

	// Create payment
	query := `
		mutation {
			paymentCreate(input: {
				clientId: $clientId
				amount: 100.0
				method: "mobile-money"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}
	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)
	paymentData := resp.Data["paymentCreate"].(map[string]interface{})
	paymentID := paymentData["id"].(string)

	// Get payment
	getQuery := `
		query {
			payment(id: $paymentId) {
				id
				amount
				method
				client {
					id
					name
				}
			}
		}
	`
	getVars := map[string]interface{}{
		"paymentId": paymentID,
	}

	getResp := ExecuteGraphQL(t, tc, getQuery, getVars, tc.AdminToken)
	AssertNoErrors(t, getResp)

	payment := getResp.Data["payment"].(map[string]interface{})
	if payment["id"].(string) != paymentID {
		t.Error("Payment ID should match")
	}
}

// TestPaymentUpdate tests updating a payment
func TestPaymentUpdate(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)

	// Create payment
	query := `
		mutation {
			paymentCreate(input: {
				clientId: $clientId
				amount: 100.0
				method: "mobile-money"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}
	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)
	paymentData := resp.Data["paymentCreate"].(map[string]interface{})
	paymentID := paymentData["id"].(string)

	// Update payment
	updateQuery := `
		mutation {
			paymentUpdate(id: $paymentId, input: {
				clientId: $clientId
				amount: 150.0
				method: "cash"
				description: "Updated payment"
			}) {
				id
				amount
				method
			}
		}
	`
	updateVars := map[string]interface{}{
		"paymentId": paymentID,
		"clientId":  clientID,
	}

	updateResp := ExecuteGraphQL(t, tc, updateQuery, updateVars, tc.AdminToken)
	AssertNoErrors(t, updateResp)

	data := updateResp.Data["paymentUpdate"].(map[string]interface{})
	if data["amount"].(float64) != 150.0 {
		t.Error("Payment amount should be updated")
	}
}

// TestPaymentDelete tests deleting a payment
func TestPaymentDelete(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)

	// Create payment
	query := `
		mutation {
			paymentCreate(input: {
				clientId: $clientId
				amount: 100.0
				method: "mobile-money"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}
	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)
	paymentData := resp.Data["paymentCreate"].(map[string]interface{})
	paymentID := paymentData["id"].(string)

	// Delete payment
	deleteQuery := `
		mutation {
			paymentDelete(id: $paymentId)
		}
	`
	deleteVars := map[string]interface{}{
		"paymentId": paymentID,
	}

	deleteResp := ExecuteGraphQL(t, tc, deleteQuery, deleteVars, tc.AdminToken)
	AssertNoErrors(t, deleteResp)
}

// TestPaymentCreate_PositiveAmountValidation tests positive amount validation
func TestPaymentCreate_PositiveAmountValidation(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)

	// Test with negative amount
	query := `
		mutation {
			paymentCreate(input: {
				clientId: $clientId
				amount: -100.0
				method: "mobile-money"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertHasErrors(t, resp)

	// Test with zero amount
	query2 := `
		mutation {
			paymentCreate(input: {
				clientId: $clientId
				amount: 0.0
				method: "mobile-money"
			}) {
				id
			}
		}
	`

	resp2 := ExecuteGraphQL(t, tc, query2, variables, tc.AdminToken)
	AssertHasErrors(t, resp2)
}


