package tests

import (
	"testing"
)

// TestClientCreate_RootClient tests creating a root client without sponsor
func TestClientCreate_RootClient(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	query := `
		mutation {
			clientCreate(input: {
				name: "Root Client"
				password: "Test123@client"
			}) {
				id
				clientId
				name
				sponsorId
				position
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["clientCreate"].(map[string]interface{})
	if data["id"] == nil {
		t.Error("Client ID should not be nil")
	}
	if data["sponsorId"] != nil {
		t.Error("Root client should not have sponsor")
	}
}

// TestClientCreate_WithSponsor tests creating a client with sponsor
func TestClientCreate_WithSponsor(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create root client first
	rootID := CreateTestClient(t, tc, "Root Client", nil)

	query := `
		mutation {
			clientCreate(input: {
				name: "Child Client"
				password: "Test123@client"
				sponsorId: $sponsorId
			}) {
				id
				clientId
				name
				sponsorId
				position
			}
		}
	`
	variables := map[string]interface{}{
		"sponsorId": rootID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["clientCreate"].(map[string]interface{})
	if data["sponsorId"] == nil {
		t.Error("Child client should have sponsor")
	}
	if data["position"] == nil {
		t.Error("Child client should have position")
	}
}

// TestClientCreate_WithPosition tests creating a client with specified position
func TestClientCreate_WithPosition(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create root client
	rootID := CreateTestClient(t, tc, "Root Client", nil)

	query := `
		mutation {
			clientCreate(input: {
				name: "Left Client"
				password: "Test123@client"
				sponsorId: $sponsorId
				position: "left"
			}) {
				id
				position
			}
		}
	`
	variables := map[string]interface{}{
		"sponsorId": rootID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["clientCreate"].(map[string]interface{})
	if data["position"].(string) != "left" {
		t.Error("Client position should be 'left'")
	}
}

// TestClientCreate_PasswordValidation tests password validation
func TestClientCreate_PasswordValidation(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Test with weak password
	query := `
		mutation {
			clientCreate(input: {
				name: "Test Client"
				password: "weak"
			}) {
				id
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertHasErrors(t, resp)
}

// TestClients_ListWithPagination tests listing clients with pagination
func TestClients_ListWithPagination(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create multiple clients
	CreateTestClient(t, tc, "Client 1", nil)
	CreateTestClient(t, tc, "Client 2", nil)
	CreateTestClient(t, tc, "Client 3", nil)

	query := `
		query {
			clients(paging: {
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

	clients := resp.Data["clients"].([]interface{})
	if len(clients) != 2 {
		t.Errorf("Expected 2 clients, got %d", len(clients))
	}
}

// TestClients_ListWithFilters tests listing clients with search filters
func TestClients_ListWithFilters(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create clients
	CreateTestClient(t, tc, "John Doe", nil)
	CreateTestClient(t, tc, "Jane Smith", nil)

	query := `
		query {
			clients(filter: {
				search: "John"
			}) {
				id
				name
			}
		}
	`

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	AssertNoErrors(t, resp)

	clients := resp.Data["clients"].([]interface{})
	if len(clients) == 0 {
		t.Error("Should find at least one client")
	}
}

// TestClient_GetByID tests getting a client by ID with relations
func TestClient_GetByID(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)

	query := `
		query {
			client(id: $clientId) {
				id
				name
				clientId
				sponsorId
				position
				totalEarnings
				walletBalance
				networkVolumeLeft
				networkVolumeRight
				binaryPairs
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	client := resp.Data["client"].(map[string]interface{})
	if client["id"].(string) != clientID {
		t.Error("Client ID should match")
	}
}

// TestClientTree tests getting client tree
func TestClientTree(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create root client
	rootID := CreateTestClient(t, tc, "Root Client", nil)

	// Create child clients
	CreateTestClient(t, tc, "Left Child", &rootID)
	CreateTestClient(t, tc, "Right Child", &rootID)

	query := `
		query {
			clientTree(id: $rootId) {
				root {
					id
					name
					clientId
				}
				nodes {
					id
					name
					position
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

	tree := resp.Data["clientTree"].(map[string]interface{})
	if tree["root"] == nil {
		t.Error("Tree root should not be nil")
	}
	nodes := tree["nodes"].([]interface{})
	if len(nodes) < 3 {
		t.Errorf("Expected at least 3 nodes, got %d", len(nodes))
	}
}

// TestClientUpdate tests updating client information
func TestClientUpdate(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Original Name", nil)

	query := `
		mutation {
			clientUpdate(id: $clientId, input: {
				name: "Updated Name"
				phone: "1234567890"
				address: "123 Test St"
			}) {
				id
				name
				phone
				address
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	data := resp.Data["clientUpdate"].(map[string]interface{})
	if data["name"].(string) != "Updated Name" {
		t.Error("Client name should be updated")
	}
}

// TestClientDelete tests deleting a client
func TestClientDelete(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Client to Delete", nil)

	query := `
		mutation {
			clientDelete(id: $clientId)
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	// Verify client is deleted
	getQuery := `
		query {
			client(id: $clientId) {
				id
			}
		}
	`
	getResp := ExecuteGraphQL(t, tc, getQuery, variables, tc.AdminToken)
	AssertHasErrors(t, getResp)
}

// TestClientCreate_AutomaticPlacement tests automatic placement in binary tree
func TestClientCreate_AutomaticPlacement(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create root
	rootID := CreateTestClient(t, tc, "Root", nil)

	// First child should go left
	_ = CreateTestClient(t, tc, "Left", &rootID)

	// Verify left placement
	query := `
		query {
			client(id: $rootId) {
				leftChildId
			}
		}
	`
	variables := map[string]interface{}{
		"rootId": rootID,
	}
	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)
	rootData := resp.Data["client"].(map[string]interface{})
	if rootData["leftChildId"] == nil {
		t.Error("Root should have left child")
	}

	// Second child should go right
	_ = CreateTestClient(t, tc, "Right", &rootID)

	// Verify right placement
	resp2 := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp2)
	rootData2 := resp2.Data["client"].(map[string]interface{})
	if rootData2["rightChildId"] == nil {
		t.Error("Root should have right child")
	}
}

// TestClientCreate_NetworkVolumeCalculation tests network volume calculation
func TestClientCreate_NetworkVolumeCalculation(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create root and children
	rootID := CreateTestClient(t, tc, "Root", nil)
	leftID := CreateTestClient(t, tc, "Left", &rootID)
	rightID := CreateTestClient(t, tc, "Right", &rootID)

	// Create product and sales to generate volume
	productID := CreateTestProduct(t, tc, "Test Product")
	CreateTestSale(t, tc, leftID, productID, 100.0, "paid")
	CreateTestSale(t, tc, rightID, productID, 100.0, "paid")

	// Check network volumes
	query := `
		query {
			client(id: $rootId) {
				networkVolumeLeft
				networkVolumeRight
			}
		}
	`
	variables := map[string]interface{}{
		"rootId": rootID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	client := resp.Data["client"].(map[string]interface{})
	leftVolume := client["networkVolumeLeft"].(float64)
	rightVolume := client["networkVolumeRight"].(float64)

	if leftVolume == 0 || rightVolume == 0 {
		t.Error("Network volumes should be updated after sales")
	}
}

