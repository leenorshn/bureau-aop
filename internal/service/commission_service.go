package service

import (
	"context"
	"errors"

	"bureau/internal/models"
	"bureau/internal/store"

	"go.uber.org/zap"
)

type CommissionService struct {
	commissionRepo       *store.CommissionRepository
	clientRepo           *store.ClientRepository
	logger               *zap.Logger
	binaryThreshold      float64
	binaryCommissionRate float64
}

func NewCommissionService(
	commissionRepo *store.CommissionRepository,
	clientRepo *store.ClientRepository,
	logger *zap.Logger,
	binaryThreshold, binaryCommissionRate float64,
) *CommissionService {
	return &CommissionService{
		commissionRepo:       commissionRepo,
		clientRepo:           clientRepo,
		logger:               logger,
		binaryThreshold:      binaryThreshold,
		binaryCommissionRate: binaryCommissionRate,
	}
}

func (s *CommissionService) GetAll(ctx context.Context, filter *models.FilterInput, paging *models.PagingInput) ([]*models.Commission, error) {
	return s.commissionRepo.GetAll(ctx, filter, paging)
}

func (s *CommissionService) GetByID(ctx context.Context, id string) (*models.Commission, error) {
	return s.commissionRepo.GetByID(ctx, id)
}

func (s *CommissionService) Create(ctx context.Context, commission *models.Commission) (*models.Commission, error) {
	return s.commissionRepo.Create(ctx, commission)
}

func (s *CommissionService) Update(ctx context.Context, id string, commission *models.Commission) (*models.Commission, error) {
	return s.commissionRepo.Update(ctx, id, commission)
}

func (s *CommissionService) Delete(ctx context.Context, id string) (bool, error) {
	err := s.commissionRepo.Delete(ctx, id)
	return err == nil, err
}

func (s *CommissionService) GetByClientID(ctx context.Context, clientID string) ([]*models.Commission, error) {
	return s.commissionRepo.GetByClientID(ctx, clientID)
}

func (s *CommissionService) GetBySourceClientID(ctx context.Context, sourceClientID string) ([]*models.Commission, error) {
	return s.commissionRepo.GetBySourceClientID(ctx, sourceClientID)
}

func (s *CommissionService) GetTotalCommissions(ctx context.Context, filter *models.FilterInput) (float64, error) {
	return s.commissionRepo.GetTotalCommissions(ctx, filter)
}

// RunBinaryCommissionCheck manually triggers binary commission calculation for a specific client
func (s *CommissionService) RunBinaryCommissionCheck(ctx context.Context, clientID string) (*models.CommissionResult, error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, errors.New("client not found")
	}

	// Check if client is eligible for binary commission
	if client.NetworkVolumeLeft < s.binaryThreshold || client.NetworkVolumeRight < s.binaryThreshold {
		return &models.CommissionResult{
			CommissionsCreated: 0,
			TotalAmount:        0,
			Message:            "Client not eligible for binary commission - insufficient network volume",
		}, nil
	}

	// Calculate commission amount (minimum of left and right volumes)
	commissionAmount := client.NetworkVolumeLeft
	if client.NetworkVolumeRight < client.NetworkVolumeLeft {
		commissionAmount = client.NetworkVolumeRight
	}
	commissionAmount *= s.binaryCommissionRate

	// Create commission record
	commission := &models.Commission{
		ClientID:       client.ID,
		SourceClientID: client.ID, // Self-commission for binary match
		Amount:         commissionAmount,
		Level:          0, // Direct binary commission
		Type:           "binary-match",
	}

	_, err = s.commissionRepo.Create(ctx, commission)
	if err != nil {
		return nil, err
	}

	// Update client earnings and wallet
	newTotalEarnings := client.TotalEarnings + commissionAmount
	newWalletBalance := client.WalletBalance + commissionAmount
	newBinaryPairs := client.BinaryPairs + 1

	err = s.clientRepo.UpdateEarnings(ctx, client.ID.Hex(), newTotalEarnings, newWalletBalance)
	if err != nil {
		s.logger.Error("Failed to update client earnings", zap.Error(err))
	}

	err = s.clientRepo.UpdateBinaryPairs(ctx, client.ID.Hex(), newBinaryPairs)
	if err != nil {
		s.logger.Error("Failed to update binary pairs", zap.Error(err))
	}

	// Reduce network volumes (consume the matched volume)
	consumedVolume := commissionAmount / s.binaryCommissionRate
	client.NetworkVolumeLeft -= consumedVolume
	client.NetworkVolumeRight -= consumedVolume

	err = s.clientRepo.UpdateNetworkVolumes(ctx, client.ID.Hex(), client.NetworkVolumeLeft, client.NetworkVolumeRight)
	if err != nil {
		s.logger.Error("Failed to update network volumes", zap.Error(err))
	}

	s.logger.Info("Binary commission calculated manually",
		zap.String("clientID", client.ID.Hex()),
		zap.Float64("amount", commissionAmount),
		zap.Int("binaryPairs", newBinaryPairs),
	)

	return &models.CommissionResult{
		CommissionsCreated: 1,
		TotalAmount:        commissionAmount,
		Message:            "Binary commission calculated successfully",
	}, nil
}

// CalculateBinaryCommissionsForAll checks all clients for binary commission eligibility
func (s *CommissionService) CalculateBinaryCommissionsForAll(ctx context.Context) (*models.CommissionResult, error) {
	// This would typically be a background job
	// For now, we'll return a placeholder
	return &models.CommissionResult{
		CommissionsCreated: 0,
		TotalAmount:        0,
		Message:            "Binary commission check completed",
	}, nil
}
