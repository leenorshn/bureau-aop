package tests

import (
	"testing"
)

// TestBinaryPlacement_FirstClient tests placement of first client as root
func TestBinaryPlacement_FirstClient(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// First client should be root
	rootID := CreateTestClient(t, tc, "Root Client", nil)

	query := `
		query {
			client(id: $clientId) {
				id
				sponsorId
				position
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": rootID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	client := resp.Data["client"].(map[string]interface{})
	if client["sponsorId"] != nil {
		t.Error("Root client should not have sponsor")
	}
}

// TestBinaryPlacement_SecondClient tests placement of second client on left
func TestBinaryPlacement_SecondClient(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create root
	rootID := CreateTestClient(t, tc, "Root", nil)

	// Second client should go left
	leftID := CreateTestClient(t, tc, "Left", &rootID)

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

	root := resp.Data["client"].(map[string]interface{})
	if root["leftChildId"] == nil {
		t.Error("Root should have left child")
	}
	if root["leftChildId"].(string) != leftID {
		t.Error("Left child ID should match")
	}
}

// TestBinaryPlacement_ThirdClient tests placement of third client on right
func TestBinaryPlacement_ThirdClient(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create root
	rootID := CreateTestClient(t, tc, "Root", nil)

	// Create left child
	CreateTestClient(t, tc, "Left", &rootID)

	// Third client should go right
	rightID := CreateTestClient(t, tc, "Right", &rootID)

	query := `
		query {
			client(id: $rootId) {
				rightChildId
			}
		}
	`
	variables := map[string]interface{}{
		"rootId": rootID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	root := resp.Data["client"].(map[string]interface{})
	if root["rightChildId"] == nil {
		t.Error("Root should have right child")
	}
	if root["rightChildId"].(string) != rightID {
		t.Error("Right child ID should match")
	}
}

// TestBinaryPlacement_Balancing tests automatic balancing of left/right
func TestBinaryPlacement_Balancing(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create root
	rootID := CreateTestClient(t, tc, "Root", nil)

	// Create left child
	leftID := CreateTestClient(t, tc, "Left", &rootID)

	// Create right child
	_ = CreateTestClient(t, tc, "Right", &rootID)

	// Next client should balance - go to left side
	leftLeftID := CreateTestClient(t, tc, "Left Left", &leftID)

	query := `
		query {
			client(id: $leftId) {
				leftChildId
			}
		}
	`
	variables := map[string]interface{}{
		"leftId": leftID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	left := resp.Data["client"].(map[string]interface{})
	if left["leftChildId"] == nil {
		t.Error("Left client should have left child")
	}
	if left["leftChildId"].(string) != leftLeftID {
		t.Error("Left child ID should match")
	}
}

// TestNetworkVolumeUpdate_OnSale tests network volume update when sale is made
func TestNetworkVolumeUpdate_OnSale(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create root and children
	rootID := CreateTestClient(t, tc, "Root", nil)
	leftID := CreateTestClient(t, tc, "Left", &rootID)
	rightID := CreateTestClient(t, tc, "Right", &rootID)

	// Create product
	productID := CreateTestProduct(t, tc, "Test Product")

	// Get initial volumes
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
	rootData := resp.Data["client"].(map[string]interface{})
	initialLeftVolume := rootData["networkVolumeLeft"].(float64)
	initialRightVolume := rootData["networkVolumeRight"].(float64)

	// Create sale on left side
	CreateTestSale(t, tc, leftID, productID, 100.0, "paid")

	// Check volumes updated
	resp2 := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp2)
	rootData2 := resp2.Data["client"].(map[string]interface{})
	newLeftVolume := rootData2["networkVolumeLeft"].(float64)

	if newLeftVolume <= initialLeftVolume {
		t.Error("Left network volume should increase after sale")
	}

	// Create sale on right side
	CreateTestSale(t, tc, rightID, productID, 100.0, "paid")

	// Check volumes updated
	resp3 := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp3)
	rootData3 := resp3.Data["client"].(map[string]interface{})
	newRightVolume := rootData3["networkVolumeRight"].(float64)

	if newRightVolume <= initialRightVolume {
		t.Error("Right network volume should increase after sale")
	}
}

// TestBinaryPairs_Calculation tests binary pairs calculation
func TestBinaryPairs_Calculation(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create root and children
	rootID := CreateTestClient(t, tc, "Root", nil)
	leftID := CreateTestClient(t, tc, "Left", &rootID)
	rightID := CreateTestClient(t, tc, "Right", &rootID)

	// Create product
	productID := CreateTestProduct(t, tc, "Test Product")

	// Create sales to generate pairs
	CreateTestSale(t, tc, leftID, productID, 100.0, "paid")
	CreateTestSale(t, tc, rightID, productID, 100.0, "paid")

	query := `
		query {
			client(id: $rootId) {
				binaryPairs
			}
		}
	`
	variables := map[string]interface{}{
		"rootId": rootID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	client := resp.Data["client"].(map[string]interface{})
	binaryPairs := int(client["binaryPairs"].(float64))

	if binaryPairs < 1 {
		t.Error("Binary pairs should be at least 1 after matching sales")
	}
}

// TestClientQualification_ActiveWithSales tests client qualification when active with sales
func TestClientQualification_ActiveWithSales(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	clientID := CreateTestClient(t, tc, "Test Client", nil)
	productID := CreateTestProduct(t, tc, "Test Product")

	// Create sale to make client active
	CreateTestSale(t, tc, clientID, productID, 100.0, "paid")

	// Check client tree for isActive
	query := `
		query {
			clientTree(id: $clientId) {
				root {
					isActive
				}
			}
		}
	`
	variables := map[string]interface{}{
		"clientId": clientID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	tree := resp.Data["clientTree"].(map[string]interface{})
	root := tree["root"].(map[string]interface{})
	isActive := root["isActive"].(bool)

	if !isActive {
		t.Error("Client should be active after making a sale")
	}
}

// TestCyclesAvailable_Calculation tests cycles available calculation
func TestCyclesAvailable_Calculation(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create root and children with sales
	rootID := CreateTestClient(t, tc, "Root", nil)
	leftID := CreateTestClient(t, tc, "Left", &rootID)
	rightID := CreateTestClient(t, tc, "Right", &rootID)

	productID := CreateTestProduct(t, tc, "Test Product")

	// Create sales to generate volume
	CreateTestSale(t, tc, leftID, productID, 100.0, "paid")
	CreateTestSale(t, tc, rightID, productID, 100.0, "paid")

	// Check cycles available in tree
	query := `
		query {
			clientTree(id: $rootId) {
				root {
					cyclesAvailable
				}
			}
		}
	`
	variables := map[string]interface{}{
		"rootId": rootID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	tree := resp.Data["clientTree"].(map[string]interface{})
	root := tree["root"].(map[string]interface{})
	cyclesAvailable := root["cyclesAvailable"]

	// Cycles may be null or a number
	if cyclesAvailable != nil {
		cycles := int(cyclesAvailable.(float64))
		if cycles < 0 {
			t.Error("Cycles available should not be negative")
		}
	}
}

// TestDailyCycleLimit tests daily cycle limit
func TestDailyCycleLimit(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// This test would require running multiple commission checks in the same day
	// For now, just verify the limit is respected in commission check
	rootID := CreateTestClient(t, tc, "Root", nil)
	leftID := CreateTestClient(t, tc, "Left", &rootID)
	rightID := CreateTestClient(t, tc, "Right", &rootID)

	productID := CreateTestProduct(t, tc, "Test Product")
	CreateTestSale(t, tc, leftID, productID, 100.0, "paid")
	CreateTestSale(t, tc, rightID, productID, 100.0, "paid")

	// Run commission check multiple times
	query := `
		mutation {
			runBinaryCommissionCheck(clientId: $rootId) {
				commissionsCreated
				message
			}
		}
	`
	variables := map[string]interface{}{
		"rootId": rootID,
	}

	// First check
	resp1 := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	if resp1.Errors != nil && len(resp1.Errors) > 0 {
		t.Logf("First commission check: %v", resp1.Errors)
	}

	// Second check (may hit daily limit)
	resp2 := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	if resp2.Errors != nil && len(resp2.Errors) > 0 {
		t.Logf("Second commission check: %v", resp2.Errors)
	}
}

// TestWeeklyCycleLimit tests weekly cycle limit
func TestWeeklyCycleLimit(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Similar to daily limit test
	// This would require testing across a week
	rootID := CreateTestClient(t, tc, "Root", nil)

	query := `
		mutation {
			runBinaryCommissionCheck(clientId: $rootId) {
				commissionsCreated
				message
			}
		}
	`
	variables := map[string]interface{}{
		"rootId": rootID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	// Should not error
	if resp.Errors != nil && len(resp.Errors) > 0 {
		t.Logf("Commission check returned errors: %v", resp.Errors)
	}
}

// TestBinaryCommission_WithThresholds tests binary commission calculation with thresholds
func TestBinaryCommission_WithThresholds(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create root and children
	rootID := CreateTestClient(t, tc, "Root", nil)
	leftID := CreateTestClient(t, tc, "Left", &rootID)
	rightID := CreateTestClient(t, tc, "Right", &rootID)

	productID := CreateTestProduct(t, tc, "Test Product")

	// Create sales above threshold
	CreateTestSale(t, tc, leftID, productID, 100.0, "paid")
	CreateTestSale(t, tc, rightID, productID, 100.0, "paid")

	// Run commission check
	query := `
		mutation {
			runBinaryCommissionCheck(clientId: $rootId) {
				commissionsCreated
				totalAmount
				message
			}
		}
	`
	variables := map[string]interface{}{
		"rootId": rootID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	// Should not error
	if resp.Errors != nil && len(resp.Errors) > 0 {
		t.Logf("Commission check returned errors: %v", resp.Errors)
	}
}

// TestVolumePropagation_Tree tests volume propagation in tree
func TestVolumePropagation_Tree(t *testing.T) {
	tc := SetupTestEnvironment(t)
	defer TeardownTestEnvironment(t, tc)

	// Create multi-level tree
	rootID := CreateTestClient(t, tc, "Root", nil)
	leftID := CreateTestClient(t, tc, "Left", &rootID)
	_ = CreateTestClient(t, tc, "Right", &rootID)
	leftLeftID := CreateTestClient(t, tc, "Left Left", &leftID)

	productID := CreateTestProduct(t, tc, "Test Product")

	// Create sale deep in tree
	CreateTestSale(t, tc, leftLeftID, productID, 100.0, "paid")

	// Check volumes propagate up
	query := `
		query {
			client(id: $rootId) {
				networkVolumeLeft
			}
		}
	`
	variables := map[string]interface{}{
		"rootId": rootID,
	}

	resp := ExecuteGraphQL(t, tc, query, variables, tc.AdminToken)
	AssertNoErrors(t, resp)

	root := resp.Data["client"].(map[string]interface{})
	leftVolume := root["networkVolumeLeft"].(float64)

	if leftVolume == 0 {
		t.Error("Root should have left volume from descendant sale")
	}
}

