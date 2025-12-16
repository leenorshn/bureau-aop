package store

import (
	"context"
	"fmt"

	"bureau/services/tree-service/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientRepository struct {
	collection *mongo.Collection
}

func NewClientRepository(db *mongo.Database) *ClientRepository {
	return &ClientRepository{
		collection: db.Collection("clients"),
	}
}

func (r *ClientRepository) GetByID(ctx context.Context, id string) (*models.Client, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid client ID: %w", err)
	}

	var clientData struct {
		ID                 primitive.ObjectID  `bson:"_id"`
		ClientID           string              `bson:"clientId"`
		Name               string              `bson:"name"`
		Phone              *string             `bson:"phone,omitempty"`
		Position           *string             `bson:"position,omitempty"`
		LeftChildID        *primitive.ObjectID `bson:"leftChildId,omitempty"`
		RightChildID       *primitive.ObjectID `bson:"rightChildId,omitempty"`
		NetworkVolumeLeft  float64             `bson:"networkVolumeLeft"`
		NetworkVolumeRight float64             `bson:"networkVolumeRight"`
		BinaryPairs        int                 `bson:"binaryPairs"`
		TotalEarnings      float64             `bson:"totalEarnings"`
		WalletBalance      float64             `bson:"walletBalance"`
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&clientData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("client not found")
		}
		return nil, err
	}

	// Convertir en modèle Client
	client := &models.Client{
		ID:                 clientData.ID.Hex(),
		ClientID:           clientData.ClientID,
		Name:               clientData.Name,
		Phone:              clientData.Phone,
		Position:           clientData.Position,
		NetworkVolumeLeft:  clientData.NetworkVolumeLeft,
		NetworkVolumeRight: clientData.NetworkVolumeRight,
		BinaryPairs:        clientData.BinaryPairs,
		TotalEarnings:     clientData.TotalEarnings,
		WalletBalance:      clientData.WalletBalance,
	}

	// Convertir les ObjectID en string pour les enfants
	if clientData.LeftChildID != nil {
		leftIDStr := clientData.LeftChildID.Hex()
		client.LeftChildID = &leftIDStr
	}
	if clientData.RightChildID != nil {
		rightIDStr := clientData.RightChildID.Hex()
		client.RightChildID = &rightIDStr
	}

	return client, nil
}

