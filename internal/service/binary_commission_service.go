package service

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"bureau/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// Interfaces pour permettre l'utilisation de mocks dans les tests
type clientRepository interface {
	GetByID(ctx context.Context, id string) (*models.Client, error)
	UpdateNetworkVolumes(ctx context.Context, id string, left, right float64) error
	UpdateEarnings(ctx context.Context, id string, totalEarnings, walletBalance float64) error
}

type commissionRepository interface {
	Create(ctx context.Context, commission *models.Commission) (*models.Commission, error)
}

type saleRepository interface {
	GetByClientID(ctx context.Context, clientID string) ([]*models.Sale, error)
}

type binaryCappingRepository interface {
	GetByClientIDAndDate(ctx context.Context, clientID primitive.ObjectID, date time.Time) (*models.BinaryCapping, error)
	Update(ctx context.Context, capping *models.BinaryCapping) error
	IncrementCycles(ctx context.Context, clientID primitive.ObjectID, date time.Time, cycles int) error
}

// BinaryCommissionService gère le calcul et le paiement des commissions binaires MLM
type BinaryCommissionService struct {
	clientRepo     clientRepository
	commissionRepo commissionRepository
	saleRepo       saleRepository
	cappingRepo    binaryCappingRepository
	logger         *zap.Logger
	config         models.BinaryConfig
	mu             sync.Mutex        // Pour éviter les doubles paiements (fallback si transactions non disponibles)
	txHelper       transactionHelper // Helper pour les transactions atomiques
}

// transactionHelper interface pour les transactions
type transactionHelper interface {
	ExecuteTransaction(ctx context.Context, fn func(context.Context) error) error
}

// NewBinaryCommissionService crée un nouveau service de commission binaire
func NewBinaryCommissionService(
	clientRepo clientRepository,
	commissionRepo commissionRepository,
	saleRepo saleRepository,
	cappingRepo binaryCappingRepository,
	logger *zap.Logger,
	config models.BinaryConfig,
	txHelper transactionHelper,
) *BinaryCommissionService {
	return &BinaryCommissionService{
		clientRepo:     clientRepo,
		commissionRepo: commissionRepo,
		saleRepo:       saleRepo,
		cappingRepo:    cappingRepo,
		logger:         logger,
		config:         config,
		txHelper:       txHelper,
	}
}

