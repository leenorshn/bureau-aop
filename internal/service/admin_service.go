package service

import (
	"context"
	"time"

	"bureau/internal/models"
	"bureau/internal/store"

	"go.uber.org/zap"
)

type AdminService struct {
	adminRepo      *store.AdminRepository
	clientRepo     *store.ClientRepository
	productRepo    *store.ProductRepository
	saleRepo       *store.SaleRepository
	commissionRepo *store.CommissionRepository
	logger         *zap.Logger
}

func NewAdminService(
	adminRepo *store.AdminRepository,
	clientRepo *store.ClientRepository,
	productRepo *store.ProductRepository,
	saleRepo *store.SaleRepository,
	commissionRepo *store.CommissionRepository,
	logger *zap.Logger,
) *AdminService {
	return &AdminService{
		adminRepo:      adminRepo,
		clientRepo:     clientRepo,
		productRepo:    productRepo,
		saleRepo:       saleRepo,
		commissionRepo: commissionRepo,
		logger:         logger,
	}
}

func (s *AdminService) GetDashboardStats(ctx context.Context, rangeArg *string) (*models.DashboardStats, error) {
	// Set default range if not provided
	if rangeArg == nil {
		rangeStr := "30d"
		rangeArg = &rangeStr
	}

	// Calculate date range
	var dateFrom *time.Time
	now := time.Now()

	switch *rangeArg {
	case "7d":
		from := now.AddDate(0, 0, -7)
		dateFrom = &from
	case "30d":
		from := now.AddDate(0, 0, -30)
		dateFrom = &from
	case "90d":
		from := now.AddDate(0, 0, -90)
		dateFrom = &from
	case "1y":
		from := now.AddDate(-1, 0, 0)
		dateFrom = &from
	}

	// Create filter for date range
	filter := &models.FilterInput{
		DateFrom: dateFrom,
		DateTo:   &now,
	}

	// Get counts
	totalClients, err := s.clientRepo.Count(ctx, nil)
	if err != nil {
		return nil, err
	}

	totalProducts, err := s.productRepo.Count(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Get totals
	totalSales, err := s.saleRepo.GetTotalSales(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalCommissions, err := s.commissionRepo.GetTotalCommissions(ctx, filter)
	if err != nil {
		return nil, err
	}

	// For active clients, we'll consider clients with sales in the range
	activeClients, err := s.clientRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &models.DashboardStats{
		TotalClients:     int(totalClients),
		TotalSales:       totalSales,
		TotalCommissions: totalCommissions,
		TotalProducts:    int(totalProducts),
		ActiveClients:    int(activeClients),
	}, nil
}





















