package domain

import (
	"context"
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

type MovieFilter struct {
	Genre  *MovieGenre
	Search string
	Limit  int
	Offset int
}

type MovieRepository interface {
	Create(ctx context.Context, movie *Movie) error
	GetByID(ctx context.Context, id uuid.UUID) (*Movie, error)
	List(ctx context.Context, filter MovieFilter) (movies []Movie, total int, err error)
	Update(ctx context.Context, movie *Movie) error
	Delete(ctx context.Context, id uuid.UUID) error
}
