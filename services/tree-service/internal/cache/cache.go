package cache

import (
	"context"
	"fmt"
	"time"

	"bureau/services/tree-service/internal/models"

	"go.uber.org/zap"
)

// TreeCache définit l'interface pour le cache
type TreeCache interface {
	Get(ctx context.Context, key string) (*models.ClientTreeResponse, error)
	Set(ctx context.Context, key string, value *models.ClientTreeResponse, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// MemoryCache implémentation en mémoire
type MemoryCache struct {
	data   map[string]cacheEntry
	logger *zap.Logger
}

type cacheEntry struct {
	value      *models.ClientTreeResponse
	expiration time.Time
}

func NewMemoryCache(logger *zap.Logger) *MemoryCache {
	cache := &MemoryCache{
		data:   make(map[string]cacheEntry),
		logger: logger,
	}
	// Nettoyer le cache toutes les minutes
	go cache.cleanup()
	return cache
}

func (c *MemoryCache) Get(ctx context.Context, key string) (*models.ClientTreeResponse, error) {
	entry, exists := c.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found")
	}

	if time.Now().After(entry.expiration) {
		delete(c.data, key)
		return nil, fmt.Errorf("key expired")
	}

	return entry.value, nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, value *models.ClientTreeResponse, ttl time.Duration) error {
	c.data[key] = cacheEntry{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	delete(c.data, key)
	return nil
}

func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for key, entry := range c.data {
			if now.After(entry.expiration) {
				delete(c.data, key)
			}
		}
	}
}

// RedisCache implémentation Redis (à implémenter si Redis est disponible)
type RedisCache struct {
	// TODO: Implémenter avec go-redis
	logger *zap.Logger
}

func NewRedisCache(redisURL string, logger *zap.Logger) *RedisCache {
	// TODO: Initialiser la connexion Redis
	return &RedisCache{
		logger: logger,
	}
}

func (c *RedisCache) Get(ctx context.Context, key string) (*models.ClientTreeResponse, error) {
	// TODO: Implémenter avec go-redis
	return nil, fmt.Errorf("Redis not implemented yet")
}

func (c *RedisCache) Set(ctx context.Context, key string, value *models.ClientTreeResponse, ttl time.Duration) error {
	// TODO: Implémenter avec go-redis
	return fmt.Errorf("Redis not implemented yet")
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	// TODO: Implémenter avec go-redis
	return fmt.Errorf("Redis not implemented yet")
}


