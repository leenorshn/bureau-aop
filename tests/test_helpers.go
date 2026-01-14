package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bureau/graph"
	"bureau/internal/auth"
	"bureau/internal/config"
	"bureau/internal/models"
	"bureau/internal/service"
	"bureau/internal/store"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap/zaptest"
)

// TestConfig holds test configuration
type TestConfig struct {
	MongoDB      *mongo.Database
	MongoClient  *mongo.Client
	Resolver     *graph.Resolver
	Server       *httptest.Server
	AdminToken   string
	ClientToken  string
	TestAdminID  string
	TestClientID string
}

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []interface{}          `json:"errors,omitempty"`
}

// SetupTestEnvironment initializes test environment with MongoDB and GraphQL server
func SetupTestEnvironment(t *testing.T) *TestConfig {
	// Load test configuration from env.test file
	if err := godotenv.Load("env.test"); err != nil {
		// If env.test doesn't exist, try to load .env or use defaults
		_ = godotenv.Load()
		t.Logf("Note: Could not load env.test, using default/test environment variables")
	}
	
	// Override with test database name
	os.Setenv("MONGO_DB_NAME", "mlm_test_db")
	cfg := config.Load()

	logger := zaptest.NewLogger(t)

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(cfg.MongoURI).
		SetMaxPoolSize(10).
		SetMinPoolSize(1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping to verify connection
	if err := mongoClient.Ping(ctx, nil); err != nil {
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := mongoClient.Database(cfg.MongoDBName)

	// Clean test database
	if err := db.Drop(ctx); err != nil {
		t.Logf("Warning: Failed to drop test database: %v", err)
	}

	// Initialize repositories
	productRepo := store.NewProductRepository(db)
	clientRepo := store.NewClientRepository(db)
	saleRepo := store.NewSaleRepository(db, logger)
	paymentRepo := store.NewPaymentRepository(db)
	commissionRepo := store.NewCommissionRepository(db)
	adminRepo := store.NewAdminRepository(db)
	caisseRepo := store.NewCaisseRepository(db)
	binaryCappingRepo := store.NewBinaryCappingRepository(db)

	// Initialize JWT service
	jwtService := auth.NewJWTService(cfg, logger)

	// Initialize services
	productService := service.NewProductService(productRepo, logger)
	clientService := service.NewClientService(clientRepo, saleRepo, commissionRepo, logger, cfg.BinaryThreshold, cfg.BinaryCommissionRate, cfg.DefaultProductPrice)
	saleService := service.NewSaleService(saleRepo, logger)
	paymentService := service.NewPaymentService(paymentRepo, logger)
	commissionService := service.NewCommissionService(commissionRepo, clientRepo, logger, cfg.BinaryCommissionRate, cfg.BinaryThreshold)
	adminService := service.NewAdminService(adminRepo, clientRepo, productRepo, saleRepo, commissionRepo, logger)
	authService := service.NewAuthService(adminRepo, jwtService, logger)
	caisseService := service.NewCaisseService(caisseRepo, logger)

	// Initialize Transaction Helper
	txHelper := store.NewTransactionHelper(mongoClient)

	// Initialize Binary Commission Service
	binaryConfig := models.BinaryConfig{
		CycleValue:         cfg.BinaryCycleValue,
		DailyCycleLimit:    cfg.BinaryDailyCycleLimit,
		WeeklyCycleLimit:   cfg.BinaryWeeklyCycleLimit,
		MinVolumePerLeg:    cfg.BinaryMinVolumePerLeg,
		RequireDirectLeft:  true,
		RequireDirectRight: true,
	}
	binaryCommissionService := service.NewBinaryCommissionService(
		clientRepo,
		commissionRepo,
		saleRepo,
		binaryCappingRepo,
		logger,
		binaryConfig,
		txHelper,
	)

	// Initialize GraphQL resolver
	resolver := graph.NewResolver(
		productService,
		clientService,
		saleService,
		paymentService,
		commissionService,
		authService,
		adminService,
		caisseService,
		binaryCommissionService,
	)

	// Create GraphQL handler
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})

	// Create authentication middleware for tests (similar to main server)
	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token := authHeader[7:]
				// Validate token and add to context if valid
				claims, err := jwtService.ValidateAccessToken(token)
				if err == nil && claims != nil {
					ctx := context.WithValue(r.Context(), "user", claims)
					r = r.WithContext(ctx)
				}
			}
			next.ServeHTTP(w, r)
		})
	}

	// Create test server with auth middleware
	testServer := httptest.NewServer(authMiddleware(srv))

	// Create test admin
	hashedPassword, _ := auth.HashPassword("Test123@admin")
	testAdmin := &models.Admin{
		Name:         "Test Admin",
		Email:        "test-admin@test.com",
		PasswordHash:  hashedPassword,
		Role:         "admin",
	}
	createdAdmin, err := adminRepo.Create(ctx, testAdmin)
	if err != nil {
		t.Fatalf("Failed to create test admin: %v", err)
	}

	// Get admin token
	authPayload, err := authService.AdminLogin(ctx, "test-admin@test.com", "Test123@admin")
	if err != nil {
		t.Fatalf("Failed to login test admin: %v", err)
	}

	return &TestConfig{
		MongoDB:     db,
		MongoClient: mongoClient,
		Resolver:   resolver,
		Server:     testServer,
		AdminToken: authPayload.AccessToken,
		TestAdminID: createdAdmin.ID.Hex(),
	}
}

