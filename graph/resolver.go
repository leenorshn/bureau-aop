package graph

import (
	"bureau/graph/model"
	"bureau/internal/service"
	"context"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// ClientResolver defines the interface for Client field resolvers
type ClientResolver interface {
	Sponsor(ctx context.Context, obj *model.Client) (*model.Client, error)
	LeftChild(ctx context.Context, obj *model.Client) (*model.Client, error)
	RightChild(ctx context.Context, obj *model.Client) (*model.Client, error)
	Transactions(ctx context.Context, obj *model.Client) ([]*model.Payment, error)
	Purchases(ctx context.Context, obj *model.Client) ([]*model.Sale, error)
}

type Resolver struct {
	productService          *service.ProductService
	clientService           *service.ClientService
	saleService             *service.SaleService
	paymentService          *service.PaymentService
	commissionService       *service.CommissionService
	authService             *service.AuthService
	adminService            *service.AdminService
	caisseService           *service.CaisseService
	binaryCommissionService *service.BinaryCommissionService
}

func NewResolver(
	productService *service.ProductService,
	clientService *service.ClientService,
	saleService *service.SaleService,
	paymentService *service.PaymentService,
	commissionService *service.CommissionService,
	authService *service.AuthService,
	adminService *service.AdminService,
	caisseService *service.CaisseService,
	binaryCommissionService *service.BinaryCommissionService,
) *Resolver {
	return &Resolver{
		productService:          productService,
		clientService:           clientService,
		saleService:             saleService,
		paymentService:          paymentService,
		commissionService:       commissionService,
		authService:             authService,
		adminService:            adminService,
		caisseService:           caisseService,
		binaryCommissionService: binaryCommissionService,
	}
}
