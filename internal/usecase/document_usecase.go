package usecase

import (
	"context"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/domain/repository"
)

// DocumentUsecase define los casos de uso para documentos
type DocumentUsecase struct {
	documentRepo repository.DocumentRepository
	userRepo     repository.UserRepository
}

// NewDocumentUsecase crea una nueva instancia de DocumentUsecase
func NewDocumentUsecase(documentRepo repository.DocumentRepository, userRepo repository.UserRepository) *DocumentUsecase {
	return &DocumentUsecase{
		documentRepo: documentRepo,
		userRepo:     userRepo,
	}
}

// GetAllDocuments obtiene todos los documentos
func (u *DocumentUsecase) GetAllDocuments(ctx context.Context) ([]*entity.Document, error) {
	return u.documentRepo.GetAll(ctx)
}

// GetDocumentByID obtiene un documento por su ID
func (u *DocumentUsecase) GetDocumentByID(ctx context.Context, id string) (*entity.Document, error) {
	if id == "" {
		return nil, entity.ErrInvalidDocumentID
	}
	return u.documentRepo.GetByID(ctx, id)
}

// CreateDocument crea un nuevo documento
func (u *DocumentUsecase) CreateDocument(ctx context.Context, document *entity.Document) error {
	if err := document.Validate(); err != nil {
		return err
	}
	return u.documentRepo.Create(ctx, document)
}

// UpdateDocument actualiza un documento existente
func (u *DocumentUsecase) UpdateDocument(ctx context.Context, document *entity.Document) error {
	if err := document.Validate(); err != nil {
		return err
	}
	return u.documentRepo.Update(ctx, document)
}

// DeleteDocument elimina un documento
func (u *DocumentUsecase) DeleteDocument(ctx context.Context, id string) error {
	if id == "" {
		return entity.ErrInvalidDocumentID
	}
	return u.documentRepo.Delete(ctx, id)
}
