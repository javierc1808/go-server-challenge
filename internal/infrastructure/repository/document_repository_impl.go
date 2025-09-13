package repository

import (
	"context"
	"math/rand"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/domain/repository"

	"github.com/brianvoe/gofakeit/v5"
)

// DocumentRepositoryImpl implements DocumentRepository
type DocumentRepositoryImpl struct {
	cache repository.CacheRepository
	// In a real implementation, this would hold a database connection
	// For now we simulate with in-memory data
}

// NewDocumentRepositoryImpl creates a new DocumentRepositoryImpl instance
func NewDocumentRepositoryImpl(cache repository.CacheRepository) repository.DocumentRepository {
	return &DocumentRepositoryImpl{
		cache: cache,
	}
}

// GetAll returns all documents (simulated)
func (r *DocumentRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Document, error) {
	// First try to fetch from cache
	cachedDocs, err := r.cache.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// If there are cached documents, return them
	if len(cachedDocs) > 0 {
		return cachedDocs, nil
	}

	// If not cached, generate some simulated ones
	var documents []*entity.Document
	count := 1 + rand.Intn(20)
	for i := 0; i < count; i++ {
		doc := r.generateRandomDocument()
		documents = append(documents, doc)

		// Store in cache
		r.cache.Set(ctx, doc.ID, doc)
	}

	return documents, nil
}

// GetByID returns a document by ID (simulated)
func (r *DocumentRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.Document, error) {
	// First try to fetch from cache
	cachedDoc, err := r.cache.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if cachedDoc != nil {
		return cachedDoc, nil
	}

	// If not cached, generate a random one
	doc := r.generateRandomDocument()
	doc.ID = id // Usar el ID solicitado

	// Store in cache
	r.cache.Set(ctx, id, doc)

	return doc, nil
}

// Create creates a new document
func (r *DocumentRepositoryImpl) Create(ctx context.Context, document *entity.Document) error {
	// Store in cache
	return r.cache.Set(ctx, document.ID, document)
}

// Update updates an existing document (simulated)
func (r *DocumentRepositoryImpl) Update(ctx context.Context, document *entity.Document) error {
	// Update in cache
	return r.cache.Set(ctx, document.ID, document)
}

// Delete removes a document (simulated)
func (r *DocumentRepositoryImpl) Delete(ctx context.Context, id string) error {
	// Remove from cache
	return r.cache.Delete(ctx, id)
}

// generateRandomDocument generates a random document for simulation
func (r *DocumentRepositoryImpl) generateRandomDocument() *entity.Document {
	doc := entity.NewDocument(
		gofakeit.UUID(),
		gofakeit.BeerName(),
		gofakeit.AppVersion(),
	)

	// Add random attachments
	attachmentCount := 1 + rand.Intn(4)
	for j := 0; j < attachmentCount; j++ {
		doc.AddAttachment(gofakeit.BeerStyle())
	}

	// Add random contributors
	contributorCount := 1 + rand.Intn(4)
	for j := 0; j < contributorCount; j++ {
		user := entity.NewUser(gofakeit.UUID(), gofakeit.Name())
		doc.AddContributor(*user)
	}

	return doc
}
