package service

import (
	"context"

	"bureau/internal/models"
	"bureau/internal/store"

	"go.uber.org/zap"
)

type PaymentService struct {
	paymentRepo *store.PaymentRepository
	logger      *zap.Logger
}

func NewPaymentService(paymentRepo *store.PaymentRepository, logger *zap.Logger) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

func (s *PaymentService) GetAll(ctx context.Context, filter *models.FilterInput, paging *models.PagingInput) ([]*models.Payment, error) {
	return s.paymentRepo.GetAll(ctx, filter, paging)
}

func (s *PaymentService) GetByID(ctx context.Context, id string) (*models.Payment, error) {
	return s.paymentRepo.GetByID(ctx, id)
}

func (s *PaymentService) Create(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	return s.paymentRepo.Create(ctx, payment)
}

func (s *PaymentService) Update(ctx context.Context, id string, payment *models.Payment) (*models.Payment, error) {
	return s.paymentRepo.Update(ctx, id, payment)
}

func (s *PaymentService) Delete(ctx context.Context, id string) (bool, error) {
	err := s.paymentRepo.Delete(ctx, id)
	return err == nil, err
}

func (s *PaymentService) GetByClientID(ctx context.Context, clientID string) ([]*models.Payment, error) {
	return s.paymentRepo.GetByClientID(ctx, clientID)
}

func (s *PaymentService) GetTotalPayments(ctx context.Context, filter *models.FilterInput) (float64, error) {
	return s.paymentRepo.GetTotalPayments(ctx, filter)
}




