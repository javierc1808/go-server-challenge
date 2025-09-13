package repository

import (
	"context"

	"frontend-challenge/internal/domain/entity"
)

// UserRepository defines the interface for the user repository
type UserRepository interface {
	// GetAll retrieves all users
	GetAll(ctx context.Context) ([]*entity.User, error)

	// GetByID retrieves a user by its ID
	GetByID(ctx context.Context, id string) (*entity.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete deletes a user
	Delete(ctx context.Context, id string) error
}
