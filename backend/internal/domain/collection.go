package domain

import (
	"context"
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

type CollectionFilter struct {
	Status *MovieStatus
	Limit  int
	Offset int
}

type CollectionRepository interface {
	Create(ctx context.Context, item *CollectionItem) error
	Get(ctx context.Context, userID, movieID uuid.UUID) (*CollectionItem, error)
	List(ctx context.Context, userID uuid.UUID, filter CollectionFilter) (items []CollectionItem, total int, err error)
	Update(ctx context.Context, item *CollectionItem) error
	Delete(ctx context.Context, userID, movieID uuid.UUID) error
}
