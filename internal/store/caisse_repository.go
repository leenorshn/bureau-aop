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

type CaisseRepository struct {
	collection         *mongo.Collection
	transactionCollection *mongo.Collection
}

func NewCaisseRepository(db *mongo.Database) *CaisseRepository {
	return &CaisseRepository{
		collection:         db.Collection("caisse"),
		transactionCollection: db.Collection("caisse_transactions"),
	}
}

// GetOrCreate gets the caisse or creates it if it doesn't exist
func (r *CaisseRepository) GetOrCreate(ctx context.Context) (*models.Caisse, error) {
	var caisse models.Caisse
	err := r.collection.FindOne(ctx, bson.M{}).Decode(&caisse)
	if err == mongo.ErrNoDocuments {
		// Create initial caisse
		caisse = models.Caisse{
			ID:           primitive.NewObjectID(),
			Balance:      0,
			TotalEntrees: 0,
			TotalSorties: 0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err = r.collection.InsertOne(ctx, caisse)
		if err != nil {
			return nil, err
		}
		return &caisse, nil
	}
	if err != nil {
		return nil, err
	}
	return &caisse, nil
}

// UpdateBalance updates the caisse balance
func (r *CaisseRepository) UpdateBalance(ctx context.Context, balance, totalEntrees, totalSorties float64) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{},
		bson.M{
			"$set": bson.M{
				"balance":      balance,
				"totalEntrees": totalEntrees,
				"totalSorties": totalSorties,
				"updatedAt":    time.Now(),
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

// AddTransaction adds a transaction to the caisse
func (r *CaisseRepository) AddTransaction(ctx context.Context, transaction *models.CaisseTransaction) (*models.CaisseTransaction, error) {
	transaction.ID = primitive.NewObjectID()
	transaction.Date = time.Now()

	result, err := r.transactionCollection.InsertOne(ctx, transaction)
	if err != nil {
		return nil, err
	}

	transaction.ID = result.InsertedID.(primitive.ObjectID)
	return transaction, nil
}

// GetTransactions gets all transactions with optional filtering
func (r *CaisseRepository) GetTransactions(ctx context.Context, filter *models.FilterInput, paging *models.PagingInput) ([]*models.CaisseTransaction, error) {
	query := bson.M{}

	if filter != nil {
		if filter.Status != nil {
			query["type"] = *filter.Status // Use status field to filter by type (entree/sortie)
		}
		if filter.DateFrom != nil {
			query["date"] = bson.M{"$gte": *filter.DateFrom}
		}
		if filter.DateTo != nil {
			if query["date"] == nil {
				query["date"] = bson.M{}
			}
			query["date"].(bson.M)["$lte"] = *filter.DateTo
		}
	}

	opts := options.Find()
	if paging != nil {
		if paging.Limit != nil {
			opts.SetLimit(int64(*paging.Limit))
		}
		if paging.Page != nil && paging.Limit != nil {
			skip := int64(*paging.Page-1) * int64(*paging.Limit)
			opts.SetSkip(skip)
		}
	}
	opts.SetSort(bson.D{{Key: "date", Value: -1}})

	cursor, err := r.transactionCollection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*models.CaisseTransaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetTransactionByID gets a transaction by ID
func (r *CaisseRepository) GetTransactionByID(ctx context.Context, id string) (*models.CaisseTransaction, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var transaction models.CaisseTransaction
	err = r.transactionCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