// ComputeBinaryCommission calcule et paie la commission binaire pour un membre
// C'est la fonction principale qui orchestre tout le processus
func (s *BinaryCommissionService) ComputeBinaryCommission(ctx context.Context, clientID string) (*models.BinaryCommissionResult, error) {
	// 1. Vérifier que le client existe
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return &models.BinaryCommissionResult{
			Success: false,
			Reason:  "Client introuvable",
		}, err
	}

	// 2. Vérifier la qualification
	qualification, err := s.checkQualification(ctx, client)
	if err != nil {
		return &models.BinaryCommissionResult{
			Success: false,
			Reason:  fmt.Sprintf("Erreur lors de la vérification de qualification: %v", err),
		}, err
	}

	if !qualification.IsQualified {
		return &models.BinaryCommissionResult{
			Success:   true,
			Qualified: false,
			Reason:    "Membre non qualifié: doit avoir au moins 1 direct actif à gauche ET 1 direct actif à droite",
		}, nil
	}

	// 3. Lire les volumes des jambes
	legs, err := s.getLegsVolumes(ctx, client)
	if err != nil {
		return &models.BinaryCommissionResult{
			Success: false,
			Reason:  fmt.Sprintf("Erreur lors de la lecture des volumes: %v", err),
		}, err
	}

	// 4. Vérifier les conditions de base (au moins 1 actif de chaque côté)
	if legs.LeftActives == 0 || legs.RightActives == 0 {
		return &models.BinaryCommissionResult{
			Success:              true,
			Qualified:            true,
			Reason:               "Jambe gauche ou droite vide - aucun cycle possible",
			LeftVolumeRemaining:  legs.LeftVolume,
			RightVolumeRemaining: legs.RightVolume,
		}, nil
	}

	// 5. Calculer les cycles possibles à partir des volumes
	cyclesAvailable := s.calculateCycles(legs)

	if cyclesAvailable == 0 {
		return &models.BinaryCommissionResult{
			Success:              true,
			Qualified:            true,
			Reason:               "Aucun cycle disponible - volumes insuffisants",
			LeftVolumeRemaining:  legs.LeftVolume,
			RightVolumeRemaining: legs.RightVolume,
		}, nil
	}

	// 6. Appliquer la limite journalière
	cyclesToPay, err := s.applyDailyLimit(ctx, client.ID, cyclesAvailable)
	if err != nil {
		return &models.BinaryCommissionResult{
			Success: false,
			Reason:  fmt.Sprintf("Erreur lors de l'application de la limite: %v", err),
		}, err
	}

	if cyclesToPay == 0 {
		return &models.BinaryCommissionResult{
			Success:              true,
			Qualified:            true,
			CyclesAvailable:      cyclesAvailable,
			CyclesPaid:           0,
			Reason:               "Limite journalière atteinte",
			LeftVolumeRemaining:  legs.LeftVolume,
			RightVolumeRemaining: legs.RightVolume,
		}, nil
	}

	// 7. Calculer le montant basé sur le volume faible
	minVolumePerLeg := s.getMinVolumePerLeg()
	volumeUsed := float64(cyclesToPay) * minVolumePerLeg
	amount := s.calculateAmount(volumeUsed)

	// 8. Enregistrer le paiement avec transaction atomique
	var commission *models.Commission
	var cyclesToPayFinal int
	var leftRemaining, rightRemaining float64
	var commissionID string

	// Utiliser une transaction atomique pour toutes les opérations critiques
	if s.txHelper != nil {
		err = s.txHelper.ExecuteTransaction(ctx, func(txCtx context.Context) error {
			// Double vérification de la limite journalière dans la transaction
			var err error
			cyclesToPayFinal, err = s.applyDailyLimit(txCtx, client.ID, cyclesAvailable)
			if err != nil {
				return fmt.Errorf("erreur lors de la vérification de la limite: %w", err)
			}

			if cyclesToPayFinal == 0 {
				return nil // Pas d'erreur, juste pas de cycles à payer
			}

			minVolumePerLeg := s.getMinVolumePerLeg()
			volumeUsed = float64(cyclesToPayFinal) * minVolumePerLeg
			amount = s.calculateAmount(volumeUsed)

			// Créer la commission
			commission, err = s.recordPayment(txCtx, client.ID, cyclesToPayFinal, amount)
			if err != nil {
				return fmt.Errorf("erreur lors de l'enregistrement du paiement: %w", err)
			}

			// Déduire les volumes utilisés
			leftRemaining, rightRemaining, err = s.deductVolume(txCtx, client.ID, legs, volumeUsed)
			if err != nil {
				return fmt.Errorf("erreur lors de la déduction des volumes: %w", err)
			}

			// Mettre à jour les gains du client
			err = s.updateClientEarnings(txCtx, client.ID.Hex(), amount)
			if err != nil {
				return fmt.Errorf("erreur lors de la mise à jour des gains: %w", err)
			}

			return nil
		})

		if err != nil {
			return &models.BinaryCommissionResult{
				Success: false,
				Reason:  fmt.Sprintf("Erreur lors de la transaction atomique: %v", err),
			}, err
		}

		if cyclesToPayFinal == 0 {
			return &models.BinaryCommissionResult{
				Success:              true,
				Qualified:            true,
				CyclesAvailable:      cyclesAvailable,
				CyclesPaid:           0,
				Reason:               "Limite journalière atteinte",
				LeftVolumeRemaining:  legs.LeftVolume,
				RightVolumeRemaining: legs.RightVolume,
			}, nil
		}

		commissionID = commission.ID.Hex()
	} else {
		// Fallback: utiliser mutex si transactions non disponibles
		s.mu.Lock()
		defer s.mu.Unlock()

		// Double vérification après verrouillage
		cyclesToPayFinal, err = s.applyDailyLimit(ctx, client.ID, cyclesAvailable)
		if err != nil {
			return &models.BinaryCommissionResult{
				Success: false,
				Reason:  fmt.Sprintf("Erreur lors de la double vérification: %v", err),
			}, err
		}

		if cyclesToPayFinal == 0 {
			return &models.BinaryCommissionResult{
				Success:              true,
				Qualified:            true,
				CyclesAvailable:      cyclesAvailable,
				CyclesPaid:           0,
				Reason:               "Limite journalière atteinte (double vérification)",
				LeftVolumeRemaining:  legs.LeftVolume,
				RightVolumeRemaining: legs.RightVolume,
			}, nil
		}

		minVolumePerLeg := s.getMinVolumePerLeg()
		volumeUsed = float64(cyclesToPayFinal) * minVolumePerLeg
		amount = s.calculateAmount(volumeUsed)

		// Créer la commission
		commission, err = s.recordPayment(ctx, client.ID, cyclesToPayFinal, amount)
		if err != nil {
			return &models.BinaryCommissionResult{
				Success: false,
				Reason:  fmt.Sprintf("Erreur lors de l'enregistrement du paiement: %v", err),
			}, err
		}

		// Déduire les volumes utilisés
		leftRemaining, rightRemaining, err = s.deductVolume(ctx, client.ID, legs, volumeUsed)
		if err != nil {
			return &models.BinaryCommissionResult{
				Success: false,
				Reason:  fmt.Sprintf("Erreur lors de la déduction des volumes: %v", err),
			}, err
		}

		// Mettre à jour les gains du client
		err = s.updateClientEarnings(ctx, client.ID.Hex(), amount)
		if err != nil {
			s.logger.Error("Failed to update client earnings", zap.Error(err))
			// Ne pas échouer complètement si c'est juste la mise à jour des gains
		}

		commissionID = commission.ID.Hex()
	}

	return &models.BinaryCommissionResult{
		Success:              true,
		Qualified:            true,
		CyclesAvailable:      cyclesAvailable,
		CyclesPaid:           cyclesToPayFinal,
		Amount:               amount,
		LeftVolumeRemaining:  leftRemaining,
		RightVolumeRemaining: rightRemaining,
		CommissionID:         &commissionID,
	}, nil
}

