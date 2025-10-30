package service

import (
	"context"
	"errors"
	"fmt"

	"bureau/internal/models"
	"bureau/internal/store"
	"bureau/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type ClientService struct {
	clientRepo           *store.ClientRepository
	saleRepo             *store.SaleRepository
	commissionRepo       *store.CommissionRepository
	logger               *zap.Logger
	binaryThreshold      float64
	binaryCommissionRate float64
	defaultProductPrice  float64
}

func NewClientService(
	clientRepo *store.ClientRepository,
	saleRepo *store.SaleRepository,
	commissionRepo *store.CommissionRepository,
	logger *zap.Logger,
	binaryThreshold, binaryCommissionRate, defaultProductPrice float64,
) *ClientService {
	return &ClientService{
		clientRepo:           clientRepo,
		saleRepo:             saleRepo,
		commissionRepo:       commissionRepo,
		logger:               logger,
		binaryThreshold:      binaryThreshold,
		binaryCommissionRate: binaryCommissionRate,
		defaultProductPrice:  defaultProductPrice,
	}
}

func (s *ClientService) GetAll(ctx context.Context, filter *models.FilterInput, paging *models.PagingInput) ([]*models.Client, error) {
	return s.clientRepo.GetAll(ctx, filter, paging)
}

func (s *ClientService) GetByID(ctx context.Context, id string) (*models.Client, error) {
	return s.clientRepo.GetByID(ctx, id)
}

func (s *ClientService) Update(ctx context.Context, id string, client *models.Client) (*models.Client, error) {
	return s.clientRepo.Update(ctx, id, client)
}

func (s *ClientService) Delete(ctx context.Context, id string) (bool, error) {
	err := s.clientRepo.Delete(ctx, id)
	return err == nil, err
}

// CreateWithBinaryPlacement creates a new client and places them in the binary tree
func (s *ClientService) CreateWithBinaryPlacement(ctx context.Context, client *models.Client, sponsorID *primitive.ObjectID) (*models.Client, error) {
	// Generate unique client ID
	clientID, err := s.generateUniqueClientID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate client ID: %w", err)
	}
	client.ClientID = clientID

	// Hash the password
	hashedPassword, err := s.HashPassword(client.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	client.PasswordHash = hashedPassword

	// If no sponsor provided, this is the root client
	if sponsorID == nil {
		client.SponsorID = nil
		client.Position = nil
		client.LeftChildID = nil
		client.RightChildID = nil
		client.NetworkVolumeLeft = 0
		client.NetworkVolumeRight = 0
		client.BinaryPairs = 0
		client.TotalEarnings = 0
		client.WalletBalance = 0

		createdClient, err := s.clientRepo.Create(ctx, client)
		if err != nil {
			return nil, err
		}

		// Note: Auto-generated sales removed to avoid frontend issues
		// Sales should be created explicitly through the SaleCreate mutation

		return createdClient, nil
	}

	// Find the sponsor
	_, err = s.clientRepo.GetByID(ctx, sponsorID.Hex())
	if err != nil {
		return nil, errors.New("sponsor not found")
	}

	// Find placement position in binary tree
	placement, err := s.findBinaryPlacement(ctx, sponsorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find binary placement: %w", err)
	}

	// Set client properties
	client.SponsorID = sponsorID
	client.Position = &placement.Position
	client.LeftChildID = nil
	client.RightChildID = nil
	client.NetworkVolumeLeft = 0
	client.NetworkVolumeRight = 0
	client.BinaryPairs = 0
	client.TotalEarnings = 0
	client.WalletBalance = 0

	// Create the client
	createdClient, err := s.clientRepo.Create(ctx, client)
	if err != nil {
		return nil, err
	}

	// Update sponsor's binary tree
	err = s.updateSponsorBinaryTree(ctx, placement.SponsorID, createdClient.ID, placement.Position)
	if err != nil {
		s.logger.Error("Failed to update sponsor binary tree", zap.Error(err))
		// Continue anyway, the client is created
	}

	// Note: Auto-generated sales removed to avoid frontend issues
	// Sales should be created explicitly through the SaleCreate mutation

	// Update network volumes and check for binary commissions
	err = s.updateNetworkVolumesAndCommissions(ctx, placement.SponsorID, s.defaultProductPrice, placement.Position)
	if err != nil {
		s.logger.Error("Failed to update network volumes and commissions", zap.Error(err))
	}

	return createdClient, nil
}

type BinaryPlacement struct {
	SponsorID primitive.ObjectID
	Position  string // "left" or "right"
}

