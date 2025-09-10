package entity

import "time"

// User representa un usuario en el dominio
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser crea una nueva instancia de User
func NewUser(id, name string) *User {
	now := time.Now()
	return &User{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate valida los datos del usuario
func (u *User) Validate() error {
	if u.ID == "" {
		return ErrInvalidUserID
	}
	if u.Name == "" {
		return ErrInvalidUserName
	}
	return nil
}