// checkQualification vérifie si un membre est qualifié pour recevoir des commissions
// Qualification = avoir au moins 1 direct actif à gauche ET 1 direct actif à droite
func (s *BinaryCommissionService) checkQualification(ctx context.Context, client *models.Client) (*models.BinaryQualification, error) {
	qualification := &models.BinaryQualification{}

	// Vérifier les directs à gauche
	if client.LeftChildID != nil {
		leftChild, err := s.clientRepo.GetByID(ctx, client.LeftChildID.Hex())
		if err == nil && leftChild != nil {
			// Un direct est actif s'il a fait au moins 1 vente
			isActive, err := s.isClientActive(ctx, leftChild.ID.Hex())
			if err == nil && isActive {
				qualification.HasDirectLeft = true
				qualification.DirectLeftCount = 1
			}
		}
	}

	// Vérifier les directs à droite
	if client.RightChildID != nil {
		rightChild, err := s.clientRepo.GetByID(ctx, client.RightChildID.Hex())
		if err == nil && rightChild != nil {
			isActive, err := s.isClientActive(ctx, rightChild.ID.Hex())
			if err == nil && isActive {
				qualification.HasDirectRight = true
				qualification.DirectRightCount = 1
			}
		}
	}

	// Qualification = avoir les deux
	qualification.IsQualified = qualification.HasDirectLeft && qualification.HasDirectRight

	return qualification, nil
}

