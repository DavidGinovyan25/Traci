package domain

import (
	"time"

	"github.com/google/uuid"
)

type CollectionItem struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	MovieID        uuid.UUID
	PersonalRating *int
	Review         *string
	Status         MovieStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewCollectionItem(userID, movieID uuid.UUID, status MovieStatus) *CollectionItem {
	now := time.Now().UTC()
	return &CollectionItem{
		ID:        uuid.New(),
		UserID:    userID,
		MovieID:   movieID,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
