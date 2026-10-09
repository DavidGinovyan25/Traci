package domain

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Role         Role
	Username     string
	FirstName    string
	SecondName   string
	Email        string
	PasswordHash string
	IsBlocked    bool
}

func NewUser(username, firstName, secondName, email, passwordHash string) *User {
	return &User{
		ID:           uuid.New(),
		Role:         RoleUser,
		Username:     username,
		FirstName:    firstName,
		SecondName:   secondName,
		Email:        email,
		PasswordHash: passwordHash,
		IsBlocked:    false,
	}
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, limit, offset int) (users []User, total int, err error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
