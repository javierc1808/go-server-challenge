package repository

import (
	"context"

	"frontend-challenge/internal/domain/entity"
)

// CacheRepository define la interfaz para el cache de documentos
type CacheRepository interface {
	// Set almacena un documento en el cache
	Set(ctx context.Context, key string, document *entity.Document) error

	// Get obtiene un documento del cache
	Get(ctx context.Context, key string) (*entity.Document, error)

	// GetAll obtiene todos los documentos del cache
	GetAll(ctx context.Context) ([]*entity.Document, error)

	// Delete elimina un documento del cache
	Delete(ctx context.Context, key string) error

	// Clear limpia todo el cache
	Clear(ctx context.Context) error

	// Exists verifica si un documento existe en el cache
	Exists(ctx context.Context, key string) bool

	// Count retorna el número de documentos en el cache
	Count(ctx context.Context) int

	// GetStats retorna estadísticas del cache
	GetStats() map[string]interface{}
}
