package service

import (
	"context"
	"sync"
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

	// Use goroutines to parallelize independent queries
	var wg sync.WaitGroup
	var mu sync.Mutex
	var statsErr error

	stats := &models.DashboardStats{}

	// Parallel queries for counts and totals
	wg.Add(5)
	
	go func() {
		defer wg.Done()
		count, err := s.clientRepo.Count(ctx, nil)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			statsErr = err
			return
		}
		stats.TotalClients = int(count)
	}()

	go func() {
		defer wg.Done()
		count, err := s.productRepo.Count(ctx, nil)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			statsErr = err
			return
		}
		stats.TotalProducts = int(count)
	}()

	go func() {
		defer wg.Done()
		total, err := s.saleRepo.GetTotalSales(ctx, filter)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			statsErr = err
			return
		}
		stats.TotalSales = total
	}()

	go func() {
		defer wg.Done()
		total, err := s.commissionRepo.GetTotalCommissions(ctx, filter)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			statsErr = err
			return
		}
		stats.TotalCommissions = total
	}()

	go func() {
		defer wg.Done()
		count, err := s.clientRepo.Count(ctx, filter)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			statsErr = err
			return
		}
		stats.ActiveClients = int(count)
	}()

	wg.Wait()

	if statsErr != nil {
		return nil, statsErr
	}

	return stats, nil
}


























