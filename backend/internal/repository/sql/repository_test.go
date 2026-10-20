package sql

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	unique := &pgconn.PgError{Code: "23505"}
	foreign := &pgconn.PgError{Code: "23503"}
	check := &pgconn.PgError{Code: "23514"}
	for _, tc := range []struct {
		name         string
		source, want error
		postgres     *pgconn.PgError
	}{
		{"no rows", pgx.ErrNoRows, domain.ErrNotFound, nil},
		{"wrapped no rows", fmt.Errorf("query: %w", pgx.ErrNoRows), domain.ErrNotFound, nil},
		{"unique violation", unique, domain.ErrAlreadyExists, unique},
		{"wrapped unique violation", fmt.Errorf("query: %w", unique), domain.ErrAlreadyExists, unique},
		{"foreign key violation", foreign, domain.ErrNotFound, foreign},
		{"check violation", check, check, check},
		{"canceled context", context.Canceled, context.Canceled, nil},
		{"deadline exceeded", context.DeadlineExceeded, context.DeadlineExceeded, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := repositoryError("operation", tc.source)
			require.ErrorIs(t, err, tc.want)
			assert.Contains(t, err.Error(), "operation")
			if tc.postgres != nil {
				var postgresError *pgconn.PgError
				require.ErrorAs(t, err, &postgresError)
				assert.Same(t, tc.postgres, postgresError)
			}
		})
	}
	require.NoError(t, repositoryError("operation", nil))
}

func TestAffectedRows(t *testing.T) {
	for _, tc := range []struct {
		name         string
		rows         int64
		source, want error
	}{
		{"missing row", 0, nil, domain.ErrNotFound},
		{"updated row", 1, nil, nil},
		{"multiple rows", 2, nil, nil},
		{"query failure before row count", 0, context.Canceled, context.Canceled},
		{"unique violation", 0, &pgconn.PgError{Code: "23505"}, domain.ErrAlreadyExists},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := affectedRows("update", tc.rows, tc.source)
			if tc.want == nil {
				require.NoError(t, err)
				return
			}
			require.ErrorIs(t, err, tc.want)
			assert.Contains(t, err.Error(), "update")
		})
	}
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
	_, err = users.GetByEmail(ctx, "missing@example.com")
	require.ErrorIs(t, err, domain.ErrNotFound)
	duplicate := domain.NewUser("other", "First", "Last", user.Email, "hash")
	{
		err := users.Create(ctx, duplicate)
		require.ErrorIs(t, err, domain.ErrAlreadyExists, "duplicate email: %v", err)
	}
	loaded, err := users.GetByEmail(ctx, user.Email)
	require.NoError(t, err, "user round trip: %v", err)
	assert.Equal(t, *user, *loaded, "user round trip: %v", err)
	duplicateUsername := domain.NewUser(user.Username, "Other", "User", "other@example.com", "hash")
	require.ErrorIs(t, users.Create(ctx, duplicateUsername), domain.ErrAlreadyExists)
	user.IsBlocked = true
	require.NoError(t, users.Update(ctx, user))
	loaded, err = users.GetByID(ctx, user.ID)
	require.NoError(t, err, "user update: %v", err)
	require.True(t, loaded.IsBlocked, "user update: %v", err)
	userPage, total, err := users.List(ctx, 10, 0)
	require.NoError(t, err, "user list: %v", err)
	assert.EqualValues(t, 1, total, "user list: %v", err)
	require.Len(t, userPage, 1, "user list: %v", err)
	assert.Equal(t, *loaded, userPage[0])
	userPage, total, err = users.List(ctx, 10, 20)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.NotNil(t, userPage)
	assert.Empty(t, userPage)
	movie := domain.NewMovie("Film", 90, domain.MovieGenreDrama, 2020, nil)
	require.NoError(t, movies.Create(ctx, movie))
	foundMovie, err := movies.GetByID(ctx, movie.ID)
	require.NoError(t, err, "movie round trip: %v", err)
	require.Nil(t, foundMovie.Description, "movie round trip: %v", err)
	assert.Equal(t, movie.ID, foundMovie.ID)
	assert.Equal(t, movie.Name, foundMovie.Name)
	assert.WithinDuration(t, movie.CreatedAt, foundMovie.CreatedAt, time.Microsecond)
	assert.WithinDuration(t, movie.UpdatedAt, foundMovie.UpdatedAt, time.Microsecond)
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
	assert.Equal(t, movie.ID, moviePage[0].ID)
	moviePage, total, err = movies.List(ctx, domain.MovieFilter{Search: "missing title", Limit: 10})
	require.NoError(t, err)
	assert.Zero(t, total)
	assert.Empty(t, moviePage)
	movie.Description = nil
	require.NoError(t, movies.Update(ctx, movie))
	foundMovie, err = movies.GetByID(ctx, movie.ID)
	require.NoError(t, err)
	require.Nil(t, foundMovie.Description)
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
	assert.Equal(t, item.ID, foundItem.ID)
	assert.Equal(t, item.UserID, foundItem.UserID)
	assert.Equal(t, item.MovieID, foundItem.MovieID)
	assert.Equal(t, domain.MovieStatusPlanned, foundItem.Status)
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
	assert.Equal(t, item.ID, items[0].ID)
	planned := domain.MovieStatusPlanned
	items, total, err = collection.List(ctx, user.ID, domain.CollectionFilter{Status: &planned, Limit: 10})
	require.NoError(t, err)
	assert.Zero(t, total)
	assert.Empty(t, items)
	items, total, err = collection.List(ctx, user.ID, domain.CollectionFilter{Limit: 10, Offset: 20})
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.NotNil(t, items)
	assert.Empty(t, items)
	stranger := uuid.New()
	_, err = collection.Get(ctx, stranger, movie.ID)
	require.ErrorIs(t, err, domain.ErrNotFound)
	items, total, err = collection.List(ctx, stranger, domain.CollectionFilter{Limit: 10})
	require.NoError(t, err)
	assert.Zero(t, total)
	assert.Empty(t, items)
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
	require.ErrorIs(t, movies.Delete(ctx, movie.ID), domain.ErrNotFound)
	_, err = movies.GetByID(ctx, movie.ID)
	require.ErrorIs(t, err, domain.ErrNotFound)
	cascadeMovie := domain.NewMovie("Cascade", 100, domain.MovieGenreDrama, 2021, nil)
	require.NoError(t, movies.Create(ctx, cascadeMovie))
	cascadeItem := domain.NewCollectionItem(user.ID, cascadeMovie.ID, domain.MovieStatusPlanned)
	require.NoError(t, collection.Create(ctx, cascadeItem))
	require.NoError(t, users.Delete(ctx, user.ID))
	_, err = collection.Get(ctx, user.ID, cascadeMovie.ID)
	require.ErrorIs(t, err, domain.ErrNotFound)
	_, err = movies.GetByID(ctx, cascadeMovie.ID)
	require.NoError(t, err)
	require.ErrorIs(t, users.Update(ctx, user), domain.ErrNotFound)
	require.ErrorIs(t, users.Delete(ctx, user.ID), domain.ErrNotFound)
	{
		_, err := users.GetByID(ctx, user.ID)
		require.ErrorIs(t, err, domain.ErrNotFound, "deleted user: %v", err)
	}
}
