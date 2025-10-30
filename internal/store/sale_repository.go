package store

import (
	"context"
	"time"

	"bureau/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *SaleRepository) GetAll(ctx context.Context, filter *models.FilterInput, paging *models.PagingInput) ([]*models.Sale, error) {
	// Build filter
	filterDoc := bson.M{}
	if filter != nil {
		if filter.Search != nil {
			filterDoc["$or"] = []bson.M{
				{"note": bson.M{"$regex": *filter.Search, "$options": "i"}},
				{"status": bson.M{"$regex": *filter.Search, "$options": "i"}},
			}
		}
		if filter.DateFrom != nil {
			filterDoc["date"] = bson.M{"$gte": *filter.DateFrom}
		}
		if filter.DateTo != nil {
			if filterDoc["date"] == nil {
				filterDoc["date"] = bson.M{}
			}
			filterDoc["date"].(bson.M)["$lte"] = *filter.DateTo
		}
		if filter.Status != nil {
			filterDoc["status"] = *filter.Status
		}
	}

	// Build options
	opts := &options.FindOptions{}
	if paging != nil {
		if paging.Limit != nil {
			opts.SetLimit(int64(*paging.Limit))
		}
		if paging.Page != nil && paging.Limit != nil {
			opts.SetSkip(int64((*paging.Page - 1) * *paging.Limit))
		}
	}

	// Sort by date descending
	opts.SetSort(bson.D{{"date", -1}})

	cursor, err := r.collection.Find(ctx, filterDoc, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []*models.Sale
	if err = cursor.All(ctx, &sales); err != nil {
		return nil, err
	}

	return sales, nil
}

func (r *SaleRepository) GetByID(ctx context.Context, id string) (*models.Sale, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var sale models.Sale
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&sale)
	if err != nil {
		return nil, err
	}

	return &sale, nil
}

func (r *SaleRepository) Create(ctx context.Context, sale *models.Sale) (*models.Sale, error) {
	sale.ID = primitive.NewObjectID()
	sale.Date = time.Now()

	_, err := r.collection.InsertOne(ctx, sale)
	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (r *SaleRepository) Update(ctx context.Context, id string, sale *models.Sale) (*models.Sale, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	sale.ID = objectID
	sale.Date = time.Now()

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, sale)
	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (r *SaleRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *SaleRepository) GetByClientID(ctx context.Context, clientID string) ([]*models.Sale, error) {
	objectID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"clientId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []*models.Sale
	if err = cursor.All(ctx, &sales); err != nil {
		return nil, err
	}

	return sales, nil
}

func (r *SaleRepository) GetBySponsorID(ctx context.Context, sponsorID string) ([]*models.Sale, error) {
	objectID, err := primitive.ObjectIDFromHex(sponsorID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"sponsorId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []*models.Sale
	if err = cursor.All(ctx, &sales); err != nil {
		return nil, err
	}

	return sales, nil
}

func (r *SaleRepository) GetTotalSales(ctx context.Context, filter *models.FilterInput) (float64, error) {
	// Build filter
	filterDoc := bson.M{}
	if filter != nil {
		if filter.Search != nil {
			filterDoc["$or"] = []bson.M{
				{"note": bson.M{"$regex": *filter.Search, "$options": "i"}},
				{"status": bson.M{"$regex": *filter.Search, "$options": "i"}},
			}
		}
		if filter.DateFrom != nil {
			filterDoc["date"] = bson.M{"$gte": *filter.DateFrom}
		}
		if filter.DateTo != nil {
			if filterDoc["date"] == nil {
				filterDoc["date"] = bson.M{}
			}
			filterDoc["date"].(bson.M)["$lte"] = *filter.DateTo
		}
		if filter.Status != nil {
			filterDoc["status"] = *filter.Status
		}
	}

	pipeline := []bson.M{
		{"$match": filterDoc},
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	total, ok := result[0]["total"].(float64)
	if !ok {
		return 0, nil
	}

	return total, nil
}
