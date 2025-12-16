package store

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type SaleRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewSaleRepository(db *mongo.Database, logger *zap.Logger) *SaleRepository {
	return &SaleRepository{
		collection: db.Collection("sales"),
		logger:     logger,
	}
}

func (r *SaleRepository) GetByClientID(ctx context.Context, clientID string) ([]interface{}, error) {
	objectID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"clientId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []interface{}
	if err := cursor.All(ctx, &sales); err != nil {
		return nil, err
	}

	return sales, nil
}



