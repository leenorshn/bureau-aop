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

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	product.ID = result.InsertedID.(primitive.ObjectID)
	return product, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*models.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product models.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetAll(ctx context.Context, filter *models.FilterInput, paging *models.PagingInput) ([]*models.Product, error) {
	query := bson.M{}

	if filter != nil {
		if filter.Search != nil {
			query["$or"] = []bson.M{
				{"name": bson.M{"$regex": *filter.Search, "$options": "i"}},
				{"description": bson.M{"$regex": *filter.Search, "$options": "i"}},
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

	var products []*models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, id string, product *models.Product) (*models.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	product.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
			"points":      product.Points,
			"imageUrl":    product.ImageURL,
			"updatedAt":   product.UpdatedAt,
		},
	}

	var updatedProduct models.Product
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedProduct)
	if err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *ProductRepository) Count(ctx context.Context, filter *models.FilterInput) (int64, error) {
	query := bson.M{}

	if filter != nil {
		if filter.Search != nil {
			query["$or"] = []bson.M{
				{"name": bson.M{"$regex": *filter.Search, "$options": "i"}},
				{"description": bson.M{"$regex": *filter.Search, "$options": "i"}},
			}
		}
	}

	return r.collection.CountDocuments(ctx, query)
}






