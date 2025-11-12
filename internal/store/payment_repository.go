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

type PaymentRepository struct {
	collection *mongo.Collection
}

func NewPaymentRepository(db *mongo.Database) *PaymentRepository {
	return &PaymentRepository{
		collection: db.Collection("payments"),
	}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	payment.Date = time.Now()

	result, err := r.collection.InsertOne(ctx, payment)
	if err != nil {
		return nil, err
	}

	payment.ID = result.InsertedID.(primitive.ObjectID)
	return payment, nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, id string) (*models.Payment, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var payment models.Payment
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&payment)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *PaymentRepository) GetAll(ctx context.Context, filter *models.FilterInput, paging *models.PagingInput) ([]*models.Payment, error) {
	query := bson.M{}

	if filter != nil {
		if filter.Search != nil {
			query["$or"] = []bson.M{
				{"method": bson.M{"$regex": *filter.Search, "$options": "i"}},
				{"description": bson.M{"$regex": *filter.Search, "$options": "i"}},
			}
		}
		if filter.Status != nil {
			query["status"] = *filter.Status
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

	var payments []*models.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepository) GetByClientID(ctx context.Context, clientID string) ([]*models.Payment, error) {
	objectID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"clientId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []*models.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepository) Update(ctx context.Context, id string, payment *models.Payment) (*models.Payment, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"amount":      payment.Amount,
			"method":      payment.Method,
			"status":      payment.Status,
			"description": payment.Description,
		},
	}

	var updatedPayment models.Payment
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedPayment)
	if err != nil {
		return nil, err
	}

	return &updatedPayment, nil
}

func (r *PaymentRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *PaymentRepository) Count(ctx context.Context, filter *models.FilterInput) (int64, error) {
	query := bson.M{}

	if filter != nil {
		if filter.Search != nil {
			query["$or"] = []bson.M{
				{"method": bson.M{"$regex": *filter.Search, "$options": "i"}},
				{"description": bson.M{"$regex": *filter.Search, "$options": "i"}},
			}
		}
		if filter.Status != nil {
			query["status"] = *filter.Status
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

func (r *PaymentRepository) GetTotalPayments(ctx context.Context, filter *models.FilterInput) (float64, error) {
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





















