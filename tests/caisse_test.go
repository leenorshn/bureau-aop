package tests

import (
	"testing"
)

// TestCaisse_GetInitialState tests getting initial caisse state
func TestCaisse_GetInitialState(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	query := `
		query {
			caisse {
				id
				balance
				totalEntrees
				totalSorties
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	caisse := resp.Data["caisse"].(map[string]interface{})
	if caisse["balance"] == nil {
		t.Error("Caisse balance should not be nil")
	}
}

// TestCaisseAddTransaction_Entree tests adding entry transaction
func TestCaisseAddTransaction_Entree(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	query := `
		mutation {
			caisseAddTransaction(input: {
				type: "entree"
				amount: 100.0
				description: "Test entry"
			}) {
				id
				type
				amount
				description
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["caisseAddTransaction"].(map[string]interface{})
	if data["type"].(string) != "entree" {
		t.Error("Transaction type should be 'entree'")
	}
}

// TestCaisseAddTransaction_Sortie tests adding exit transaction
func TestCaisseAddTransaction_Sortie(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// First add some money
	entryQuery := `
		mutation {
			caisseAddTransaction(input: {
				type: "entree"
				amount: 200.0
				description: "Initial deposit"
			}) {
				id
			}
		}
	`
	ExecuteGraphQL(t, tc, entryQuery, nil, tc.AdminToken)

	// Then create exit
	query := `
		mutation {
			caisseAddTransaction(input: {
				type: "sortie"
				amount: 50.0
				description: "Test exit"
			}) {
				id
				type
				amount
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["caisseAddTransaction"].(map[string]interface{})
	if data["type"].(string) != "sortie" {
		t.Error("Transaction type should be 'sortie'")
	}
}

// TestCaisseAddTransaction_SaleReference tests transaction with sale reference
func TestCaisseAddTransaction_SaleReference(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")
	saleID := CreateTestSale(t, tc, clientID, productID, 100.0, "paid")

	query := `
		mutation {
			caisseAddTransaction(input: {
				type: "entree"
				amount: 100.0
				description: "Sale transaction"
				reference: $saleId
				referenceType: "sale"
			}) {
				id
				reference
				referenceType
			}
		}
	`
	variables := map[string]interface{}{
		"saleId": saleID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["caisseAddTransaction"].(map[string]interface{})
	if data["reference"].(string) != saleID {
		t.Error("Transaction reference should match sale ID")
	}
}

// TestCaisseAddTransaction_PaymentReference tests transaction with payment reference
func TestCaisseAddTransaction_PaymentReference(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)

	// Create payment
	paymentQuery := `
		mutation {
			paymentCreate(input: {
				clientId: $clientId
				amount: 50.0
				method: "cash"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}
	paymentResp := ExecuteGraphQL(t, tc, paymentQuery, variables, tc.AdminToken)
	AssertNoErrors(t, paymentResp)
	paymentData := paymentResp.Data["paymentCreate"].(map[string]interface{})
	paymentID := paymentData["id"].(string)

	query := `
		mutation {
			caisseAddTransaction(input: {
				type: "sortie"
				amount: 50.0
				description: "Payment transaction"
				reference: $paymentId
				referenceType: "payment"
			}) {
				id
				reference
				referenceType
			}
		}
	`
	variables["paymentId"] = paymentID

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["caisseAddTransaction"].(map[string]interface{})
	if data["referenceType"].(string) != "payment" {
		t.Error("Transaction reference type should be 'payment'")
	}
}

// TestCaisseUpdateBalance tests manual balance update
func TestCaisseUpdateBalance(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	query := `
		mutation {
			caisseUpdateBalance(balance: 500.0) {
				id
				balance
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["caisseUpdateBalance"].(map[string]interface{})
	if data["balance"].(float64) != 500.0 {
		t.Error("Caisse balance should be updated to 500.0")
	}
}

// TestCaisseTransactions_ListWithFilters tests listing transactions with filters
func TestCaisseTransactions_ListWithFilters(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create some transactions
	entryQuery := `
		mutation {
			caisseAddTransaction(input: {
				type: "entree"
				amount: 100.0
				description: "Entry 1"
			}) {
				id
			}
		}
	`
	ExecuteGraphQL(t, tc, entryQuery, nil, tc.AdminToken)

	query := `
		query {
			caisseTransactions(filter: {
				status: "entree"
			}) {
				id
				type
				amount
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	transactions := resp.Data["caisseTransactions"].([]interface{})
	if len(transactions) == 0 {
		t.Error("Should find at least one transaction")
	}
}

// TestCaisse_AutoCalculateTotals tests automatic calculation of totalEntrees/totalSorties
func TestCaisse_AutoCalculateTotals(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Add entries
	entryQuery := `
		mutation {
			caisseAddTransaction(input: {
				type: "entree"
				amount: 100.0
				description: "Entry 1"
			}) {
				id
			}
		}
	`
	ExecuteGraphQL(t, tc, entryQuery, nil, tc.AdminToken)

	entryQuery2 := `
		mutation {
			caisseAddTransaction(input: {
				type: "entree"
				amount: 50.0
				description: "Entry 2"
			}) {
				id
			}
		}
	`
	ExecuteGraphQL(t, tc, entryQuery2, nil, tc.AdminToken)

	// Add exit
	exitQuery := `
		mutation {
			caisseAddTransaction(input: {
				type: "sortie"
				amount: 30.0
				description: "Exit 1"
			}) {
				id
			}
		}
	`
	ExecuteGraphQL(t, tc, exitQuery, nil, tc.AdminToken)

	// Check totals
	query := `
		query {
			caisse {
				totalEntrees
				totalSorties
				balance
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	caisse := resp.Data["caisse"].(map[string]interface{})
	totalEntrees := caisse["totalEntrees"].(float64)
	totalSorties := caisse["totalSorties"].(float64)
	balance := caisse["balance"].(float64)

	if totalEntrees != 150.0 {
		t.Errorf("Total entrees should be 150.0, got %.2f", totalEntrees)
	}
	if totalSorties != 30.0 {
		t.Errorf("Total sorties should be 30.0, got %.2f", totalSorties)
	}
	if balance != 120.0 {
		t.Errorf("Balance should be 120.0, got %.2f", balance)
	}
}

// TestCaisse_BalanceCalculation tests balance = totalEntrees - totalSorties
func TestCaisse_BalanceCalculation(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Add transactions
	entryQuery := `
		mutation {
			caisseAddTransaction(input: {
				type: "entree"
				amount: 200.0
				description: "Entry"
			}) {
				id
			}
		}
	`
	ExecuteGraphQL(t, tc, entryQuery, nil, tc.AdminToken)

	exitQuery := `
		mutation {
			caisseAddTransaction(input: {
				type: "sortie"
				amount: 75.0
				description: "Exit"
			}) {
				id
			}
		}
	`
	ExecuteGraphQL(t, tc, exitQuery, nil, tc.AdminToken)

	// Verify balance calculation
	query := `
		query {
			caisse {
				balance
				totalEntrees
				totalSorties
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	caisse := resp.Data["caisse"].(map[string]interface{})
	balance := caisse["balance"].(float64)
	totalEntrees := caisse["totalEntrees"].(float64)
	totalSorties := caisse["totalSorties"].(float64)

	expectedBalance := totalEntrees - totalSorties
	if balance != expectedBalance {
		t.Errorf("Balance should be %.2f, got %.2f", expectedBalance, balance)
	}
}

// TestCaisse_SaleLinkedTransaction tests transaction linked to sale
func TestCaisse_SaleLinkedTransaction(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	// Create paid sale (should auto-create transaction)
	CreateTestSale(t, tc, clientID, productID, 100.0, "paid")

	// Check transaction exists
	query := `
		query {
			caisseTransactions {
				id
				type
				amount
				referenceType
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	transactions := resp.Data["caisseTransactions"].([]interface{})
	if len(transactions) == 0 {
		t.Error("Should have transaction from sale")
	}
}

// TestCaisse_PaymentLinkedTransaction tests transaction linked to payment
func TestCaisse_PaymentLinkedTransaction(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// First add money to caisse
	entryQuery := `
		mutation {
			caisseAddTransaction(input: {
				type: "entree"
				amount: 200.0
				description: "Initial"
			}) {
				id
			}
		}
	`
	ExecuteGraphQL(t, tc, entryQuery, nil, tc.AdminToken)

	clientID := CreateTestClient(t, tc, "Test Client", nil)

	// Create payment (should auto-create transaction)
	paymentQuery := `
		mutation {
			paymentCreate(input: {
				clientId: $clientId
				amount: 50.0
				method: "cash"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}
	ExecuteGraphQL(t, tc, paymentQuery, variables, tc.AdminToken)

	// Check transaction exists
	query := `
		query {
			caisseTransactions {
				id
				type
				amount
				referenceType
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	transactions := resp.Data["caisseTransactions"].([]interface{})
	// Should have at least entry + payment exit
	if len(transactions) < 2 {
		t.Error("Should have transactions from entry and payment")
	}
}


