package domain

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID          uuid.UUID
	Name        string
	DurationMin int
	Genre       MovieGenre
	ReleaseYear int
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewMovie(name string, durationMin int, genre MovieGenre, releaseYear int, description *string) *Movie {
	now := time.Now().UTC()
	return &Movie{
		ID:          uuid.New(),
		Name:        name,
		DurationMin: durationMin,
		Genre:       genre,
		ReleaseYear: releaseYear,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
