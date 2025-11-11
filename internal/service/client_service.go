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
// If requestedPosition is provided ("left" or "right"), it will try to place the client at that position.
// If the requested position is not available, it returns an error.
// If requestedPosition is nil, it uses the first available position (left then right).
// If both positions are taken, it returns an error asking user to choose another sponsor.
func (s *ClientService) CreateWithBinaryPlacement(ctx context.Context, client *models.Client, sponsorID *primitive.ObjectID, requestedPosition *string) (*models.Client, error) {
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
		client.Points = 0

		createdClient, err := s.clientRepo.Create(ctx, client)
		if err != nil {
			return nil, err
		}
		return createdClient, nil
	}

	// Find the sponsor
	sponsor, err := s.clientRepo.GetByID(ctx, sponsorID.Hex())
	if err != nil {
		return nil, errors.New("sponsor introuvable")
	}

	// Check if sponsor has available positions
	hasLeftChild := sponsor.LeftChildID != nil
	hasRightChild := sponsor.RightChildID != nil

	// If both positions are taken, return error
	if hasLeftChild && hasRightChild {
		return nil, errors.New("ce sponsor a déjà 2 enfants (les positions gauche et droite sont prises). Veuillez choisir un autre sponsor qui a une position disponible")
	}

	// Determine position
	var position string
	if requestedPosition != nil {
		// User specified a position
		requestedPos := *requestedPosition
		if requestedPos != "left" && requestedPos != "right" {
			return nil, errors.New("la position doit être 'left' ou 'right'")
		}

		// Check if requested position is available
		if requestedPos == "left" && hasLeftChild {
			return nil, errors.New("la position gauche est déjà prise sur ce sponsor. Veuillez choisir 'right' ou sélectionner un autre sponsor")
		}
		if requestedPos == "right" && hasRightChild {
			return nil, errors.New("la position droite est déjà prise sur ce sponsor. Veuillez choisir 'left' ou sélectionner un autre sponsor")
		}

		position = requestedPos
	} else {
		// No position specified, use first available
		if !hasLeftChild {
			position = "left"
		} else {
			position = "right"
		}
	}

	// Set client properties
	client.SponsorID = sponsorID
	client.Position = &position
	client.LeftChildID = nil
	client.RightChildID = nil
	client.NetworkVolumeLeft = 0
	client.NetworkVolumeRight = 0
	client.BinaryPairs = 0
	client.TotalEarnings = 0
	client.WalletBalance = 0
	client.Points = 0

	// Create the client
	createdClient, err := s.clientRepo.Create(ctx, client)
	if err != nil {
		return nil, err
	}

	// Update sponsor's binary tree
	err = s.updateSponsorBinaryTree(ctx, *sponsorID, createdClient.ID, position)
	if err != nil {
		s.logger.Error("Failed to update sponsor binary tree", zap.Error(err))
		// Continue anyway, the client is created
	}

	// Update network volumes and check for binary commissions
	err = s.updateNetworkVolumesAndCommissions(ctx, *sponsorID, s.defaultProductPrice, position)
	if err != nil {
		s.logger.Error("Failed to update network volumes and commissions", zap.Error(err))
	}

	return createdClient, nil
}

// updateSponsorBinaryTree updates the sponsor's binary tree with the new client
func (s *ClientService) updateSponsorBinaryTree(ctx context.Context, sponsorID primitive.ObjectID, clientID primitive.ObjectID, position string) error {
	// Verify sponsor exists and check binary tree constraint
	sponsor, err := s.clientRepo.GetByID(ctx, sponsorID.Hex())
	if err != nil {
		return fmt.Errorf("sponsor introuvable: %w", err)
	}

	// Enforce binary tree constraint: a sponsor can only have 2 children
	if position == "left" {
		if sponsor.LeftChildID != nil {
			return errors.New("le sponsor a déjà un enfant à gauche - contrainte d'arbre binaire violée")
		}
		return s.clientRepo.UpdateBinaryFields(ctx, sponsorID.Hex(), &clientID, nil, &position)
	} else {
		if sponsor.RightChildID != nil {
			return errors.New("le sponsor a déjà un enfant à droite - contrainte d'arbre binaire violée")
		}
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

// AddPoints adds points to a client
func (s *ClientService) AddPoints(ctx context.Context, clientID string, pointsToAdd float64) error {
	return s.clientRepo.AddPoints(ctx, clientID, pointsToAdd)
}
