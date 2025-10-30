package tests

import (
	"context"
	"testing"

	"bureau/internal/auth"
	"bureau/internal/config"
	"bureau/internal/models"
	"bureau/internal/service"
	"bureau/internal/store"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap/zaptest"
)

func TestJWTService(t *testing.T) {
	cfg := config.Load()
	logger := zaptest.NewLogger(t)

	jwtService := auth.NewJWTService(cfg, logger)

	// Create test admin
	admin := &models.Admin{
		ID:    primitive.NewObjectID(),
		Name:  "Test Admin",
		Email: "test@example.com",
		Role:  "admin",
	}

	// Test access token generation
	t.Run("GenerateAccessToken", func(t *testing.T) {
		token, err := jwtService.GenerateAccessToken(admin)
		if err != nil {
			t.Fatalf("Failed to generate access token: %v", err)
		}

		if token == "" {
			t.Error("Access token should not be empty")
		}
	})

	// Test refresh token generation
	t.Run("GenerateRefreshToken", func(t *testing.T) {
		token, err := jwtService.GenerateRefreshToken(admin)
		if err != nil {
			t.Fatalf("Failed to generate refresh token: %v", err)
		}

		if token == "" {
			t.Error("Refresh token should not be empty")
		}
	})

	// Test token validation
	t.Run("ValidateAccessToken", func(t *testing.T) {
		token, err := jwtService.GenerateAccessToken(admin)
		if err != nil {
			t.Fatalf("Failed to generate access token: %v", err)
		}

		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			t.Fatalf("Failed to validate access token: %v", err)
		}

		if claims.AdminID != admin.ID.Hex() {
			t.Errorf("Expected admin ID %s, got %s", admin.ID.Hex(), claims.AdminID)
		}

		if claims.Email != admin.Email {
			t.Errorf("Expected email %s, got %s", admin.Email, claims.Email)
		}

		if claims.Role != admin.Role {
			t.Errorf("Expected role %s, got %s", admin.Role, claims.Role)
		}
	})

	// Test refresh token validation
	t.Run("ValidateRefreshToken", func(t *testing.T) {
		token, err := jwtService.GenerateRefreshToken(admin)
		if err != nil {
			t.Fatalf("Failed to generate refresh token: %v", err)
		}

		claims, err := jwtService.ValidateRefreshToken(token)
		if err != nil {
			t.Fatalf("Failed to validate refresh token: %v", err)
		}

		if claims.AdminID != admin.ID.Hex() {
			t.Errorf("Expected admin ID %s, got %s", admin.ID.Hex(), claims.AdminID)
		}
	})

	// Test invalid token
	t.Run("InvalidToken", func(t *testing.T) {
		_, err := jwtService.ValidateAccessToken("invalid-token")
		if err == nil {
			t.Error("Should return error for invalid token")
		}
	})
}

func TestBcryptFunctions(t *testing.T) {
	password := "test-password-123"

	// Test password hashing
	t.Run("HashPassword", func(t *testing.T) {
		hash, err := auth.HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}

		if hash == "" {
			t.Error("Hash should not be empty")
		}

		if hash == password {
			t.Error("Hash should not be the same as password")
		}
	})

	// Test password verification
	t.Run("CheckPasswordHash", func(t *testing.T) {
		hash, err := auth.HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}

		// Test correct password
		if !auth.CheckPasswordHash(password, hash) {
			t.Error("Password verification should succeed for correct password")
		}

		// Test incorrect password
		if auth.CheckPasswordHash("wrong-password", hash) {
			t.Error("Password verification should fail for incorrect password")
		}
	})
}

func TestAuthService(t *testing.T) {
	cfg := config.Load()
	logger := zaptest.NewLogger(t)

	// Initialize MongoDB
	mongoDB, err := store.NewMongoDB(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// Initialize repositories and services
	adminRepo := store.NewAdminRepository(mongoDB.Database)
	jwtService := auth.NewJWTService(cfg, logger)
	authService := service.NewAuthService(adminRepo, jwtService, logger)

	ctx := context.Background()

	// Create test admin
	hashedPassword, err := auth.HashPassword("test-password")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	admin := &models.Admin{
		Name:         "Test Admin",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Role:         "admin",
	}

	createdAdmin, err := adminRepo.Create(ctx, admin)
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	// Test admin login
	t.Run("AdminLogin", func(t *testing.T) {
		authPayload, err := authService.AdminLogin(ctx, "test@example.com", "test-password")
		if err != nil {
			t.Fatalf("Failed to login admin: %v", err)
		}

		if authPayload.AccessToken == "" {
			t.Error("Access token should not be empty")
		}

		if authPayload.RefreshToken == "" {
			t.Error("Refresh token should not be empty")
		}

		if authPayload.Admin == nil {
			t.Error("Admin should not be nil")
		}

		if authPayload.Admin.ID.Hex() != createdAdmin.ID.Hex() {
			t.Errorf("Expected admin ID %s, got %s", createdAdmin.ID.Hex(), authPayload.Admin.ID.Hex())
		}
	})

	// Test invalid login
	t.Run("InvalidLogin", func(t *testing.T) {
		_, err := authService.AdminLogin(ctx, "test@example.com", "wrong-password")
		if err == nil {
			t.Error("Should return error for invalid password")
		}
	})

	// Test refresh token
	t.Run("RefreshToken", func(t *testing.T) {
		// First login to get refresh token
		authPayload, err := authService.AdminLogin(ctx, "test@example.com", "test-password")
		if err != nil {
			t.Fatalf("Failed to login admin: %v", err)
		}

		// Use refresh token to get new tokens
		newAuthPayload, err := authService.RefreshToken(ctx, authPayload.RefreshToken)
		if err != nil {
			t.Fatalf("Failed to refresh token: %v", err)
		}

		if newAuthPayload.AccessToken == "" {
			t.Error("New access token should not be empty")
		}

		if newAuthPayload.RefreshToken == "" {
			t.Error("New refresh token should not be empty")
		}
	})

	// Test invalid refresh token
	t.Run("InvalidRefreshToken", func(t *testing.T) {
		_, err := authService.RefreshToken(ctx, "invalid-refresh-token")
		if err == nil {
			t.Error("Should return error for invalid refresh token")
		}
	})
}






