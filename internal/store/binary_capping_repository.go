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

// BinaryCappingRepository gère les limites journalières/hebdomadaires des cycles
type BinaryCappingRepository struct {
	collection *mongo.Collection
}

// NewBinaryCappingRepository crée un nouveau repository pour le capping
func NewBinaryCappingRepository(db *mongo.Database) *BinaryCappingRepository {
	return &BinaryCappingRepository{
		collection: db.Collection("binary_capping"),
	}
}

// GetByClientIDAndDate récupère le capping pour un client et une date donnée
func (r *BinaryCappingRepository) GetByClientIDAndDate(ctx context.Context, clientID primitive.ObjectID, date time.Time) (*models.BinaryCapping, error) {
	dateStart := date.Truncate(24 * time.Hour)
	
	var capping models.BinaryCapping
	err := r.collection.FindOne(ctx, bson.M{
		"clientId": clientID,
		"date":     dateStart,
	}).Decode(&capping)

	if err == mongo.ErrNoDocuments {
		// Créer un nouveau capping
		capping = models.BinaryCapping{
			ID:                primitive.NewObjectID(),
			ClientID:          clientID,
			Date:              dateStart,
			CyclesPaidToday:   0,
			CyclesPaidThisWeek: 0,
			LastResetDate:     dateStart,
		}
		_, err = r.collection.InsertOne(ctx, capping)
		if err != nil {
			return nil, err
		}
		return &capping, nil
	}

	if err != nil {
		return nil, err
	}

	// Vérifier si on doit reset (nouveau jour)
	if capping.LastResetDate.Before(dateStart) {
		capping.CyclesPaidToday = 0
		capping.LastResetDate = dateStart
		err = r.Update(ctx, &capping)
		if err != nil {
			return nil, err
		}
	}

	return &capping, nil
}

// Update met à jour un capping
func (r *BinaryCappingRepository) Update(ctx context.Context, capping *models.BinaryCapping) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": capping.ID},
		bson.M{"$set": bson.M{
			"cyclesPaidToday":    capping.CyclesPaidToday,
			"cyclesPaidThisWeek": capping.CyclesPaidThisWeek,
			"lastResetDate":      capping.LastResetDate,
		}},
		options.Update().SetUpsert(true),
	)
	return err
}

// IncrementCycles incrémente le nombre de cycles payés
func (r *BinaryCappingRepository) IncrementCycles(ctx context.Context, clientID primitive.ObjectID, date time.Time, cycles int) error {
	dateStart := date.Truncate(24 * time.Hour)
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"clientId": clientID,
			"date":     dateStart,
		},
		bson.M{
			"$inc": bson.M{
				"cyclesPaidToday": cycles,
			},
			"$setOnInsert": bson.M{
				"clientId":          clientID,
				"date":              dateStart,
				"cyclesPaidThisWeek": 0,
				"lastResetDate":     dateStart,
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}










