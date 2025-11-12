package service

import (
	"context"
	"errors"

	"bureau/internal/models"
	"bureau/internal/store"

	"go.uber.org/zap"
)

type CaisseService struct {
	caisseRepo *store.CaisseRepository
	logger     *zap.Logger
}

func NewCaisseService(caisseRepo *store.CaisseRepository, logger *zap.Logger) *CaisseService {
	return &CaisseService{
		caisseRepo: caisseRepo,
		logger:     logger,
	}
}

// GetCaisse gets the caisse (creates it if it doesn't exist)
func (s *CaisseService) GetCaisse(ctx context.Context) (*models.Caisse, error) {
	return s.caisseRepo.GetOrCreate(ctx)
}

// AddTransaction adds a transaction and updates the caisse balance
func (s *CaisseService) AddTransaction(ctx context.Context, transaction *models.CaisseTransaction) (*models.CaisseTransaction, error) {
	// Validate transaction type
	if transaction.Type != "entree" && transaction.Type != "sortie" {
		return nil, errors.New("le type de transaction doit être 'entree' ou 'sortie'")
	}

	if transaction.Amount <= 0 {
		return nil, errors.New("le montant doit être supérieur à 0")
	}

	// Get current caisse
	caisse, err := s.caisseRepo.GetOrCreate(ctx)
	if err != nil {
		return nil, err
	}

	// Add transaction
	createdTransaction, err := s.caisseRepo.AddTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// Update caisse balance
	var newBalance, newTotalEntrees, newTotalSorties float64
	if transaction.Type == "entree" {
		newBalance = caisse.Balance + transaction.Amount
		newTotalEntrees = caisse.TotalEntrees + transaction.Amount
		newTotalSorties = caisse.TotalSorties
	} else {
		newBalance = caisse.Balance - transaction.Amount
		newTotalEntrees = caisse.TotalEntrees
		newTotalSorties = caisse.TotalSorties + transaction.Amount
	}

	err = s.caisseRepo.UpdateBalance(ctx, newBalance, newTotalEntrees, newTotalSorties)
	if err != nil {
		s.logger.Error("Failed to update caisse balance", zap.Error(err))
		// Continue anyway, transaction is created
	}

	return createdTransaction, nil
}

// GetTransactions gets all transactions
func (s *CaisseService) GetTransactions(ctx context.Context, filter *models.FilterInput, paging *models.PagingInput) ([]*models.CaisseTransaction, error) {
	return s.caisseRepo.GetTransactions(ctx, filter, paging)
}

// GetTransactionByID gets a transaction by ID
func (s *CaisseService) GetTransactionByID(ctx context.Context, id string) (*models.CaisseTransaction, error) {
	return s.caisseRepo.GetTransactionByID(ctx, id)
}

// UpdateBalance manually updates the caisse balance (admin only)
func (s *CaisseService) UpdateBalance(ctx context.Context, balance float64) (*models.Caisse, error) {
	// Get current caisse to preserve totalEntrees and totalSorties
	caisse, err := s.caisseRepo.GetOrCreate(ctx)
	if err != nil {
		return nil, err
	}

	err = s.caisseRepo.UpdateBalance(ctx, balance, caisse.TotalEntrees, caisse.TotalSorties)
	if err != nil {
		return nil, err
	}

	caisse.Balance = balance
	caisse.UpdatedAt = caisse.UpdatedAt // Will be updated by repository
	return caisse, nil
}

