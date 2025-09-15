package entity

import "time"

// User represents a user in the domain
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewUser creates a new instance of User
func NewUser(id, name string) *User {
	now := time.Now()
	return &User{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate validates the user's data
func (u *User) Validate() error {
	if u.ID == "" {
		return ErrInvalidUserID
	}
	if u.Name == "" {
		return ErrInvalidUserName
	}
	return nil
}
