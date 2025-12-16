package service

import (
	"context"
	"fmt"
	"time"

	"bureau/services/tree-service/internal/cache"
	"bureau/services/tree-service/internal/models"
	"bureau/services/tree-service/internal/store"

	"go.uber.org/zap"
)

type TreeService struct {
	clientRepo *store.ClientRepository
	saleRepo   *store.SaleRepository
	cache      cache.TreeCache
	logger     *zap.Logger
}

func NewTreeService(
	clientRepo *store.ClientRepository,
	saleRepo *store.SaleRepository,
	cache cache.TreeCache,
	logger *zap.Logger,
) *TreeService {
	return &TreeService{
		clientRepo: clientRepo,
		saleRepo:   saleRepo,
		cache:      cache,
		logger:     logger,
	}
}

// GetClientTree récupère l'arbre client avec cache et optimisations
func (s *TreeService) GetClientTree(ctx context.Context, clientID string) (*models.ClientTreeResponse, error) {
	// Vérifier le cache d'abord
	cacheKey := fmt.Sprintf("tree:%s", clientID)
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		s.logger.Info("Cache hit for client tree", zap.String("clientId", clientID))
		return cached, nil
	}

	// Récupérer le client racine
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("client introuvable: %w", err)
	}

	// Cache pour les vérifications d'activité
	activeCache := make(map[string]bool)

	// Créer le nœud racine
	rootNode := s.buildTreeNode(client, 0, nil, activeCache, ctx)

	// Collecter tous les descendants
	nodes := []*models.TreeNode{rootNode}
	maxLevel := 0

	// Parcourir l'arbre de manière optimisée
	s.collectTreeNodes(ctx, client, 0, &nodes, &maxLevel, activeCache)

	response := &models.ClientTreeResponse{
		Root:       rootNode,
		Nodes:      nodes,
		TotalNodes: len(nodes),
		MaxLevel:   maxLevel,
	}

	// Mettre en cache pour 5 minutes
	if err := s.cache.Set(ctx, cacheKey, response, 5*time.Minute); err != nil {
		s.logger.Warn("Failed to cache tree", zap.Error(err))
	}

	return response, nil
}

// buildTreeNode construit un nœud de l'arbre avec les informations de base
func (s *TreeService) buildTreeNode(
	client *models.Client,
	level int,
	parentID *string,
	activeCache map[string]bool,
	ctx context.Context,
) *models.TreeNode {
	// Vérifier si le client est actif (avec cache)
	isActive := s.isClientActiveCached(ctx, client.ID, activeCache)

	node := &models.TreeNode{
		ID:                client.ID,
		ClientID:          client.ClientID,
		Name:              client.Name,
		Phone:             client.Phone,
		ParentID:          parentID,
		Level:             level,
		Position:          client.Position,
		NetworkVolumeLeft: client.NetworkVolumeLeft,
		NetworkVolumeRight: client.NetworkVolumeRight,
		BinaryPairs:       client.BinaryPairs,
		TotalEarnings:     client.TotalEarnings,
		WalletBalance:     client.WalletBalance,
		IsActive:          isActive,
		// Les champs calculés seront ajoutés seulement pour les 3 premiers niveaux
		LeftActives:  0,
		RightActives: 0,
		IsQualified:  false,
	}

	// Calculer les actifs seulement pour les 3 premiers niveaux
	if level < 3 {
		leftActives, rightActives := s.countActivesInLegs(ctx, client, activeCache, 3-level)
		node.LeftActives = leftActives
		node.RightActives = rightActives
		node.IsQualified = s.checkQualification(ctx, client, activeCache)
		
		// Calculer les cycles disponibles
		if leftActives > 0 && rightActives > 0 {
			if leftActives < rightActives {
				cycles := leftActives
				node.CyclesAvailable = &cycles
			} else {
				cycles := rightActives
				node.CyclesAvailable = &cycles
			}
		} else {
			zero := 0
			node.CyclesAvailable = &zero
		}
	} else {
		zero := 0
		node.CyclesAvailable = &zero
	}

	zero := 0
	node.CyclesPaidToday = &zero // Sera calculé par le Binary Commission Service si nécessaire

	return node
}

