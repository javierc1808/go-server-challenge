package repository

import (
	"context"

	"frontend-challenge/internal/domain/entity"
)

// DocumentRepository define la interfaz para el repositorio de documentos
type DocumentRepository interface {
	// GetAll obtiene todos los documentos
	GetAll(ctx context.Context) ([]*entity.Document, error)

	// GetByID obtiene un documento por su ID
	GetByID(ctx context.Context, id string) (*entity.Document, error)

	// Create crea un nuevo documento
	Create(ctx context.Context, document *entity.Document) error

	// Update actualiza un documento existente
	Update(ctx context.Context, document *entity.Document) error

	// Delete elimina un documento
	Delete(ctx context.Context, id string) error
}