// isClientActive vérifie si un client est actif (a fait au moins 1 vente)
func (s *BinaryCommissionService) isClientActive(ctx context.Context, clientID string) (bool, error) {
	sales, err := s.saleRepo.GetByClientID(ctx, clientID)
	if err != nil {
		return false, err
	}
	return len(sales) > 0, nil
}

// getLegsVolumes récupère les volumes et actifs des jambes gauche et droite
func (s *BinaryCommissionService) getLegsVolumes(ctx context.Context, client *models.Client) (*models.BinaryLegs, error) {
	legs := &models.BinaryLegs{
		LeftVolume:  client.NetworkVolumeLeft,
		RightVolume: client.NetworkVolumeRight,
	}

	// Compter les actifs dans chaque jambe
	leftActives, err := s.countActivesInLeg(ctx, client.LeftChildID, "left")
	if err != nil {
		return nil, err
	}
	legs.LeftActives = leftActives

	rightActives, err := s.countActivesInLeg(ctx, client.RightChildID, "right")
	if err != nil {
		return nil, err
	}
	legs.RightActives = rightActives

	return legs, nil
}

// countActivesInLeg compte récursivement les membres actifs dans une jambe
func (s *BinaryCommissionService) countActivesInLeg(ctx context.Context, rootID *primitive.ObjectID, side string) (int, error) {
	if rootID == nil {
		return 0, nil
	}

	count := 0
	visited := make(map[string]bool) // Pour éviter de compter deux fois
	queue := []*primitive.ObjectID{rootID}

	for len(queue) > 0 {
		currentID := queue[0]
		queue = queue[1:]

		// Éviter les doublons
		if visited[currentID.Hex()] {
			continue
		}
		visited[currentID.Hex()] = true

		client, err := s.clientRepo.GetByID(ctx, currentID.Hex())
		if err != nil {
			continue // Ignorer les erreurs et continuer
		}

		// Vérifier si ce client est actif
		isActive, err := s.isClientActive(ctx, currentID.Hex())
		if err == nil && isActive {
			count++
		}

		// Ajouter tous les enfants à la queue (pour compter récursivement)
		if client.LeftChildID != nil {
			queue = append(queue, client.LeftChildID)
		}
		if client.RightChildID != nil {
			queue = append(queue, client.RightChildID)
		}
	}

	return count, nil
}

// calculateCycles calcule le nombre de cycles possibles
// cycles = floor(min(leftVolume, rightVolume) / minVolumePerLeg)
func (s *BinaryCommissionService) calculateCycles(legs *models.BinaryLegs) int {
	if legs.LeftVolume <= 0 || legs.RightVolume <= 0 {
		return 0
	}

	minVolumePerLeg := s.getMinVolumePerLeg()
	weakVolume := math.Min(legs.LeftVolume, legs.RightVolume)
	return int(math.Floor(weakVolume / minVolumePerLeg))
}

// applyDailyLimit applique la limite journalière de cycles
func (s *BinaryCommissionService) applyDailyLimit(ctx context.Context, clientID primitive.ObjectID, cyclesAvailable int) (int, error) {
	if s.config.DailyCycleLimit <= 0 {
		return cyclesAvailable, nil // Pas de limite
	}

	// Récupérer ou créer le capping pour aujourd'hui
	today := time.Now().Truncate(24 * time.Hour)
	capping, err := s.getOrCreateCapping(ctx, clientID, today)
	if err != nil {
		return 0, err
	}

	// Vérifier si la limite est atteinte
	if capping.CyclesPaidToday >= s.config.DailyCycleLimit {
		return 0, nil
	}

	// Calculer combien de cycles on peut encore payer
	remainingLimit := s.config.DailyCycleLimit - capping.CyclesPaidToday
	cyclesToPay := int(math.Min(float64(cyclesAvailable), float64(remainingLimit)))

	// Mettre à jour le capping (incrémenter les cycles payés)
	err = s.cappingRepo.IncrementCycles(ctx, clientID, today, cyclesToPay)
	if err != nil {
		return 0, err
	}

	return cyclesToPay, nil
}

