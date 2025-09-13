package repository

import (
	"context"
	"sync"
	"time"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/domain/repository"
)

// MemoryCache implements CacheRepository using in-memory storage
type MemoryCache struct {
	documents map[string]*entity.Document
	mutex     sync.RWMutex
	ttl       time.Duration
	expiry    map[string]time.Time
}

// NewMemoryCache creates a new MemoryCache instance
func NewMemoryCache(ttl time.Duration) repository.CacheRepository {
	cache := &MemoryCache{
		documents: make(map[string]*entity.Document),
		ttl:       ttl,
		expiry:    make(map[string]time.Time),
	}

	// Start automatic cleanup of expired entries
	go cache.startCleanup()

	return cache
}

// Set stores a document in the cache
func (c *MemoryCache) Set(ctx context.Context, key string, document *entity.Document) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.documents[key] = document
	c.expiry[key] = time.Now().Add(c.ttl)

	return nil
}

// Get retrieves a document from the cache
func (c *MemoryCache) Get(ctx context.Context, key string) (*entity.Document, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	document, exists := c.documents[key]
	if !exists {
		return nil, nil
	}

	// Check if it has expired
	if time.Now().After(c.expiry[key]) {
		return nil, nil
	}

	return document, nil
}

// GetAll returns all documents from the cache
func (c *MemoryCache) GetAll(ctx context.Context) ([]*entity.Document, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var documents []*entity.Document
	now := time.Now()

	for key, document := range c.documents {
		// Include only non-expired documents
		if now.Before(c.expiry[key]) {
			documents = append(documents, document)
		}
	}

	return documents, nil
}

// Delete removes a document from the cache
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.documents, key)
	delete(c.expiry, key)

	return nil
}

// Clear clears the entire cache
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.documents = make(map[string]*entity.Document)
	c.expiry = make(map[string]time.Time)

	return nil
}

// Exists checks if a document exists in the cache
func (c *MemoryCache) Exists(ctx context.Context, key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	_, exists := c.documents[key]
	if !exists {
		return false
	}

	// Check if it has expired
	return time.Now().Before(c.expiry[key])
}

// Count returns the number of documents in the cache
func (c *MemoryCache) Count(ctx context.Context) int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	count := 0
	now := time.Now()

	for _, expiry := range c.expiry {
		if now.Before(expiry) {
			count++
		}
	}

	return count
}

// startCleanup starts the automatic cleanup of expired entries
func (c *MemoryCache) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanupExpired()
	}
}

// cleanupExpired removes expired entries from the cache
func (c *MemoryCache) cleanupExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, expiry := range c.expiry {
		if now.After(expiry) {
			delete(c.documents, key)
			delete(c.expiry, key)
		}
	}
}

// GetStats returns cache statistics
func (c *MemoryCache) GetStats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	now := time.Now()
	activeCount := 0
	expiredCount := 0

	for _, expiry := range c.expiry {
		if now.Before(expiry) {
			activeCount++
		} else {
			expiredCount++
		}
	}

	return map[string]interface{}{
		"total_documents":   len(c.documents),
		"active_documents":  activeCount,
		"expired_documents": expiredCount,
		"ttl_seconds":       c.ttl.Seconds(),
	}
}
