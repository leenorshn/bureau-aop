package tests

import (
	"testing"
)

// TestSaleCreate_WithProduct tests sale creation with product
func TestSaleCreate_WithProduct(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	query := `
		mutation {
			saleCreate(input: {
				clientId: $clientId
				productId: $productId
				quantity: 2
				amount: 200.0
				status: "paid"
			}) {
				id
				clientId
				productId
				quantity
				amount
				status
			}
		}
	`
	variables := map[string]interface{}{
		"clientId":  clientID,
		"productId": productID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["saleCreate"].(map[string]interface{})
	if data["id"] == nil {
		t.Error("Sale ID should not be nil")
	}
}

// TestSaleCreate_WithoutProduct tests sale creation without product (manual sale)
func TestSaleCreate_WithoutProduct(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)

	query := `
		mutation {
			saleCreate(input: {
				clientId: $clientId
				quantity: 1
				amount: 100.0
				status: "paid"
			}) {
				id
				amount
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)
}

// TestSaleCreate_StatusPaid tests sale creation with paid status
func TestSaleCreate_StatusPaid(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	query := `
		mutation {
			saleCreate(input: {
				clientId: $clientId
				productId: $productId
				quantity: 1
				amount: 100.0
				status: "paid"
			}) {
				id
				status
			}
		}
	`
	variables := map[string]interface{}{
		"clientId":  clientID,
		"productId": productID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["saleCreate"].(map[string]interface{})
	if data["status"].(string) != "paid" {
		t.Error("Sale status should be 'paid'")
	}
}

// TestSaleCreate_StatusPartial tests sale creation with partial status
func TestSaleCreate_StatusPartial(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	query := `
		mutation {
			saleCreate(input: {
				clientId: $clientId
				productId: $productId
				quantity: 1
				amount: 100.0
				paidAmount: 50.0
				status: "partial"
			}) {
				id
				status
				paidAmount
			}
		}
	`
	variables := map[string]interface{}{
		"clientId":  clientID,
		"productId": productID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["saleCreate"].(map[string]interface{})
	if data["status"].(string) != "partial" {
		t.Error("Sale status should be 'partial'")
	}
}

// TestSaleCreate_StatusPending tests sale creation with pending status
func TestSaleCreate_StatusPending(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	query := `
		mutation {
			saleCreate(input: {
				clientId: $clientId
				productId: $productId
				quantity: 1
				amount: 100.0
				status: "pending"
			}) {
				id
				status
			}
		}
	`
	variables := map[string]interface{}{
		"clientId":  clientID,
		"productId": productID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)
}

// TestSaleCreate_InsufficientStock tests sale creation with insufficient stock
func TestSaleCreate_InsufficientStock(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	
	// Create product with low stock
	query := `
		mutation {
			productCreate(input: {
				name: "Low Stock Product"
				description: "Test"
				price: 100.0
				stock: 5
				points: 10.0
				imageUrl: "https://example.com/image.jpg"
			}) {
				id
			}
		}
	`
	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)
	productData := resp.Data["productCreate"].(map[string]interface{})
	productID := productData["id"].(string)

	// Try to buy more than available
	saleQuery := `
		mutation {
			saleCreate(input: {
				clientId: $clientId
				productId: $productId
				quantity: 10
				amount: 1000.0
				status: "paid"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId":  clientID,
		"productId": productID,
	}

	saleResp := ExecuteGraphQL(t, tc, saleQuery, variables, tc.AdminToken)
	AssertHasErrors(t, saleResp)
}

// TestSaleCreate_PointsAddition tests automatic points addition to client
func TestSaleCreate_PointsAddition(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	// Get initial points
	query := `
		query {
			client(id: $clientId) {
				points
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}
	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)
	clientData := resp.Data["client"].(map[string]interface{})
	initialPoints := clientData["points"].(float64)

	// Create sale
	CreateTestSale(t, tc, clientID, productID, 100.0, "paid")

	// Check points increased
	resp2 := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp2)
	clientData2 := resp2.Data["client"].(map[string]interface{})
	newPoints := clientData2["points"].(float64)

	if newPoints <= initialPoints {
		t.Error("Client points should increase after sale")
	}
}

// TestSaleCreate_StockUpdate tests product stock update after sale
func TestSaleCreate_StockUpdate(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	// Get initial stock
	query := `
		query {
			product(id: $productId) {
				stock
			}
		}
	`
	variables := map[string]interface{}{
		"productId": productID,
	}
	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)
	productData := resp.Data["product"].(map[string]interface{})
	initialStock := int(productData["stock"].(float64))

	// Create sale with quantity 2
	CreateTestSale(t, tc, clientID, productID, 200.0, "paid")

	// Check stock decreased
	resp2 := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp2)
	productData2 := resp2.Data["product"].(map[string]interface{})
	newStock := int(productData2["stock"].(float64))

	if newStock != initialStock-1 {
		t.Errorf("Stock should decrease by 1, expected %d, got %d", initialStock-1, newStock)
	}
}

// TestSaleCreate_CaisseEntry tests caisse entry for paid sale
func TestSaleCreate_CaisseEntry(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	// Get initial caisse balance
	query := `
		query {
			caisse {
				balance
				totalEntrees
			}
		}
	`
	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)
	caisseData := resp.Data["caisse"].(map[string]interface{})
	initialBalance := caisseData["balance"].(float64)

	// Create paid sale
	CreateTestSale(t, tc, clientID, productID, 100.0, "paid")

	// Check caisse balance increased
	resp2 := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp2)
	caisseData2 := resp2.Data["caisse"].(map[string]interface{})
	newBalance := caisseData2["balance"].(float64)

	if newBalance <= initialBalance {
		t.Error("Caisse balance should increase after paid sale")
	}
}

// TestSales_ListWithFilters tests listing sales with filters
func TestSales_ListWithFilters(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")
	CreateTestSale(t, tc, clientID, productID, 100.0, "paid")

	query := `
		query {
			sales(filter: {
				status: "paid"
			}) {
				id
				amount
				status
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	sales := resp.Data["sales"].([]interface{})
	if len(sales) == 0 {
		t.Error("Should find at least one sale")
	}
}

// TestSale_GetByID tests getting a sale by ID
func TestSale_GetByID(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")
	saleID := CreateTestSale(t, tc, clientID, productID, 100.0, "paid")

	query := `
		query {
			sale(id: $saleId) {
				id
				amount
				status
				client {
					id
					name
				}
			}
		}
	`
	variables := map[string]interface{}{
		"saleId": saleID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	sale := resp.Data["sale"].(map[string]interface{})
	if sale["id"].(string) != saleID {
		t.Error("Sale ID should match")
	}
}

// TestSaleUpdate tests updating a sale
func TestSaleUpdate(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")
	saleID := CreateTestSale(t, tc, clientID, productID, 100.0, "pending")

	query := `
		mutation {
			saleUpdate(id: $saleId, input: {
				clientId: $clientId
				productId: $productId
				quantity: 1
				amount: 100.0
				status: "paid"
			}) {
				id
				status
			}
		}
	`
	variables := map[string]interface{}{
		"saleId":    saleID,
		"clientId":  clientID,
		"productId": productID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["saleUpdate"].(map[string]interface{})
	if data["status"].(string) != "paid" {
		t.Error("Sale status should be updated to 'paid'")
	}
}

// TestSaleDelete tests deleting a sale
func TestSaleDelete(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")
	saleID := CreateTestSale(t, tc, clientID, productID, 100.0, "pending")

	query := `
		mutation {
			saleDelete(id: $saleId)
		}
	`
	variables := map[string]interface{}{
		"saleId": saleID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)
}

// TestSaleCreate_PartialPaidAmountValidation tests paidAmount validation for partial status
func TestSaleCreate_PartialPaidAmountValidation(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	// Test without paidAmount for partial status
	query := `
		mutation {
			saleCreate(input: {
				clientId: $clientId
				productId: $productId
				quantity: 1
				amount: 100.0
				status: "partial"
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


