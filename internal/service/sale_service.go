package service

import (
	"context"

	"bureau/internal/models"
	"bureau/internal/store"

	"go.uber.org/zap"
)

type SaleService struct {
	saleRepo *store.SaleRepository
	logger   *zap.Logger
}

func NewSaleService(saleRepo *store.SaleRepository, logger *zap.Logger) *SaleService {
	return &SaleService{
		saleRepo: saleRepo,
		logger:   logger,
	}
}

func (s *SaleService) GetAll(ctx context.Context, filter *models.FilterInput, paging *models.PagingInput) ([]*models.Sale, error) {
	return s.saleRepo.GetAll(ctx, filter, paging)
}

func (s *SaleService) GetByID(ctx context.Context, id string) (*models.Sale, error) {
	return s.saleRepo.GetByID(ctx, id)
}

func (s *SaleService) Create(ctx context.Context, sale *models.Sale) (*models.Sale, error) {
	return s.saleRepo.Create(ctx, sale)
}

func (s *SaleService) Update(ctx context.Context, id string, sale *models.Sale) (*models.Sale, error) {
	return s.saleRepo.Update(ctx, id, sale)
}

func (s *SaleService) Delete(ctx context.Context, id string) (bool, error) {
	err := s.saleRepo.Delete(ctx, id)
	return err == nil, err
}

func (s *SaleService) GetByClientID(ctx context.Context, clientID string) ([]*models.Sale, error) {
	return s.saleRepo.GetByClientID(ctx, clientID)
}

func (s *SaleService) GetBySponsorID(ctx context.Context, sponsorID string) ([]*models.Sale, error) {
	return s.saleRepo.GetBySponsorID(ctx, sponsorID)
}

func (s *SaleService) GetTotalSales(ctx context.Context, filter *models.FilterInput) (float64, error) {
	return s.saleRepo.GetTotalSales(ctx, filter)
}

