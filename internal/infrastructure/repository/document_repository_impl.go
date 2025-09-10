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
	// En una implementación real, aquí tendríamos una conexión a base de datos
	// Por ahora simulamos con datos en memoria
}

// NewDocumentRepositoryImpl crea una nueva instancia de DocumentRepositoryImpl
func NewDocumentRepositoryImpl() repository.DocumentRepository {
	return &DocumentRepositoryImpl{}
}

// GetAll obtiene todos los documentos (simulado)
func (r *DocumentRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Document, error) {
	var documents []*entity.Document

	// Generar documentos simulados
	count := 1 + rand.Intn(20)
	for i := 0; i < count; i++ {
		doc := r.generateRandomDocument()
		documents = append(documents, doc)
	}

	return documents, nil
}

// GetByID obtiene un documento por su ID (simulado)
func (r *DocumentRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.Document, error) {
	// En una implementación real, buscaríamos en la base de datos
	// Por ahora generamos un documento aleatorio
	return r.generateRandomDocument(), nil
}

// Create crea un nuevo documento (simulado)
func (r *DocumentRepositoryImpl) Create(ctx context.Context, document *entity.Document) error {
	// En una implementación real, insertaríamos en la base de datos
	return nil
}

// Update actualiza un documento existente (simulado)
func (r *DocumentRepositoryImpl) Update(ctx context.Context, document *entity.Document) error {
	// En una implementación real, actualizaríamos en la base de datos
	return nil
}

// Delete elimina un documento (simulado)
func (r *DocumentRepositoryImpl) Delete(ctx context.Context, id string) error {
	// En una implementación real, eliminaríamos de la base de datos
	return nil
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
