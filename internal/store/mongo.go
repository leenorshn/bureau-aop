package store

import (
	"context"
	"time"

	"bureau/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
	Logger   *zap.Logger
}

func NewMongoDB(cfg *config.Config, logger *zap.Logger) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(cfg.MongoDBName)

	// Create indexes
	if err := createIndexes(ctx, db); err != nil {
		logger.Warn("Failed to create indexes", zap.Error(err))
	}

	return &MongoDB{
		Client:   client,
		Database: db,
		Logger:   logger,
	}, nil
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	// Clients indexes
	clientsCollection := db.Collection("clients")
	_, err := clientsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    map[string]interface{}{"email": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: map[string]interface{}{"sponsorId": 1},
		},
		{
			Keys: map[string]interface{}{"leftChildId": 1},
		},
		{
			Keys: map[string]interface{}{"rightChildId": 1},
		},
	})
	if err != nil {
		return err
	}

	// Sales indexes
	salesCollection := db.Collection("sales")
	_, err = salesCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"clientId": 1},
		},
		{
			Keys: map[string]interface{}{"sponsorId": 1},
		},
		{
			Keys: map[string]interface{}{"date": -1},
		},
	})
	if err != nil {
		return err
	}

	// Payments indexes
	paymentsCollection := db.Collection("payments")
	_, err = paymentsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"clientId": 1},
		},
		{
			Keys: map[string]interface{}{"date": -1},
		},
	})
	if err != nil {
		return err
	}

	// Commissions indexes
	commissionsCollection := db.Collection("commissions")
	_, err = commissionsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"clientId": 1},
		},
		{
			Keys: map[string]interface{}{"date": -1},
		},
	})
	if err != nil {
		return err
	}

	// Admins indexes
	adminsCollection := db.Collection("admins")
	_, err = adminsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    map[string]interface{}{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return err
	}

	return nil
}