// TeardownTestEnvironment cleans up test environment
func TeardownTestEnvironment(t *testing.T, tc *TestConfig) {
	if tc.Server != nil {
		tc.Server.Close()
	}
	if tc.MongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tc.MongoDB.Drop(ctx); err != nil {
			t.Logf("Warning: Failed to drop test database: %v", err)
		}
	}
	if tc.MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tc.MongoClient.Disconnect(ctx); err != nil {
			t.Logf("Warning: Failed to disconnect MongoDB client: %v", err)
		}
	}
}

// ExecuteGraphQL executes a GraphQL request
func ExecuteGraphQL(t *testing.T, tc *TestConfig, query string, variables map[string]interface{}, token string) *GraphQLResponse {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", tc.Server.URL+"/query", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	var gqlResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return &gqlResp
}

// CreateTestProduct creates a test product
func CreateTestProduct(t *testing.T, tc *TestConfig, name string) string {
	query := fmt.Sprintf(`
		mutation {
			productCreate(input: {
				name: "%s"
				description: "Test product description"
				price: 100.0
				stock: 50
				points: 10.0
				imageUrl: "https://example.com/image.jpg"
			}) {
				id
			}
		}
	`, name)

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	if resp.Errors != nil && len(resp.Errors) > 0 {
		t.Fatalf("Failed to create product: %v", resp.Errors)
	}

	data := resp.Data["productCreate"].(map[string]interface{})
	return data["id"].(string)
}

// CreateTestClient creates a test client
func CreateTestClient(t *testing.T, tc *TestConfig, name string, sponsorID *string) string {
	var query string
	if sponsorID != nil {
		query = fmt.Sprintf(`
			mutation {
				clientCreate(input: {
					name: "%s"
					password: "Test123@client"
					sponsorId: "%s"
				}) {
					id
					clientId
				}
			}
		`, name, *sponsorID)
	} else {
		query = fmt.Sprintf(`
			mutation {
				clientCreate(input: {
					name: "%s"
					password: "Test123@client"
				}) {
					id
					clientId
				}
			}
		`, name)
	}

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	if resp.Errors != nil && len(resp.Errors) > 0 {
		t.Fatalf("Failed to create client: %v", resp.Errors)
	}

	data := resp.Data["clientCreate"].(map[string]interface{})
	return data["id"].(string)
}

// CreateTestSale creates a test sale
func CreateTestSale(t *testing.T, tc *TestConfig, clientID, productID string, amount float64, status string) string {
	query := fmt.Sprintf(`
		mutation {
			saleCreate(input: {
				clientId: "%s"
				productId: "%s"
				quantity: 1
				amount: %.2f
				status: "%s"
			}) {
				id
			}
		}
	`, clientID, productID, amount, status)

	resp := ExecuteGraphQL(t, tc, query, nil, tc.AdminToken)
	if resp.Errors != nil && len(resp.Errors) > 0 {
		t.Fatalf("Failed to create sale: %v", resp.Errors)
	}

	data := resp.Data["saleCreate"].(map[string]interface{})
	return data["id"].(string)
}

// AssertNoErrors checks that there are no GraphQL errors
func AssertNoErrors(t *testing.T, resp *GraphQLResponse) {
	if resp.Errors != nil && len(resp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %v", resp.Errors)
	}
}

// AssertHasErrors checks that there are GraphQL errors
func AssertHasErrors(t *testing.T, resp *GraphQLResponse) {
	if resp.Errors == nil || len(resp.Errors) == 0 {
		t.Error("Expected GraphQL errors but got none")
	}
}

