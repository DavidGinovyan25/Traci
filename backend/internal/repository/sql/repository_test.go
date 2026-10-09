package sql

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"traci/backend/internal/domain"
	"traci/backend/internal/gen/db"
	"traci/backend/internal/testutil"
)

func TestRepositoryErrors(t *testing.T) {
	for _, tc := range []struct{ source, want error }{
		{pgx.ErrNoRows, domain.ErrNotFound},
		{&pgconn.PgError{Code: "23505"}, domain.ErrAlreadyExists},
		{&pgconn.PgError{Code: "23503"}, domain.ErrNotFound},
		{context.Canceled, context.Canceled},
	} {
		{
			err := repositoryError("operation", tc.source)
			require.ErrorIs(t, err, tc.want, "got %v, want %v", err, tc.want)
		}
	}
	require.ErrorIs(t, affectedRows("update", 0, nil), domain.ErrNotFound, "missing row not detected")
	require.NoError(t, affectedRows("update", 1, nil))
}

func TestRepositoriesPostgres(t *testing.T) {
	url := testutil.Postgres(t)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, url)
	require.NoError(t, err)
	defer conn.Close(ctx)
	queries := db.New(conn)
	users := NewUserRepository(queries)
	movies := NewMovieRepository(queries)
	collection := NewCollectionRepository(queries)
	user := domain.NewUser("user", "First", "Last", "user@example.com", "hash")
	require.NoError(t, users.Create(ctx, user))
	duplicate := domain.NewUser("other", "First", "Last", user.Email, "hash")
	{
		err := users.Create(ctx, duplicate)
		require.ErrorIs(t, err, domain.ErrAlreadyExists, "duplicate email: %v", err)
	}
	loaded, err := users.GetByEmail(ctx, user.Email)
	require.NoError(t, err, "user round trip: %v", err)
	assert.Equal(t, *user, *loaded, "user round trip: %v", err)
	user.IsBlocked = true
	require.NoError(t, users.Update(ctx, user))
	loaded, err = users.GetByID(ctx, user.ID)
	require.NoError(t, err, "user update: %v", err)
	require.True(t, loaded.IsBlocked, "user update: %v", err)
	userPage, total, err := users.List(ctx, 10, 0)
	require.NoError(t, err, "user list: %v", err)
	assert.EqualValues(t, 1, total, "user list: %v", err)
	require.Len(t, userPage, 1, "user list: %v", err)
	movie := domain.NewMovie("Film", 90, domain.MovieGenreDrama, 2020, nil)
	require.NoError(t, movies.Create(ctx, movie))
	foundMovie, err := movies.GetByID(ctx, movie.ID)
	require.NoError(t, err, "movie round trip: %v", err)
	require.Nil(t, foundMovie.Description, "movie round trip: %v", err)
	description := "description"
	movie.Description = &description
	require.NoError(t, movies.Update(ctx, movie))
	foundMovie, err = movies.GetByID(ctx, movie.ID)
	require.NoError(t, err, "movie description: %v", err)
	require.NotNil(t, foundMovie.Description, "movie description: %v", err)
	assert.Equal(t, description, *foundMovie.Description, "movie description: %v", err)
	genre := domain.MovieGenreDrama
	moviePage, total, err := movies.List(ctx, domain.MovieFilter{Genre: &genre, Search: "film", Limit: 10})
	require.NoError(t, err, "movie list: %v", err)
	assert.EqualValues(t, 1, total, "movie list: %v", err)
	require.Len(t, moviePage, 1, "movie list: %v", err)
	moviePage, total, err = movies.List(ctx, domain.MovieFilter{Limit: 10, Offset: 20})
	require.NoError(t, err, "empty page: %v", err)
	assert.EqualValues(t, 1, total, "empty page: %v", err)
	require.NotNil(t, moviePage, "empty page: %v", err)
	require.Len(t, moviePage, 0, "empty page: %v", err)
	item := domain.NewCollectionItem(user.ID, movie.ID, domain.MovieStatusPlanned)
	require.NoError(t, collection.Create(ctx, item))
	duplicateItem := domain.NewCollectionItem(user.ID, movie.ID, domain.MovieStatusPlanned)
	{
		err := collection.Create(ctx, duplicateItem)
		require.ErrorIs(t, err, domain.ErrAlreadyExists, "duplicate collection: %v", err)
	}
	missingMovie := domain.NewCollectionItem(user.ID, uuid.New(), domain.MovieStatusPlanned)
	{
		err := collection.Create(ctx, missingMovie)
		require.ErrorIs(t, err, domain.ErrNotFound, "foreign key: %v", err)
	}
	foundItem, err := collection.Get(ctx, user.ID, movie.ID)
	require.NoError(t, err, "nullable fields: %v", err)
	require.Nil(t, foundItem.PersonalRating, "nullable fields: %v", err)
	require.Nil(t, foundItem.Review, "nullable fields: %v", err)
	rating, review := 8, "review"
	item.PersonalRating, item.Review, item.Status = &rating, &review, domain.MovieStatusWatched
	require.NoError(t, collection.Update(ctx, item))
	foundItem, err = collection.Get(ctx, user.ID, movie.ID)
	require.NoError(t, err, "collection update: %v", err)
	require.NotNil(t, foundItem.PersonalRating, "collection update: %v", err)
	assert.Equal(t, rating, *foundItem.PersonalRating, "collection update: %v", err)
	require.NotNil(t, foundItem.Review, "collection update: %v", err)
	assert.Equal(t, review, *foundItem.Review, "collection update: %v", err)
	status := domain.MovieStatusWatched
	items, total, err := collection.List(ctx, user.ID, domain.CollectionFilter{Status: &status, Limit: 10})
	require.NoError(t, err, "collection list: %v", err)
	assert.EqualValues(t, 1, total, "collection list: %v", err)
	require.Len(t, items, 1, "collection list: %v", err)
	stranger := uuid.New()
	{
		err := collection.Delete(ctx, stranger, movie.ID)
		require.ErrorIs(t, err, domain.ErrNotFound, "owner scope: %v", err)
	}
	foreignItem := *item
	foreignItem.UserID = stranger
	{
		err := collection.Update(ctx, &foreignItem)
		require.ErrorIs(t, err, domain.ErrNotFound, "owner update scope: %v", err)
	}
	item.PersonalRating, item.Review = nil, nil
	require.NoError(t, collection.Update(ctx, item))
	foundItem, err = collection.Get(ctx, user.ID, movie.ID)
	require.NoError(t, err, "clear fields: %v", err)
	require.Nil(t, foundItem.PersonalRating, "clear fields: %v", err)
	require.Nil(t, foundItem.Review, "clear fields: %v", err)
	require.NoError(t, collection.Delete(ctx, user.ID, movie.ID))
	{
		_, err := collection.Get(ctx, user.ID, movie.ID)
		require.ErrorIs(t, err, domain.ErrNotFound, "deleted collection: %v", err)
	}
	require.NoError(t, collection.Create(ctx, item))
	require.NoError(t, movies.Delete(ctx, movie.ID))
	{
		_, err := collection.Get(ctx, user.ID, movie.ID)
		require.ErrorIs(t, err, domain.ErrNotFound, "cascade: %v", err)
	}
	{
		err := movies.Update(ctx, movie)
		require.ErrorIs(t, err, domain.ErrNotFound, "missing movie update: %v", err)
	}
	require.NoError(t, users.Delete(ctx, user.ID))
	{
		_, err := users.GetByID(ctx, user.ID)
		require.ErrorIs(t, err, domain.ErrNotFound, "deleted user: %v", err)
	}
}
