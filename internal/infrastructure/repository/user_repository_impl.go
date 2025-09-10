package repository

import (
	"context"
	"math/rand"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/domain/repository"

	"github.com/brianvoe/gofakeit/v5"
)

// UserRepositoryImpl implementa UserRepository
type UserRepositoryImpl struct {
	// En una implementación real, aquí tendríamos una conexión a base de datos
}

// NewUserRepositoryImpl crea una nueva instancia de UserRepositoryImpl
func NewUserRepositoryImpl() repository.UserRepository {
	return &UserRepositoryImpl{}
}

// GetAll obtiene todos los usuarios (simulado)
func (r *UserRepositoryImpl) GetAll(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User

	// Generar usuarios simulados
	count := 1 + rand.Intn(10)
	for i := 0; i < count; i++ {
		user := entity.NewUser(gofakeit.UUID(), gofakeit.Name())
		users = append(users, user)
	}

	return users, nil
}

// GetByID obtiene un usuario por su ID (simulado)
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.User, error) {
	// En una implementación real, buscaríamos en la base de datos
	// Por ahora generamos un usuario aleatorio
	return entity.NewUser(gofakeit.UUID(), gofakeit.Name()), nil
}

// Create crea un nuevo usuario (simulado)
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	// En una implementación real, insertaríamos en la base de datos
	return nil
}

// Update actualiza un usuario existente (simulado)
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	// En una implementación real, actualizaríamos en la base de datos
	return nil
}

// Delete elimina un usuario (simulado)
func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	// En una implementación real, eliminaríamos de la base de datos
	return nil
}
