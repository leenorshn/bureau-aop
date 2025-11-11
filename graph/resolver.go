package graph

import (
	"bureau/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	productService    *service.ProductService
	clientService     *service.ClientService
	saleService       *service.SaleService
	paymentService    *service.PaymentService
	commissionService *service.CommissionService
	authService       *service.AuthService
	adminService      *service.AdminService
}

func NewResolver(
	productService *service.ProductService,
	clientService *service.ClientService,
	saleService *service.SaleService,
	paymentService *service.PaymentService,
	commissionService *service.CommissionService,
	authService *service.AuthService,
	adminService *service.AdminService,
) *Resolver {
	return &Resolver{
		productService:    productService,
		clientService:     clientService,
		saleService:       saleService,
		paymentService:    paymentService,
		commissionService: commissionService,
		authService:       authService,
		adminService:      adminService,
	}
}
