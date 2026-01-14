package tests

import (
	"testing"
)

// TestProductCreate_ValidData tests product creation with valid data
func TestProductCreate_ValidData(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	query := `
		mutation {
			productCreate(input: {
				name: "Test Product"
				description: "Test product description"
				price: 100.0
				stock: 50
				points: 10.0
				imageUrl: "https://example.com/image.jpg"
			}) {
				id
				name
				description
				price
				stock
				points
				imageUrl
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["productCreate"].(map[string]interface{})
	if data["id"] == nil {
		t.Error("Product ID should not be nil")
	}
	if data["name"].(string) != "Test Product" {
		t.Error("Product name should match")
	}
}

// TestProductCreate_InvalidData tests product creation with invalid data
func TestProductCreate_InvalidData(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Test with negative price
	query := `
		mutation {
			productCreate(input: {
				name: "Test Product"
				description: "Test description"
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

	// Test with negative stock
	query2 := `
		mutation {
			productCreate(input: {
				name: "Test Product"
				description: "Test description"
				price: 100.0
				stock: -50
				points: 10.0
				imageUrl: "https://example.com/image.jpg"
			}) {
				id
			}
		}
	`

	resp2 := ExecuteGraphQL(t, tc, query2, nil, tc.AdminToken)
	AssertHasErrors(t, resp2)
}

// TestProducts_ListWithPagination tests listing products with pagination
func TestProducts_ListWithPagination(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create multiple products
	CreateTestProduct(t, tc, "Product 1")
	CreateTestProduct(t, tc, "Product 2")
	CreateTestProduct(t, tc, "Product 3")

	query := `
		query {
			products(paging: {
				page: 1
				limit: 2
			}) {
				id
				name
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	products := resp.Data["products"].([]interface{})
	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}
}

// TestProducts_ListWithFilters tests listing products with filters
func TestProducts_ListWithFilters(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create products
	CreateTestProduct(t, tc, "Test Product")
	CreateTestProduct(t, tc, "Another Product")

	query := `
		query {
			products(filter: {
				search: "Test"
			}) {
				id
				name
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	products := resp.Data["products"].([]interface{})
	if len(products) == 0 {
		t.Error("Should find at least one product")
	}
}

// TestProduct_GetByID tests getting a product by ID
func TestProduct_GetByID(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	productID := CreateTestProduct(t, tc, "Test Product")

	query := `
		query {
			product(id: $productId) {
				id
				name
				description
				price
				stock
			}
		}
	`
	variables := map[string]interface{}{
		"productId": productID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	product := resp.Data["product"].(map[string]interface{})
	if product["id"].(string) != productID {
		t.Error("Product ID should match")
	}
}

// TestProductUpdate_ValidData tests product update with valid data
func TestProductUpdate_ValidData(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	productID := CreateTestProduct(t, tc, "Original Name")

	query := `
		mutation {
			productUpdate(id: $productId, input: {
				name: "Updated Name"
				description: "Updated description"
				price: 150.0
				stock: 75
				points: 15.0
				imageUrl: "https://example.com/updated.jpg"
			}) {
				id
				name
				price
				stock
			}
		}
	`
	variables := map[string]interface{}{
		"productId": productID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["productUpdate"].(map[string]interface{})
	if data["name"].(string) != "Updated Name" {
		t.Error("Product name should be updated")
	}
	if data["price"].(float64) != 150.0 {
		t.Error("Product price should be updated")
	}
}

// TestProductUpdate_NonExistentID tests product update with non-existent ID
func TestProductUpdate_NonExistentID(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	nonExistentID := "507f1f77bcf86cd799439011"

	query := `
		mutation {
			productUpdate(id: $productId, input: {
				name: "Updated Name"
				description: "Updated description"
				price: 150.0
				stock: 75
				points: 15.0
				imageUrl: "https://example.com/updated.jpg"
			}) {
				id
			}
		}
	`
	variables := map[string]interface{}{
		"productId": nonExistentID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertHasErrors(t, resp)
}

// TestProductDelete_ValidID tests product deletion with valid ID
func TestProductDelete_ValidID(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	productID := CreateTestProduct(t, tc, "Product to Delete")

	query := `
		mutation {
			productDelete(id: $productId)
		}
	`
	variables := map[string]interface{}{
		"productId": productID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	// Verify product is deleted
	getQuery := `
		query {
			product(id: $productId) {
				id
			}
		}
	`
	getResp := ExecuteGraphQL(t, tc, getQuery, variables, tc.AdminToken)
	AssertHasErrors(t, getResp)
}

// TestProductDelete_NonExistentID tests product deletion with non-existent ID
func TestProductDelete_NonExistentID(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	nonExistentID := "507f1f77bcf86cd799439011"

	query := `
		mutation {
			productDelete(id: $productId)
		}
	`
	variables := map[string]interface{}{
		"productId": nonExistentID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertHasErrors(t, resp)
}

// TestProductCreate_RequiredFields tests product creation with missing required fields
func TestProductCreate_RequiredFields(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Missing name
	query := `
		mutation {
			productCreate(input: {
				description: "Test description"
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
	AssertHasErrors(t, resp)
}


