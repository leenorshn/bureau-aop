package service

import (
	"context"

	"bureau/internal/models"
	"bureau/internal/store"

	"go.uber.org/zap"
)

type ProductService struct {
	productRepo *store.ProductRepository
	logger      *zap.Logger
}

func NewProductService(productRepo *store.ProductRepository, logger *zap.Logger) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		logger:      logger,
	}
}

func (s *ProductService) GetAll(ctx context.Context, filter *models.FilterInput, paging *models.PagingInput) ([]*models.Product, error) {
	return s.productRepo.GetAll(ctx, filter, paging)
}

func (s *ProductService) GetByID(ctx context.Context, id string) (*models.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *ProductService) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	return s.productRepo.Create(ctx, product)
}

func (s *ProductService) Update(ctx context.Context, id string, product *models.Product) (*models.Product, error) {
	return s.productRepo.Update(ctx, id, product)
}

func (s *ProductService) Delete(ctx context.Context, id string) (bool, error) {
	err := s.productRepo.Delete(ctx, id)
	return err == nil, err
}

