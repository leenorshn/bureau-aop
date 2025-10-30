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

type AdminRepository struct {
	collection *mongo.Collection
}

func NewAdminRepository(db *mongo.Database) *AdminRepository {
	return &AdminRepository{
		collection: db.Collection("admins"),
	}
}

func (r *AdminRepository) Create(ctx context.Context, admin *models.Admin) (*models.Admin, error) {
	admin.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, admin)
	if err != nil {
		return nil, err
	}

	admin.ID = result.InsertedID.(primitive.ObjectID)
	return admin, nil
}

func (r *AdminRepository) GetByID(ctx context.Context, id string) (*models.Admin, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var admin models.Admin
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&admin)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (r *AdminRepository) GetByEmail(ctx context.Context, email string) (*models.Admin, error) {
	var admin models.Admin
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&admin)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (r *AdminRepository) GetAll(ctx context.Context, filter *models.FilterInput, paging *models.PagingInput) ([]*models.Admin, error) {
	query := bson.M{}

	if filter != nil {
		if filter.Search != nil {
			query["$or"] = []bson.M{
				{"name": bson.M{"$regex": *filter.Search, "$options": "i"}},
				{"email": bson.M{"$regex": *filter.Search, "$options": "i"}},
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
	opts.SetSort(bson.D{{"createdAt", -1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var admins []*models.Admin
	if err = cursor.All(ctx, &admins); err != nil {
		return nil, err
	}

	return admins, nil
}

func (r *AdminRepository) Update(ctx context.Context, id string, admin *models.Admin) (*models.Admin, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"name":  admin.Name,
			"email": admin.Email,
			"role":  admin.Role,
		},
	}

	var updatedAdmin models.Admin
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedAdmin)
	if err != nil {
		return nil, err
	}

	return &updatedAdmin, nil
}

func (r *AdminRepository) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$set": bson.M{
			"passwordHash": passwordHash,
		},
	})
	return err
}

func (r *AdminRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *AdminRepository) Count(ctx context.Context, filter *models.FilterInput) (int64, error) {
	query := bson.M{}

	if filter != nil {
		if filter.Search != nil {
			query["$or"] = []bson.M{
				{"name": bson.M{"$regex": *filter.Search, "$options": "i"}},
				{"email": bson.M{"$regex": *filter.Search, "$options": "i"}},
			}
		}
	}

	return r.collection.CountDocuments(ctx, query)
}













