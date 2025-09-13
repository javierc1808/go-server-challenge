package repository

import (
	"context"

	"frontend-challenge/internal/domain/entity"
)

// CacheRepository defines the interface for the document cache
type CacheRepository interface {
	// Set stores a document in the cache
	Set(ctx context.Context, key string, document *entity.Document) error

	// Get retrieves a document from the cache
	Get(ctx context.Context, key string) (*entity.Document, error)

	// GetAll retrieves all documents from the cache
	GetAll(ctx context.Context) ([]*entity.Document, error)

	// Delete removes a document from the cache
	Delete(ctx context.Context, key string) error

	// Clear clears the entire cache
	Clear(ctx context.Context) error

	// Exists checks if a document exists in the cache
	Exists(ctx context.Context, key string) bool

	// Count returns the number of documents in the cache
	Count(ctx context.Context) int

	// GetStats returns cache statistics
	GetStats() map[string]interface{}
}
