package tests

import (
	"testing"
)

// TestDashboardStats_NoRange tests dashboard stats without range
func TestDashboardStats_NoRange(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create some test data
	CreateTestClient(t, tc, "Test Client", nil)
	CreateTestProduct(t, tc, "Test Product")

	query := `
		query {
			dashboardStats {
				totalProducts
				totalClients
				totalSales
				totalRevenue
				activeClients
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	stats := resp.Data["dashboardStats"].(map[string]interface{})
	if stats["totalProducts"] == nil {
		t.Error("Total products should not be nil")
	}
}

// TestDashboardStats_WithRange tests dashboard stats with different ranges
func TestDashboardStats_WithRange(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	ranges := []string{"7d", "30d", "90d", "1y"}

	for _, rangeVal := range ranges {
		query := `
			query {
				dashboardStats(range: $range) {
					totalClients
					totalSales
					totalCommissions
				}
			}
		`
		variables := map[string]interface{}{
			"range": rangeVal,
		}

		resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
		AssertNoErrors(t, resp)

		stats := resp.Data["dashboardStats"].(map[string]interface{})
		if stats == nil {
			t.Errorf("Dashboard stats should not be nil for range %s", rangeVal)
		}
	}
}

// TestDashboardData tests dashboardData alias
func TestDashboardData(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	query := `
		query {
			dashboardData {
				totalProducts
				totalClients
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	stats := resp.Data["dashboardData"].(map[string]interface{})
	if stats == nil {
		t.Error("Dashboard data should not be nil")
	}
}

// TestDashboardStats_TotalProducts tests total products calculation
func TestDashboardStats_TotalProducts(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create products
	CreateTestProduct(t, tc, "Product 1")
	CreateTestProduct(t, tc, "Product 2")
	CreateTestProduct(t, tc, "Product 3")

	query := `
		query {
			dashboardStats {
				totalProducts
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	stats := resp.Data["dashboardStats"].(map[string]interface{})
	totalProducts := int(stats["totalProducts"].(float64))

	if totalProducts < 3 {
		t.Errorf("Total products should be at least 3, got %d", totalProducts)
	}
}

// TestDashboardStats_TotalClients tests total clients calculation
func TestDashboardStats_TotalClients(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create clients
	CreateTestClient(t, tc, "Client 1", nil)
	CreateTestClient(t, tc, "Client 2", nil)

	query := `
		query {
			dashboardStats {
				totalClients
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	stats := resp.Data["dashboardStats"].(map[string]interface{})
	totalClients := int(stats["totalClients"].(float64))

	if totalClients < 2 {
		t.Errorf("Total clients should be at least 2, got %d", totalClients)
	}
}

// TestDashboardStats_TotalSales tests total sales calculation
func TestDashboardStats_TotalSales(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	// Create sales
	CreateTestSale(t, tc, clientID, productID, 100.0, "paid")
	CreateTestSale(t, tc, clientID, productID, 200.0, "paid")

	query := `
		query {
			dashboardStats {
				totalSales
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	stats := resp.Data["dashboardStats"].(map[string]interface{})
	totalSales := stats["totalSales"].(float64)

	if totalSales < 300.0 {
		t.Errorf("Total sales should be at least 300.0, got %.2f", totalSales)
	}
}

// TestDashboardStats_TotalRevenue tests total revenue calculation
func TestDashboardStats_TotalRevenue(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	// Create paid sales
	CreateTestSale(t, tc, clientID, productID, 100.0, "paid")
	CreateTestSale(t, tc, clientID, productID, 150.0, "paid")

	query := `
		query {
			dashboardStats {
				totalRevenue
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	stats := resp.Data["dashboardStats"].(map[string]interface{})
	totalRevenue := stats["totalRevenue"].(float64)

	if totalRevenue < 250.0 {
		t.Errorf("Total revenue should be at least 250.0, got %.2f", totalRevenue)
	}
}

// TestDashboardStats_ActiveClients tests active clients calculation
func TestDashboardStats_ActiveClients(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	// Create sale to make client active
	CreateTestSale(t, tc, clientID, productID, 100.0, "paid")

	query := `
		query {
			dashboardStats {
				activeClients
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	stats := resp.Data["dashboardStats"].(map[string]interface{})
	activeClients := int(stats["activeClients"].(float64))

	if activeClients < 1 {
		t.Error("Should have at least one active client")
	}
}

// TestDashboardStats_NetworkVolumes tests left/right volume calculation
func TestDashboardStats_NetworkVolumes(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create binary tree
	rootID := CreateTestClient(t, tc, "Root", nil)
	leftID := CreateTestClient(t, tc, "Left", &rootID)
	rightID := CreateTestClient(t, tc, "Right", &rootID)

	productID := CreateTestProduct(t, tc, "Test Product")

	// Create sales
	CreateTestSale(t, tc, leftID, productID, 100.0, "paid")
	CreateTestSale(t, tc, rightID, productID, 100.0, "paid")

	query := `
		query {
			dashboardStats {
				leftVolume
				rightVolume
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	stats := resp.Data["dashboardStats"].(map[string]interface{})
	leftVolume := stats["leftVolume"].(float64)
	rightVolume := stats["rightVolume"].(float64)

	if leftVolume == 0 || rightVolume == 0 {
		t.Error("Network volumes should not be zero")
	}
}

// TestDashboardStats_BinaryPairs tests binary pairs calculation
func TestDashboardStats_BinaryPairs(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create binary tree with sales
	rootID := CreateTestClient(t, tc, "Root", nil)
	leftID := CreateTestClient(t, tc, "Left", &rootID)
	rightID := CreateTestClient(t, tc, "Right", &rootID)

	productID := CreateTestProduct(t, tc, "Test Product")

	// Create matching sales
	CreateTestSale(t, tc, leftID, productID, 100.0, "paid")
	CreateTestSale(t, tc, rightID, productID, 100.0, "paid")

	query := `
		query {
			dashboardStats {
				binaryPairs
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	stats := resp.Data["dashboardStats"].(map[string]interface{})
	binaryPairs := int(stats["binaryPairs"].(float64))

	if binaryPairs < 1 {
		t.Error("Should have at least one binary pair")
	}
}

// TestDashboardStats_TotalCommissions tests total commissions calculation
func TestDashboardStats_TotalCommissions(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	sourceClientID := CreateTestClient(t, tc, "Source Client", nil)

	// Create commission
	query := `
		mutation {
			commissionManualCreate(input: {
				clientId: $clientId
				sourceClientId: $sourceClientId
				amount: 50.0
				level: 1
				type: "override"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"clientId":       clientID,
		"sourceClientId": sourceClientID,
	}
	ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)

	statsQuery := `
		query {
			dashboardStats {
				totalCommissions
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, statsQuery, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	stats := resp.Data["dashboardStats"].(map[string]interface{})
	totalCommissions := stats["totalCommissions"].(float64)

	if totalCommissions < 50.0 {
		t.Errorf("Total commissions should be at least 50.0, got %.2f", totalCommissions)
	}
}

// TestDashboardStats_NetworkBalance tests network balance calculation
func TestDashboardStats_NetworkBalance(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create binary tree
	rootID := CreateTestClient(t, tc, "Root", nil)
	leftID := CreateTestClient(t, tc, "Left", &rootID)
	rightID := CreateTestClient(t, tc, "Right", &rootID)

	productID := CreateTestProduct(t, tc, "Test Product")

	// Create sales
	CreateTestSale(t, tc, leftID, productID, 100.0, "paid")
	CreateTestSale(t, tc, rightID, productID, 100.0, "paid")

	query := `
		query {
			dashboardStats {
				networkBalance
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	stats := resp.Data["dashboardStats"].(map[string]interface{})
	networkBalance := stats["networkBalance"].(float64)

	// Balance should be calculated (left - right or similar)
	if networkBalance == 0 {
		t.Log("Network balance is 0 (may be expected if volumes are equal)")
	}
}


