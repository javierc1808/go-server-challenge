package repository

import (
	"context"
	"sync"
	"time"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/domain/repository"
)

// MemoryCache implementa CacheRepository usando memoria
type MemoryCache struct {
	documents map[string]*entity.Document
	mutex     sync.RWMutex
	ttl       time.Duration
	expiry    map[string]time.Time
}

// NewMemoryCache crea una nueva instancia de MemoryCache
func NewMemoryCache(ttl time.Duration) repository.CacheRepository {
	cache := &MemoryCache{
		documents: make(map[string]*entity.Document),
		ttl:       ttl,
		expiry:    make(map[string]time.Time),
	}

	// Iniciar limpieza automática de elementos expirados
	go cache.startCleanup()

	return cache
}

// Set almacena un documento en el cache
func (c *MemoryCache) Set(ctx context.Context, key string, document *entity.Document) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.documents[key] = document
	c.expiry[key] = time.Now().Add(c.ttl)

	return nil
}

// Get obtiene un documento del cache
func (c *MemoryCache) Get(ctx context.Context, key string) (*entity.Document, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	document, exists := c.documents[key]
	if !exists {
		return nil, nil
	}

	// Verificar si ha expirado
	if time.Now().After(c.expiry[key]) {
		return nil, nil
	}

	return document, nil
}

// GetAll obtiene todos los documentos del cache
func (c *MemoryCache) GetAll(ctx context.Context) ([]*entity.Document, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var documents []*entity.Document
	now := time.Now()

	for key, document := range c.documents {
		// Solo incluir documentos no expirados
		if now.Before(c.expiry[key]) {
			documents = append(documents, document)
		}
	}

	return documents, nil
}

// Delete elimina un documento del cache
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.documents, key)
	delete(c.expiry, key)

	return nil
}

// Clear limpia todo el cache
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.documents = make(map[string]*entity.Document)
	c.expiry = make(map[string]time.Time)

	return nil
}

// Exists verifica si un documento existe en el cache
func (c *MemoryCache) Exists(ctx context.Context, key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	_, exists := c.documents[key]
	if !exists {
		return false
	}

	// Verificar si ha expirado
	return time.Now().Before(c.expiry[key])
}

// Count retorna el número de documentos en el cache
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

// startCleanup inicia la limpieza automática de elementos expirados
func (c *MemoryCache) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanupExpired()
	}
}

// cleanupExpired elimina elementos expirados del cache
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

// GetStats retorna estadísticas del cache
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
