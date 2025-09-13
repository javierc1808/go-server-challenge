package repository

import (
	"context"

	"frontend-challenge/internal/domain/entity"
)

// DocumentRepository defines the interface for the document repository
type DocumentRepository interface {
	// GetAll retrieves all documents
	GetAll(ctx context.Context) ([]*entity.Document, error)

	// GetByID retrieves a document by its ID
	GetByID(ctx context.Context, id string) (*entity.Document, error)

	// Create creates a new document
	Create(ctx context.Context, document *entity.Document) error

	// Update updates an existing document
	Update(ctx context.Context, document *entity.Document) error

	// Delete deletes a document
	Delete(ctx context.Context, id string) error
}
