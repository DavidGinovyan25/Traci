package domain

import "github.com/google/uuid"

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
