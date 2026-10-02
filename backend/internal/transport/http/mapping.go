package httptransport

import (
	"github.com/oapi-codegen/nullable"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"traci/backend/internal/application"
	"traci/backend/internal/domain"
	"traci/backend/internal/gen"
)

func nullableValue[T any](value *T) nullable.Nullable[T] {
	if value == nil {
		return nullable.NewNullNullable[T]()
	}
	return nullable.NewNullableWithValue(*value)
}

func nullablePointer[T any](value nullable.Nullable[T]) *T {
	if value.IsNull() || !value.IsSpecified() {
		return nil
	}
	result := value.MustGet()
	return &result
}

func movieResponse(movie domain.Movie) gen.Movie {
	return gen.Movie{Id: movie.ID, Name: movie.Name, DurationMin: movie.DurationMin,
		Genre: gen.MovieGenre(movie.Genre), ReleaseYear: movie.ReleaseYear, Description: nullableValue(movie.Description),
		CreatedAt: movie.CreatedAt, UpdatedAt: movie.UpdatedAt}
}

func userResponse(user domain.User) gen.User {
	return gen.User{Id: user.ID, Role: gen.Role(user.Role), Username: user.Username,
		FirstName: user.FirstName, SecondName: user.SecondName, Email: openapi_types.Email(user.Email), IsBlocked: user.IsBlocked}
}

func collectionResponse(entry application.CollectionEntry) gen.CollectionItem {
	return gen.CollectionItem{Movie: movieResponse(entry.Movie), Status: gen.MovieStatus(entry.Item.Status),
		PersonalRating: nullableValue(entry.Item.PersonalRating), Review: nullableValue(entry.Item.Review),
		CreatedAt: entry.Item.CreatedAt, UpdatedAt: entry.Item.UpdatedAt}
}
