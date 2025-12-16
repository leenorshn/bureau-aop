package store

import (
	"context"
	"time"

	"bureau/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommissionRepository struct {
	collection *mongo.Collection
}

func NewCommissionRepository(db *mongo.Database) *CommissionRepository {
	return &CommissionRepository{
		collection: db.Collection("commissions"),
	}
}

func (r *CommissionRepository) Create(ctx context.Context, commission *models.Commission) (*models.Commission, error) {
	commission.Date = time.Now()

	result, err := r.collection.InsertOne(ctx, commission)
	if err != nil {
		return nil, err
	}

	commission.ID = result.InsertedID.(primitive.ObjectID)
	return commission, nil
}

func (r *CommissionRepository) GetByID(ctx context.Context, id string) (*models.Commission, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var commission models.Commission
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&commission)
	if err != nil {
		return nil, err
	}

	return &commission, nil
}

func (r *CommissionRepository) GetAll(ctx context.Context, filter *models.FilterInput, paging *models.PagingInput) ([]*models.Commission, error) {
	query := bson.M{}

	if filter != nil {
		if filter.Search != nil {
			query["$or"] = []bson.M{
				{"type": bson.M{"$regex": *filter.Search, "$options": "i"}},
			}
		}
		if filter.DateFrom != nil {
			query["date"] = bson.M{"$gte": *filter.DateFrom}
		}
		if filter.DateTo != nil {
			if query["date"] == nil {
				query["date"] = bson.M{"$lte": *filter.DateTo}
			} else {
				query["date"].(bson.M)["$lte"] = *filter.DateTo
			}
		}
	}

	opts := options.Find()
	if paging != nil {
		if paging.Limit != nil {
			opts.SetLimit(int64(*paging.Limit))
		}
		if paging.Page != nil {
			skip := int64(*paging.Page-1) * int64(*paging.Limit)
			opts.SetSkip(skip)
		}
	}
	opts.SetSort(bson.D{{"date", -1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var commissions []*models.Commission
	if err = cursor.All(ctx, &commissions); err != nil {
		return nil, err
	}

	return commissions, nil
}

func (r *CommissionRepository) GetByClientID(ctx context.Context, clientID string) ([]*models.Commission, error) {
	objectID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"clientId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var commissions []*models.Commission
	if err = cursor.All(ctx, &commissions); err != nil {
		return nil, err
	}

	return commissions, nil
}

func (r *CommissionRepository) GetBySourceClientID(ctx context.Context, sourceClientID string) ([]*models.Commission, error) {
	objectID, err := primitive.ObjectIDFromHex(sourceClientID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"sourceClientId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var commissions []*models.Commission
	if err = cursor.All(ctx, &commissions); err != nil {
		return nil, err
	}

	return commissions, nil
}

func (r *CommissionRepository) Update(ctx context.Context, id string, commission *models.Commission) (*models.Commission, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"amount": commission.Amount,
			"level":  commission.Level,
			"type":   commission.Type,
		},
	}

	var updatedCommission models.Commission
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedCommission)
	if err != nil {
		return nil, err
	}

	return &updatedCommission, nil
}

func (r *CommissionRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *CommissionRepository) Count(ctx context.Context, filter *models.FilterInput) (int64, error) {
	query := bson.M{}

	if filter != nil {
		if filter.Search != nil {
			query["$or"] = []bson.M{
				{"type": bson.M{"$regex": *filter.Search, "$options": "i"}},
			}
		}
		if filter.DateFrom != nil {
			query["date"] = bson.M{"$gte": *filter.DateFrom}
		}
		if filter.DateTo != nil {
			if query["date"] == nil {
				query["date"] = bson.M{"$lte": *filter.DateTo}
			} else {
				query["date"].(bson.M)["$lte"] = *filter.DateTo
			}
		}
	}

	return r.collection.CountDocuments(ctx, query)
}

func (r *CommissionRepository) GetTotalCommissions(ctx context.Context, filter *models.FilterInput) (float64, error) {
	query := bson.M{}

	if filter != nil {
		if filter.DateFrom != nil {
			query["date"] = bson.M{"$gte": *filter.DateFrom}
		}
		if filter.DateTo != nil {
			if query["date"] == nil {
				query["date"] = bson.M{"$lte": *filter.DateTo}
			} else {
				query["date"].(bson.M)["$lte"] = *filter.DateTo
			}
		}
	}

	pipeline := []bson.M{
		{"$match": query},
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


























