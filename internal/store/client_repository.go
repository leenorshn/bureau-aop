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

type ClientRepository struct {
	collection *mongo.Collection
}

func NewClientRepository(db *mongo.Database) *ClientRepository {
	return &ClientRepository{
		collection: db.Collection("clients"),
	}
}

func (r *ClientRepository) Create(ctx context.Context, client *models.Client) (*models.Client, error) {
	client.JoinDate = time.Now()

	result, err := r.collection.InsertOne(ctx, client)
	if err != nil {
		return nil, err
	}

	client.ID = result.InsertedID.(primitive.ObjectID)
	return client, nil
}

func (r *ClientRepository) GetByID(ctx context.Context, id string) (*models.Client, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var client models.Client
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&client)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *ClientRepository) GetByClientID(ctx context.Context, clientID string) (*models.Client, error) {
	var client models.Client
	err := r.collection.FindOne(ctx, bson.M{"clientId": clientID}).Decode(&client)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *ClientRepository) GetAll(ctx context.Context, filter *models.FilterInput, paging *models.PagingInput) ([]*models.Client, error) {
	query := bson.M{}

	if filter != nil {
		if filter.Search != nil {
			query["$or"] = []bson.M{
				{"name": bson.M{"$regex": *filter.Search, "$options": "i"}},
				{"clientId": bson.M{"$regex": *filter.Search, "$options": "i"}},
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
	opts.SetSort(bson.D{{"joinDate", -1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var clients []*models.Client
	if err = cursor.All(ctx, &clients); err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *ClientRepository) Update(ctx context.Context, id string, client *models.Client) (*models.Client, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"name": client.Name,
		},
	}

	var updatedClient models.Client
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedClient)
	if err != nil {
		return nil, err
	}

	return &updatedClient, nil
}

func (r *ClientRepository) UpdateBinaryFields(ctx context.Context, id string, leftChildID, rightChildID *primitive.ObjectID, position *string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{}
	if leftChildID != nil {
		update["leftChildId"] = leftChildID
	}
	if rightChildID != nil {
		update["rightChildId"] = rightChildID
	}
	if position != nil {
		update["position"] = position
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": update})
	return err
}

func (r *ClientRepository) UpdateNetworkVolumes(ctx context.Context, id string, leftVolume, rightVolume float64) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$set": bson.M{
			"networkVolumeLeft":  leftVolume,
			"networkVolumeRight": rightVolume,
		},
	})
	return err
}

func (r *ClientRepository) UpdateEarnings(ctx context.Context, id string, totalEarnings, walletBalance float64) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$set": bson.M{
			"totalEarnings": totalEarnings,
			"walletBalance": walletBalance,
		},
	})
	return err
}

func (r *ClientRepository) UpdatePoints(ctx context.Context, id string, points float64) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$set": bson.M{
			"points": points,
		},
	})
	return err
}

func (r *ClientRepository) AddPoints(ctx context.Context, id string, pointsToAdd float64) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$inc": bson.M{
			"points": pointsToAdd,
		},
	})
	return err
}

func (r *ClientRepository) UpdateBinaryPairs(ctx context.Context, id string, pairs int) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$set": bson.M{
			"binaryPairs": pairs,
		},
	})
	return err
}

func (r *ClientRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *ClientRepository) Count(ctx context.Context, filter *models.FilterInput) (int64, error) {
	query := bson.M{}

	if filter != nil {
		if filter.Search != nil {
			query["$or"] = []bson.M{
				{"name": bson.M{"$regex": *filter.Search, "$options": "i"}},
				{"clientId": bson.M{"$regex": *filter.Search, "$options": "i"}},
			}
		}
	}

	return r.collection.CountDocuments(ctx, query)
}

func (r *ClientRepository) GetBySponsorID(ctx context.Context, sponsorID string) ([]*models.Client, error) {
	objectID, err := primitive.ObjectIDFromHex(sponsorID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"sponsorId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var clients []*models.Client
	if err = cursor.All(ctx, &clients); err != nil {
		return nil, err
	}

	return clients, nil
}