// findBinaryPlacement finds the appropriate position for a new client in the binary tree
func (s *ClientService) findBinaryPlacement(ctx context.Context, sponsorID *primitive.ObjectID) (*BinaryPlacement, error) {
	// Start from the sponsor and traverse down to find the first available position
	queue := []primitive.ObjectID{*sponsorID}

	for len(queue) > 0 {
		currentID := queue[0]
		queue = queue[1:]

		// Get current client
		client, err := s.clientRepo.GetByID(ctx, currentID.Hex())
		if err != nil {
			continue
		}

		// Check if left position is available
		if client.LeftChildID == nil {
			return &BinaryPlacement{
				SponsorID: currentID,
				Position:  "left",
			}, nil
		}

		// Check if right position is available
		if client.RightChildID == nil {
			return &BinaryPlacement{
				SponsorID: currentID,
				Position:  "right",
			}, nil
		}

		// Both positions are filled, add children to queue for BFS traversal
		queue = append(queue, *client.LeftChildID, *client.RightChildID)
	}

	return nil, errors.New("no available position found in binary tree")
}

// updateSponsorBinaryTree updates the sponsor's binary tree with the new client
func (s *ClientService) updateSponsorBinaryTree(ctx context.Context, sponsorID primitive.ObjectID, clientID primitive.ObjectID, position string) error {
	if position == "left" {
		return s.clientRepo.UpdateBinaryFields(ctx, sponsorID.Hex(), &clientID, nil, &position)
	} else {
		return s.clientRepo.UpdateBinaryFields(ctx, sponsorID.Hex(), nil, &clientID, &position)
	}
}

// updateNetworkVolumesAndCommissions updates network volumes and calculates binary commissions
func (s *ClientService) updateNetworkVolumesAndCommissions(ctx context.Context, sponsorID primitive.ObjectID, amount float64, side string) error {
	// Traverse up the sponsor chain and update volumes
	currentID := sponsorID

	for currentID != primitive.NilObjectID {
		client, err := s.clientRepo.GetByID(ctx, currentID.Hex())
		if err != nil {
			break
		}

		// Update network volumes
		if side == "left" {
			client.NetworkVolumeLeft += amount
		} else {
			client.NetworkVolumeRight += amount
		}

		// Update in database
		err = s.clientRepo.UpdateNetworkVolumes(ctx, currentID.Hex(), client.NetworkVolumeLeft, client.NetworkVolumeRight)
		if err != nil {
			s.logger.Error("Failed to update network volumes", zap.Error(err))
		}

		// Check for binary commission eligibility
		if client.NetworkVolumeLeft >= s.binaryThreshold && client.NetworkVolumeRight >= s.binaryThreshold {
			err = s.calculateBinaryCommission(ctx, client)
			if err != nil {
				s.logger.Error("Failed to calculate binary commission", zap.Error(err))
			}
		}

		// Move to parent sponsor
		if client.SponsorID != nil {
			currentID = *client.SponsorID
		} else {
			break
		}
	}

	return nil
}

// calculateBinaryCommission calculates and creates binary commission
func (s *ClientService) calculateBinaryCommission(ctx context.Context, client *models.Client) error {
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

	_, err := s.commissionRepo.Create(ctx, commission)
	if err != nil {
		return err
	}

	// Update client earnings and wallet
	newTotalEarnings := client.TotalEarnings + commissionAmount
	newWalletBalance := client.WalletBalance + commissionAmount
	newBinaryPairs := client.BinaryPairs + 1

	err = s.clientRepo.UpdateEarnings(ctx, client.ID.Hex(), newTotalEarnings, newWalletBalance)
	if err != nil {
		return err
	}

	err = s.clientRepo.UpdateBinaryPairs(ctx, client.ID.Hex(), newBinaryPairs)
	if err != nil {
		return err
	}

	// Reduce network volumes (consume the matched volume)
	consumedVolume := commissionAmount / s.binaryCommissionRate
	client.NetworkVolumeLeft -= consumedVolume
	client.NetworkVolumeRight -= consumedVolume

	err = s.clientRepo.UpdateNetworkVolumes(ctx, client.ID.Hex(), client.NetworkVolumeLeft, client.NetworkVolumeRight)
	if err != nil {
		return err
	}

	s.logger.Info("Binary commission calculated",
		zap.String("clientID", client.ID.Hex()),
		zap.Float64("amount", commissionAmount),
		zap.Int("binaryPairs", newBinaryPairs),
	)

	return nil
}

// generateUniqueClientID generates a unique 8-digit client ID
func (s *ClientService) generateUniqueClientID(ctx context.Context) (string, error) {
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		clientID, err := utils.GenerateClientID()
		if err != nil {
			return "", err
		}

		// Check if this ID already exists
		_, err = s.clientRepo.GetByClientID(ctx, clientID)
		if err != nil {
			// ID doesn't exist, we can use it
			return clientID, nil
		}
	}

	return "", errors.New("failed to generate unique client ID after multiple attempts")
}

// AuthenticateClient authenticates a client using clientId and password
func (s *ClientService) AuthenticateClient(ctx context.Context, clientID, password string) (*models.Client, error) {
	client, err := s.clientRepo.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(client.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return client, nil
}

// HashPassword hashes a password using bcrypt
func (s *ClientService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
