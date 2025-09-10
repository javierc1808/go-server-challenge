package repository

import (
	"context"

	"frontend-challenge/internal/domain/entity"
)

// UserRepository define la interfaz para el repositorio de usuarios
type UserRepository interface {
	// GetAll obtiene todos los usuarios
	GetAll(ctx context.Context) ([]*entity.User, error)

	// GetByID obtiene un usuario por su ID
	GetByID(ctx context.Context, id string) (*entity.User, error)

	// Create crea un nuevo usuario
	Create(ctx context.Context, user *entity.User) error

	// Update actualiza un usuario existente
	Update(ctx context.Context, user *entity.User) error

	// Delete elimina un usuario
	Delete(ctx context.Context, id string) error
}
