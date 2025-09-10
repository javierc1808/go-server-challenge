package repository

import (
	"context"
	"math/rand"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/domain/repository"

	"github.com/brianvoe/gofakeit/v5"
)

// DocumentRepositoryImpl implementa DocumentRepository
type DocumentRepositoryImpl struct {
	cache repository.CacheRepository
	// En una implementación real, aquí tendríamos una conexión a base de datos
	// Por ahora simulamos con datos en memoria
}

// NewDocumentRepositoryImpl crea una nueva instancia de DocumentRepositoryImpl
func NewDocumentRepositoryImpl(cache repository.CacheRepository) repository.DocumentRepository {
	return &DocumentRepositoryImpl{
		cache: cache,
	}
}

// GetAll obtiene todos los documentos (simulado)
func (r *DocumentRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Document, error) {
	// Primero intentar obtener del cache
	cachedDocs, err := r.cache.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Si hay documentos en cache, devolverlos
	if len(cachedDocs) > 0 {
		return cachedDocs, nil
	}

	// Si no hay documentos en cache, generar algunos simulados
	var documents []*entity.Document
	count := 1 + rand.Intn(20)
	for i := 0; i < count; i++ {
		doc := r.generateRandomDocument()
		documents = append(documents, doc)

		// Almacenar en cache
		r.cache.Set(ctx, doc.ID, doc)
	}

	return documents, nil
}

// GetByID obtiene un documento por su ID (simulado)
func (r *DocumentRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.Document, error) {
	// Primero intentar obtener del cache
	cachedDoc, err := r.cache.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if cachedDoc != nil {
		return cachedDoc, nil
	}

	// Si no está en cache, generar uno aleatorio
	doc := r.generateRandomDocument()
	doc.ID = id // Usar el ID solicitado

	// Almacenar en cache
	r.cache.Set(ctx, id, doc)

	return doc, nil
}

// Create crea un nuevo documento
func (r *DocumentRepositoryImpl) Create(ctx context.Context, document *entity.Document) error {
	// Almacenar en cache
	return r.cache.Set(ctx, document.ID, document)
}

// Update actualiza un documento existente (simulado)
func (r *DocumentRepositoryImpl) Update(ctx context.Context, document *entity.Document) error {
	// Actualizar en cache
	return r.cache.Set(ctx, document.ID, document)
}

// Delete elimina un documento (simulado)
func (r *DocumentRepositoryImpl) Delete(ctx context.Context, id string) error {
	// Eliminar del cache
	return r.cache.Delete(ctx, id)
}

// generateRandomDocument genera un documento aleatorio para simulación
func (r *DocumentRepositoryImpl) generateRandomDocument() *entity.Document {
	doc := entity.NewDocument(
		gofakeit.UUID(),
		gofakeit.BeerName(),
		gofakeit.AppVersion(),
	)

	// Agregar adjuntos aleatorios
	attachmentCount := 1 + rand.Intn(4)
	for j := 0; j < attachmentCount; j++ {
		doc.AddAttachment(gofakeit.BeerStyle())
	}

	// Agregar contribuidores aleatorios
	contributorCount := 1 + rand.Intn(4)
	for j := 0; j < contributorCount; j++ {
		user := entity.NewUser(gofakeit.UUID(), gofakeit.Name())
		doc.AddContributor(*user)
	}

	return doc
}
