package store

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// TransactionHelper gère les transactions MongoDB atomiques
type TransactionHelper struct {
	client *mongo.Client
}

// NewTransactionHelper crée un nouveau helper pour les transactions
func NewTransactionHelper(client *mongo.Client) *TransactionHelper {
	return &TransactionHelper{
		client: client,
	}
}

// ExecuteTransaction exécute une fonction dans une transaction MongoDB
// Si les transactions ne sont pas supportées (standalone), exécute sans transaction
func (h *TransactionHelper) ExecuteTransaction(ctx context.Context, fn func(context.Context) error) error {
	// Démarrer une session
	session, err := h.client.StartSession()
	if err != nil {
		// Si les sessions ne sont pas supportées, exécuter sans transaction
		return fn(ctx)
	}
	defer session.EndSession(ctx)

	// Exécuter dans une transaction
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			// Si les transactions ne sont pas supportées, exécuter sans transaction
			return fn(ctx)
		}

		// Exécuter la fonction
		if err := fn(sc); err != nil {
			// Rollback en cas d'erreur
			if abortErr := session.AbortTransaction(sc); abortErr != nil {
				return abortErr
			}
			return err
		}

		// Commit la transaction
		if err := session.CommitTransaction(sc); err != nil {
			return err
		}

		return nil
	})

	return err
}

// GetSessionContext retourne un contexte avec session pour les opérations atomiques
func (h *TransactionHelper) GetSessionContext(ctx context.Context) (context.Context, mongo.Session, error) {
	session, err := h.client.StartSession()
	if err != nil {
		return ctx, nil, err
	}

	sessionCtx := mongo.NewSessionContext(ctx, session)
	return sessionCtx, session, nil
}

