package usecase

import (
	"context"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/domain/repository"
)

// DocumentUsecase defines the use cases for documents
type DocumentUsecase struct {
	documentRepo repository.DocumentRepository
	userRepo     repository.UserRepository
}

// NewDocumentUsecase creates a new instance of DocumentUsecase
func NewDocumentUsecase(documentRepo repository.DocumentRepository, userRepo repository.UserRepository) *DocumentUsecase {
	return &DocumentUsecase{
		documentRepo: documentRepo,
		userRepo:     userRepo,
	}
}

// GetAllDocuments retrieves all documents
func (u *DocumentUsecase) GetAllDocuments(ctx context.Context) ([]*entity.Document, error) {
	return u.documentRepo.GetAll(ctx)
}

// GetDocumentByID retrieves a document by its ID
func (u *DocumentUsecase) GetDocumentByID(ctx context.Context, id string) (*entity.Document, error) {
	if id == "" {
		return nil, entity.ErrInvalidDocumentID
	}
	return u.documentRepo.GetByID(ctx, id)
}

// CreateDocument creates a new document
func (u *DocumentUsecase) CreateDocument(ctx context.Context, document *entity.Document) error {
	if err := document.Validate(); err != nil {
		return err
	}
	return u.documentRepo.Create(ctx, document)
}

// UpdateDocument updates an existing document
func (u *DocumentUsecase) UpdateDocument(ctx context.Context, document *entity.Document) error {
	if err := document.Validate(); err != nil {
		return err
	}
	return u.documentRepo.Update(ctx, document)
}

// DeleteDocument deletes a document
func (u *DocumentUsecase) DeleteDocument(ctx context.Context, id string) error {
	if id == "" {
		return entity.ErrInvalidDocumentID
	}
	return u.documentRepo.Delete(ctx, id)
}
