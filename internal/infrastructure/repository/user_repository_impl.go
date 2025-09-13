package repository

import (
	"context"
	"math/rand"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/domain/repository"

	"github.com/brianvoe/gofakeit/v5"
)

// UserRepositoryImpl implements UserRepository
type UserRepositoryImpl struct {
	// In a real implementation, here we would have a database connection
}

// NewUserRepositoryImpl creates a new instance of UserRepositoryImpl
func NewUserRepositoryImpl() repository.UserRepository {
	return &UserRepositoryImpl{}
}

// GetAll gets all users (simulated)
func (r *UserRepositoryImpl) GetAll(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User

	// Generate simulated users
	count := 1 + rand.Intn(10)
	for i := 0; i < count; i++ {
		user := entity.NewUser(gofakeit.UUID(), gofakeit.Name())
		users = append(users, user)
	}

	return users, nil
}

// GetByID gets a user by its ID (simulated)
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.User, error) {
	// In a real implementation, we would search in the database
	// For now, we generate a random user
	return entity.NewUser(gofakeit.UUID(), gofakeit.Name()), nil
}

// Create creates a new user (simulated)
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	// In a real implementation, we would insert into the database
	return nil
}

// Update updates an existing user (simulated)
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	// In a real implementation, we would update in the database
	return nil
}

// Delete deletes a user (simulated)
func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	// In a real implementation, we would delete from the database
	return nil
}
