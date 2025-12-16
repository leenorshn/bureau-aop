package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"bureau/services/tree-service/internal/service"

	"go.uber.org/zap"
)

type TreeHandler struct {
	treeService *service.TreeService
	logger      *zap.Logger
}

func NewTreeHandler(treeService *service.TreeService, logger *zap.Logger) *TreeHandler {
	return &TreeHandler{
		treeService: treeService,
		logger:      logger,
	}
}

// HandleTreeRequest gère les requêtes pour l'arbre client
func (h *TreeHandler) HandleTreeRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraire l'ID du client depuis l'URL
	// Format: /api/v1/tree/{clientId}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/tree/")
	clientID := strings.TrimSuffix(path, "/")

	if clientID == "" {
		http.Error(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	// Récupérer l'arbre
	tree, err := h.treeService.GetClientTree(r.Context(), clientID)
	if err != nil {
		h.logger.Error("Failed to get client tree", zap.Error(err), zap.String("clientId", clientID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retourner la réponse JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tree); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}