// getOrCreateCapping récupère ou crée un enregistrement de capping
func (s *BinaryCommissionService) getOrCreateCapping(ctx context.Context, clientID primitive.ObjectID, date time.Time) (*models.BinaryCapping, error) {
	return s.cappingRepo.GetByClientIDAndDate(ctx, clientID, date)
}

// updateCapping met à jour le capping dans la DB
func (s *BinaryCommissionService) updateCapping(ctx context.Context, capping *models.BinaryCapping) error {
	return s.cappingRepo.Update(ctx, capping)
}

// recordPayment enregistre le paiement de commission
func (s *BinaryCommissionService) recordPayment(ctx context.Context, clientID primitive.ObjectID, cycles int, amount float64) (*models.Commission, error) {
	commission := &models.Commission{
		ID:             primitive.NewObjectID(),
		ClientID:       clientID,
		SourceClientID: clientID, // Auto-commission pour binaire
		Amount:         amount,
		Level:          0,
		Type:           "binary-cycle",
		Date:           time.Now(),
	}

	created, err := s.commissionRepo.Create(ctx, commission)
	if err != nil {
		return nil, fmt.Errorf("failed to create commission: %w", err)
	}

	// Enregistrer aussi dans BinaryCycle pour l'historique
	// TODO: Créer un repository pour BinaryCycle si nécessaire

	return created, nil
}

// deductVolume déduit les volumes utilisés des jambes
func (s *BinaryCommissionService) deductVolume(ctx context.Context, clientID primitive.ObjectID, legs *models.BinaryLegs, volumeUsed float64) (float64, float64, error) {
	leftRemaining := legs.LeftVolume - volumeUsed
	rightRemaining := legs.RightVolume - volumeUsed

	// S'assurer que les volumes ne deviennent pas négatifs
	if leftRemaining < 0 {
		leftRemaining = 0
	}
	if rightRemaining < 0 {
		rightRemaining = 0
	}

	// Mettre à jour dans la DB
	err := s.clientRepo.UpdateNetworkVolumes(ctx, clientID.Hex(), leftRemaining, rightRemaining)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to update network volumes: %w", err)
	}

	return leftRemaining, rightRemaining, nil
}

func (s *BinaryCommissionService) getMinVolumePerLeg() float64 {
	if s.config.MinVolumePerLeg <= 0 {
		return 1.0
	}
	return s.config.MinVolumePerLeg
}

func (s *BinaryCommissionService) calculateAmount(volumeUsed float64) float64 {
	amount := volumeUsed * s.config.CommissionRate
	return math.Round(amount*100) / 100
}

// updateClientEarnings met à jour les gains totaux et le wallet du client
func (s *BinaryCommissionService) updateClientEarnings(ctx context.Context, clientID string, amount float64) error {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return err
	}

	newTotalEarnings := client.TotalEarnings + amount
	newWalletBalance := client.WalletBalance + amount

	return s.clientRepo.UpdateEarnings(ctx, clientID, newTotalEarnings, newWalletBalance)
}

// GetLegsVolumes récupère les volumes et actifs des jambes gauche et droite (méthode publique)
func (s *BinaryCommissionService) GetLegsVolumes(ctx context.Context, client *models.Client) (*models.BinaryLegs, error) {
	return s.getLegsVolumes(ctx, client)
}

