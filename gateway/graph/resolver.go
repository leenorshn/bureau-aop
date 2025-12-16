package graph

import (
	"bureau/gateway/internal/client"

	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	treeServiceClient *client.TreeServiceClient
	logger            *zap.Logger
}

func NewResolver(
	treeServiceClient *client.TreeServiceClient,
	logger *zap.Logger,
) *Resolver {
	return &Resolver{
		treeServiceClient: treeServiceClient,
		logger:            logger,
	}
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

