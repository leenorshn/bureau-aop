package tests

import (
	"testing"
	"time"
)

// TestPerformance_MultipleClientsCreation tests creating multiple clients
func TestPerformance_MultipleClientsCreation(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	start := time.Now()

	// Create 10 clients
	for i := 0; i < 10; i++ {
		CreateTestClient(t, tc, "Client", nil)
	}

	elapsed := time.Since(start)
	t.Logf("Created 10 clients in %v", elapsed)

	if elapsed > 5*time.Second {
		t.Errorf("Creating 10 clients took too long: %v", elapsed)
	}
}

// TestPerformance_DashboardWithLargeData tests dashboard query with large dataset
func TestPerformance_DashboardWithLargeData(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create test data
	for i := 0; i < 5; i++ {
		CreateTestClient(t, tc, "Client", nil)
		CreateTestProduct(t, tc, "Product")
	}

	start := time.Now()

	query := `
		query {
			dashboardStats {
				totalProducts
				totalClients
				totalSales
				totalCommissions
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	elapsed := time.Since(start)
	t.Logf("Dashboard query took %v", elapsed)

	if elapsed > 2*time.Second {
		t.Errorf("Dashboard query took too long: %v", elapsed)
	}
}

// TestPerformance_ClientTreeDeep tests clientTree with deep tree
func TestPerformance_ClientTreeDeep(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create deep tree (3 levels)
	rootID := CreateTestClient(t, tc, "Root", nil)
	leftID := CreateTestClient(t, tc, "Left", &rootID)
	rightID := CreateTestClient(t, tc, "Right", &rootID)
	_ = CreateTestClient(t, tc, "Left Left", &leftID)
	_ = CreateTestClient(t, tc, "Left Right", &leftID)
	_ = CreateTestClient(t, tc, "Right Left", &rightID)
	_ = CreateTestClient(t, tc, "Right Right", &rightID)

	start := time.Now()

	query := `
		query {
			clientTree(id: $rootId) {
				root {
					id
					name
				}
				nodes {
					id
					name
					level
				}
				totalNodes
				maxLevel
			}
		}
	`
	variables := map[string]interface{}{
		"rootId": rootID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	elapsed := time.Since(start)
	t.Logf("ClientTree query took %v", elapsed)

	if elapsed > 3*time.Second {
		t.Errorf("ClientTree query took too long: %v", elapsed)
	}
}

// TestPerformance_SimpleQueryResponseTime tests simple query response time
func TestPerformance_SimpleQueryResponseTime(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	start := time.Now()

	query := `
		query {
			products {
				id
				name
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	elapsed := time.Since(start)
	t.Logf("Simple query took %v", elapsed)

	if elapsed > 1*time.Second {
		t.Errorf("Simple query took too long: %v", elapsed)
	}
}

// TestPerformance_MutationResponseTime tests mutation response time
func TestPerformance_MutationResponseTime(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	start := time.Now()

	query := `
		mutation {
			productCreate(input: {
				name: "Performance Test Product"
				description: "Test"
				price: 100.0
				stock: 50
				points: 10.0
				imageUrl: "https://example.com/image.jpg"
			}) {
				id
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	elapsed := time.Since(start)
	t.Logf("Mutation took %v", elapsed)

	if elapsed > 1*time.Second {
		t.Errorf("Mutation took too long: %v", elapsed)
	}
}

// TestPerformance_PaginationLargeList tests pagination with large list
func TestPerformance_PaginationLargeList(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create multiple products
	for i := 0; i < 20; i++ {
		CreateTestProduct(t, tc, "Product")
	}

	start := time.Now()

	query := `
		query {
			products(paging: {
				page: 1
				limit: 10
			}) {
				id
				name
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	elapsed := time.Since(start)
	t.Logf("Pagination query took %v", elapsed)

	products := resp.Data["products"].([]interface{})
	if len(products) != 10 {
		t.Errorf("Should return 10 products, got %d", len(products))
	}

	if elapsed > 1*time.Second {
		t.Errorf("Pagination query took too long: %v", elapsed)
	}
}