// GetLegsVolumesWithCache récupère les volumes et actifs avec un cache d'activité (version optimisée)
func (s *BinaryCommissionService) GetLegsVolumesWithCache(ctx context.Context, client *models.Client, activeCache map[string]bool, maxDepth int) (*models.BinaryLegs, error) {
	legs := &models.BinaryLegs{
		LeftVolume:  client.NetworkVolumeLeft,
		RightVolume: client.NetworkVolumeRight,
	}

	// Compter les actifs dans chaque jambe avec cache et limite de profondeur
	leftActives, err := s.countActivesInLegWithCache(ctx, client.LeftChildID, "left", activeCache, maxDepth, 0)
	if err != nil {
		return nil, err
	}
	legs.LeftActives = leftActives

	rightActives, err := s.countActivesInLegWithCache(ctx, client.RightChildID, "right", activeCache, maxDepth, 0)
	if err != nil {
		return nil, err
	}
	legs.RightActives = rightActives

	return legs, nil
}

// countActivesInLegWithCache compte récursivement les membres actifs avec cache et limite de profondeur
func (s *BinaryCommissionService) countActivesInLegWithCache(ctx context.Context, rootID *primitive.ObjectID, side string, activeCache map[string]bool, maxDepth int, currentDepth int) (int, error) {
	if rootID == nil {
		return 0, nil
	}

	// Limiter la profondeur pour éviter les calculs trop coûteux
	if maxDepth > 0 && currentDepth >= maxDepth {
		return 0, nil
	}

	count := 0
	visited := make(map[string]bool) // Pour éviter de compter deux fois
	queue := []struct {
		id    *primitive.ObjectID
		depth int
	}{{rootID, currentDepth}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Limiter la profondeur
		if maxDepth > 0 && current.depth >= maxDepth {
			continue
		}

		// Éviter les doublons
		if visited[current.id.Hex()] {
			continue
		}
		visited[current.id.Hex()] = true

		client, err := s.clientRepo.GetByID(ctx, current.id.Hex())
		if err != nil {
			continue // Ignorer les erreurs et continuer
		}

		// Vérifier si ce client est actif (utiliser le cache si disponible)
		var isActive bool
		if activeCache != nil {
			if cached, found := activeCache[current.id.Hex()]; found {
				isActive = cached
			} else {
				// Si pas dans le cache, vérifier et mettre en cache
				active, err := s.isClientActive(ctx, current.id.Hex())
				if err == nil {
					isActive = active
					activeCache[current.id.Hex()] = active
				}
			}
		} else {
			active, err := s.isClientActive(ctx, current.id.Hex())
			if err == nil {
				isActive = active
			}
		}

		if isActive {
			count++
		}

		// Ajouter tous les enfants à la queue (pour compter récursivement)
		nextDepth := current.depth + 1
		if maxDepth == 0 || nextDepth < maxDepth {
			if client.LeftChildID != nil {
				queue = append(queue, struct {
					id    *primitive.ObjectID
					depth int
				}{client.LeftChildID, nextDepth})
			}
			if client.RightChildID != nil {
				queue = append(queue, struct {
					id    *primitive.ObjectID
					depth int
				}{client.RightChildID, nextDepth})
			}
		}
	}

	return count, nil
}

// CheckQualification vérifie si un membre est qualifié (méthode publique)
func (s *BinaryCommissionService) CheckQualification(ctx context.Context, client *models.Client) (*models.BinaryQualification, error) {
	return s.checkQualification(ctx, client)
}

// IsClientActive vérifie si un client est actif (méthode publique)
func (s *BinaryCommissionService) IsClientActive(ctx context.Context, clientID string) (bool, error) {
	return s.isClientActive(ctx, clientID)
}

// GetOrCreateCapping récupère ou crée un enregistrement de capping (méthode publique)
func (s *BinaryCommissionService) GetOrCreateCapping(ctx context.Context, clientID primitive.ObjectID, date time.Time) (*models.BinaryCapping, error) {
	return s.getOrCreateCapping(ctx, clientID, date)
}

// CalculateCycles calcule le nombre de cycles possibles (méthode publique)
func (s *BinaryCommissionService) CalculateCycles(legs *models.BinaryLegs) int {
	return s.calculateCycles(legs)
}
