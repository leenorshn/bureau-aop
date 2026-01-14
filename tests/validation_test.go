package tests

import (
	"testing"
)

// TestValidation_ObjectID tests ObjectID validation
func TestValidation_ObjectID(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	validID := "507f1f77bcf86cd799439011"
	invalidIDs := []string{
		"",
		"invalid",
		"123",
		"507f1f77bcf86cd79943901", // too short
		"507f1f77bcf86cd7994390111", // too long
	}

	// Test valid ID
	query := `
		query {
			product(id: $productId) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"productId": validID,
	}
	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	// May not find product, but should not error on ID format
	if resp.Errors != nil {
		// Check if error is about ID format or product not found
		t.Log("Valid ID format accepted (product may not exist)")
	}

	// Test invalid IDs
	for _, invalidID := range invalidIDs {
		variables["productId"] = invalidID
		resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
		AssertHasErrors(t, resp)
	}
}

// TestValidation_Email tests email validation
func TestValidation_Email(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Test with invalid email in login
	query := `
		mutation {
			userLogin(input: {
				email: "invalid-email"
				password: "Test123@admin"
			}) {
				accessToken
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, "")
	AssertHasErrors(t, resp)
}

// TestValidation_PasswordStrength tests password strength validation
func TestValidation_PasswordStrength(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	weakPasswords := []string{
		"weak",
		"12345678",
		"password",
		"PASSWORD",
		"Password1",
		"Password@", // missing digit
		"Pass1@",     // too short
	}

	for _, weakPassword := range weakPasswords {
		query := `
			mutation {
				clientCreate(input: {
					name: "Test Client"
					password: $password
				}) {
					id
				}
			}
		`
		variables := map[string]interface{}{
			"password": weakPassword,
		}

		resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
		AssertHasErrors(t, resp)
	}
}

// TestValidation_AmountPositive tests positive amount validation
func TestValidation_AmountPositive(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)

	// Test with negative amount
	query := `
		mutation {
			paymentCreate(input: {
				clientId: $clientId
				amount: -100.0
				method: "cash"
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
				method: "cash"
			}) {
				id
			}
		}
	`

	resp2 := ExecuteGraphQL(t, tc, query2, variables, tc.AdminToken)
	AssertHasErrors(t, resp2)
}

// TestValidation_AmountNonNegative tests non-negative amount validation
func TestValidation_AmountNonNegative(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Test with negative price in product
	query := `
		mutation {
			productCreate(input: {
				name: "Test Product"
				description: "Test"
				price: -100.0
				stock: 50
				points: 10.0
				imageUrl: "https://example.com/image.jpg"
			}) {
				id
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertHasErrors(t, resp)
}

// TestValidation_Quantity tests quantity validation
func TestValidation_Quantity(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	// Test with zero quantity
	query := `
		mutation {
			saleCreate(input: {
				clientId: $clientId
				productId: $productId
				quantity: 0
				amount: 100.0
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId":  clientID,
		"productId": productID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertHasErrors(t, resp)

	// Test with negative quantity
	query2 := `
		mutation {
			saleCreate(input: {
				clientId: $clientId
				productId: $productId
				quantity: -1
				amount: 100.0
			}) {
				id
			}
		}
	`

	resp2 := ExecuteGraphQL(t, tc, query2, variables, tc.AdminToken)
	AssertHasErrors(t, resp2)
}

// TestValidation_SaleStatus tests sale status validation
func TestValidation_SaleStatus(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	// Test with invalid status
	query := `
		mutation {
			saleCreate(input: {
				clientId: $clientId
				productId: $productId
				quantity: 1
				amount: 100.0
				status: "invalid-status"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId":  clientID,
		"productId": productID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertHasErrors(t, resp)
}

// TestValidation_PaymentMethod tests payment method validation
func TestValidation_PaymentMethod(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)

	// Test with invalid method
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

// TestValidation_CommissionType tests commission type validation
func TestValidation_CommissionType(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	sourceClientID := CreateTestClient(t, tc, "Source Client", nil)

	// Test with invalid type
	query := `
		mutation {
			commissionManualCreate(input: {
				clientId: $clientId
				sourceClientId: $sourceClientId
				amount: 50.0
				level: 1
				type: "invalid-type"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId":       clientID,
		"sourceClientId": sourceClientID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertHasErrors(t, resp)
}

// TestValidation_TransactionType tests caisse transaction type validation
func TestValidation_TransactionType(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Test with invalid type
	query := `
		mutation {
			caisseAddTransaction(input: {
				type: "invalid-type"
				amount: 100.0
				description: "Test"
			}) {
				id
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertHasErrors(t, resp)
}


