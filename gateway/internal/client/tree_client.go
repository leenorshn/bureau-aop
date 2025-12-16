package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"bureau/gateway/internal/models"

	"go.uber.org/zap"
)

type TreeServiceClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
}

func NewTreeServiceClient(baseURL string, logger *zap.Logger) *TreeServiceClient {
	return &TreeServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// GetClientTree récupère l'arbre client depuis le Tree Service
func (c *TreeServiceClient) GetClientTree(ctx context.Context, clientID string) (*models.ClientTreeResponse, error) {
	url := fmt.Sprintf("%s/api/v1/tree/%s", c.baseURL, clientID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call tree service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tree service returned error: %s", string(body))
	}

	var treeResponse models.ClientTreeResponse
	if err := json.NewDecoder(resp.Body).Decode(&treeResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &treeResponse, nil
}

