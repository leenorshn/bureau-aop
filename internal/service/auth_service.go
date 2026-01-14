package service

import (
	"context"
	"errors"
	"fmt"

	"bureau/internal/auth"
	"bureau/internal/models"
	"bureau/internal/store"

	"go.uber.org/zap"
)

type AuthService struct {
	adminRepo  *store.AdminRepository
	jwtService *auth.JWTService
	logger     *zap.Logger
}

func NewAuthService(adminRepo *store.AdminRepository, jwtService *auth.JWTService, logger *zap.Logger) *AuthService {
	return &AuthService{
		adminRepo:  adminRepo,
		jwtService: jwtService,
		logger:     logger,
	}
}

// AdminLogin authenticates an admin and returns JWT tokens
func (s *AuthService) AdminLogin(ctx context.Context, email, password string) (*models.AuthPayload, error) {
	// Get admin by email
	admin, err := s.adminRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !auth.CheckPasswordHash(password, admin.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(admin)
	if err != nil {
		s.logger.Error("Failed to generate access token", zap.Error(err))
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(admin)
	if err != nil {
		s.logger.Error("Failed to generate refresh token", zap.Error(err))
		return nil, errors.New("failed to generate refresh token")
	}

	return &models.AuthPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Admin:        admin,
	}, nil
}

// RefreshToken validates a refresh token and returns new tokens
func (s *AuthService) RefreshToken(ctx context.Context, tokenString string) (*models.AuthPayload, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get admin by ID
	admin, err := s.adminRepo.GetByID(ctx, claims.AdminID)
	if err != nil {
		return nil, errors.New("admin not found")
	}

	// Generate new tokens
	accessToken, err := s.jwtService.GenerateAccessToken(admin)
	if err != nil {
		s.logger.Error("Failed to generate access token", zap.Error(err))
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(admin)
	if err != nil {
		s.logger.Error("Failed to generate refresh token", zap.Error(err))
		return nil, errors.New("failed to generate refresh token")
	}

	return &models.AuthPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Admin:        admin,
	}, nil
}

// ValidateToken validates an access token and returns the admin
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*models.Admin, error) {
	// Validate access token
	claims, err := s.jwtService.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid access token")
	}

	// Get admin by ID
	admin, err := s.adminRepo.GetByID(ctx, claims.AdminID)
	if err != nil {
		return nil, errors.New("admin not found")
	}

	return admin, nil
}

// GetJWTService returns the JWT service for external use
func (s *AuthService) GetJWTService() *auth.JWTService {
	return s.jwtService
}

// UpdateAdminPassword updates an admin's password
func (s *AuthService) UpdateAdminPassword(ctx context.Context, id string, newPassword string) error {
	// Validate password strength
	if err := auth.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Hash the new password
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.adminRepo.UpdatePassword(ctx, id, hashedPassword)
}

// UpdateAdminPasswordByEmail updates an admin's password by email
func (s *AuthService) UpdateAdminPasswordByEmail(ctx context.Context, email string, newPassword string) error {
	// Get admin by email
	admin, err := s.adminRepo.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("admin introuvable")
	}

	// Validate password strength
	if err := auth.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Hash the new password
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.adminRepo.UpdatePassword(ctx, admin.ID.Hex(), hashedPassword)
}