// collectTreeNodes collecte récursivement tous les nœuds de l'arbre
func (s *TreeService) collectTreeNodes(
	ctx context.Context,
	client *models.Client,
	level int,
	nodes *[]*models.TreeNode,
	maxLevel *int,
	activeCache map[string]bool,
) {
	currentLevel := level + 1
	if currentLevel > *maxLevel {
		*maxLevel = currentLevel
	}

		if client.LeftChildID != nil {
			leftChild, err := s.clientRepo.GetByID(ctx, *client.LeftChildID)
			if err == nil && leftChild != nil {
				parentID := client.ID
				position := "left"
				leftChild.Position = &position
				node := s.buildTreeNode(leftChild, currentLevel, &parentID, activeCache, ctx)
				*nodes = append(*nodes, node)
				s.collectTreeNodes(ctx, leftChild, currentLevel, nodes, maxLevel, activeCache)
			}
		}

		if client.RightChildID != nil {
			rightChild, err := s.clientRepo.GetByID(ctx, *client.RightChildID)
			if err == nil && rightChild != nil {
				parentID := client.ID
				position := "right"
				rightChild.Position = &position
				node := s.buildTreeNode(rightChild, currentLevel, &parentID, activeCache, ctx)
				*nodes = append(*nodes, node)
				s.collectTreeNodes(ctx, rightChild, currentLevel, nodes, maxLevel, activeCache)
			}
		}
}

// isClientActiveCached vérifie si un client est actif avec cache
func (s *TreeService) isClientActiveCached(ctx context.Context, clientID string, activeCache map[string]bool) bool {
	if active, found := activeCache[clientID]; found {
		return active
	}

	sales, err := s.saleRepo.GetByClientID(ctx, clientID)
	isActive := err == nil && len(sales) > 0
	activeCache[clientID] = isActive
	return isActive
}

// countActivesInLegs compte les actifs dans chaque jambe avec limite de profondeur
func (s *TreeService) countActivesInLegs(
	ctx context.Context,
	client *models.Client,
	activeCache map[string]bool,
	maxDepth int,
) (leftActives, rightActives int) {
	if maxDepth <= 0 {
		return 0, 0
	}

	if client.LeftChildID != nil {
		leftActives = s.countActivesInLeg(ctx, client.LeftChildID, activeCache, maxDepth-1, 0)
	}
	if client.RightChildID != nil {
		rightActives = s.countActivesInLeg(ctx, client.RightChildID, activeCache, maxDepth-1, 0)
	}

	return
}

// countActivesInLeg compte récursivement les actifs dans une jambe
func (s *TreeService) countActivesInLeg(
	ctx context.Context,
	rootIDStr *string,
	activeCache map[string]bool,
	maxDepth int,
	currentDepth int,
) int {
	if rootIDStr == nil || *rootIDStr == "" || (maxDepth > 0 && currentDepth >= maxDepth) {
		return 0
	}

	count := 0
	visited := make(map[string]bool)
	queue := []struct {
		id    string
		depth int
	}{{*rootIDStr, currentDepth}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if maxDepth > 0 && current.depth >= maxDepth {
			continue
		}

		if visited[current.id] {
			continue
		}
		visited[current.id] = true

		client, err := s.clientRepo.GetByID(ctx, current.id)
		if err != nil {
			continue
		}

		if s.isClientActiveCached(ctx, current.id, activeCache) {
			count++
		}

		nextDepth := current.depth + 1
		if maxDepth == 0 || nextDepth < maxDepth {
			if client.LeftChildID != nil {
				queue = append(queue, struct {
					id    string
					depth int
				}{*client.LeftChildID, nextDepth})
			}
			if client.RightChildID != nil {
				queue = append(queue, struct {
					id    string
					depth int
				}{*client.RightChildID, nextDepth})
			}
		}
	}

	return count
}

// checkQualification vérifie si un client est qualifié
func (s *TreeService) checkQualification(ctx context.Context, client *models.Client, activeCache map[string]bool) bool {
	hasDirectLeft := false
	hasDirectRight := false

	if client.LeftChildID != nil {
		hasDirectLeft = s.isClientActiveCached(ctx, *client.LeftChildID, activeCache)
	}
	if client.RightChildID != nil {
		hasDirectRight = s.isClientActiveCached(ctx, *client.RightChildID, activeCache)
	}

	return hasDirectLeft && hasDirectRight
}

// InvalidateCache invalide le cache pour un client
func (s *TreeService) InvalidateCache(ctx context.Context, clientID string) error {
	cacheKey := fmt.Sprintf("tree:%s", clientID)
	return s.cache.Delete(ctx, cacheKey)
}

