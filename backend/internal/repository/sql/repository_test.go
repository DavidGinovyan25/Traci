package sql

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"traci/backend/internal/domain"
	"traci/backend/internal/gen/db"
)

func TestRepositoryErrors(t *testing.T) {
	for _, tc := range []struct{ source, want error }{
		{pgx.ErrNoRows, domain.ErrNotFound},
		{&pgconn.PgError{Code: "23505"}, domain.ErrAlreadyExists},
		{&pgconn.PgError{Code: "23503"}, domain.ErrNotFound},
		{context.Canceled, context.Canceled},
	} {
		if err := repositoryError("operation", tc.source); !errors.Is(err, tc.want) {
			t.Fatalf("got %v, want %v", err, tc.want)
		}
	}
	if !errors.Is(affectedRows("update", 0, nil), domain.ErrNotFound) {
		t.Fatal("missing row not detected")
	}
	if err := affectedRows("update", 1, nil); err != nil {
		t.Fatal(err)
	}
}

func TestRepositoriesPostgres(t *testing.T) {
	url := os.Getenv("POSTGRES_TEST_URL")
	if url == "" {
		t.Skip("POSTGRES_TEST_URL is not set")
	}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(ctx)
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	schema := "test_" + uuid.New().String()
	name := pgx.Identifier{schema}.Sanitize()
	if _, err := tx.Exec(ctx, "CREATE SCHEMA "+name); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, "SET LOCAL search_path TO "+name); err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../../../db/postgres/migrations/20260930180010_init_db.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, string(migration)); err != nil {
		t.Fatal(err)
	}
	queries := db.New(tx)
	users := NewUserRepository(queries)
	movies := NewMovieRepository(queries)
	collection := NewCollectionRepository(queries)
	user := domain.NewUser("user", "First", "Last", "user@example.com", "hash")
	if err := users.Create(ctx, user); err != nil {
		t.Fatal(err)
	}
	loaded, err := users.GetByEmail(ctx, user.Email)
	if err != nil || *loaded != *user {
		t.Fatalf("user round trip: %v", err)
	}
	user.IsBlocked = true
	if err := users.Update(ctx, user); err != nil {
		t.Fatal(err)
	}
	loaded, err = users.GetByID(ctx, user.ID)
	if err != nil || !loaded.IsBlocked {
		t.Fatalf("user update: %v", err)
	}
	userPage, total, err := users.List(ctx, 10, 0)
	if err != nil || total != 1 || len(userPage) != 1 {
		t.Fatalf("user list: %v", err)
	}
	movie := domain.NewMovie("Film", 90, domain.MovieGenreDrama, 2020, nil)
	if err := movies.Create(ctx, movie); err != nil {
		t.Fatal(err)
	}
	foundMovie, err := movies.GetByID(ctx, movie.ID)
	if err != nil || foundMovie.Description != nil {
		t.Fatalf("movie round trip: %v", err)
	}
	description := "description"
	movie.Description = &description
	if err := movies.Update(ctx, movie); err != nil {
		t.Fatal(err)
	}
	foundMovie, err = movies.GetByID(ctx, movie.ID)
	if err != nil || foundMovie.Description == nil || *foundMovie.Description != description {
		t.Fatalf("movie description: %v", err)
	}
	genre := domain.MovieGenreDrama
	moviePage, total, err := movies.List(ctx, domain.MovieFilter{Genre: &genre, Search: "film", Limit: 10})
	if err != nil || total != 1 || len(moviePage) != 1 {
		t.Fatalf("movie list: %v", err)
	}
	moviePage, total, err = movies.List(ctx, domain.MovieFilter{Limit: 10, Offset: 20})
	if err != nil || total != 1 || moviePage == nil || len(moviePage) != 0 {
		t.Fatalf("empty page: %v", err)
	}
	item := domain.NewCollectionItem(user.ID, movie.ID, domain.MovieStatusPlanned)
	if err := collection.Create(ctx, item); err != nil {
		t.Fatal(err)
	}
	foundItem, err := collection.Get(ctx, user.ID, movie.ID)
	if err != nil || foundItem.PersonalRating != nil || foundItem.Review != nil {
		t.Fatalf("nullable fields: %v", err)
	}
	rating, review := 8, "review"
	item.PersonalRating, item.Review, item.Status = &rating, &review, domain.MovieStatusWatched
	if err := collection.Update(ctx, item); err != nil {
		t.Fatal(err)
	}
	foundItem, err = collection.Get(ctx, user.ID, movie.ID)
	if err != nil || foundItem.PersonalRating == nil || *foundItem.PersonalRating != rating || foundItem.Review == nil || *foundItem.Review != review {
		t.Fatalf("collection update: %v", err)
	}
	status := domain.MovieStatusWatched
	items, total, err := collection.List(ctx, user.ID, domain.CollectionFilter{Status: &status, Limit: 10})
	if err != nil || total != 1 || len(items) != 1 {
		t.Fatalf("collection list: %v", err)
	}
	stranger := uuid.New()
	if err := collection.Delete(ctx, stranger, movie.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("owner scope: %v", err)
	}
	foreignItem := *item
	foreignItem.UserID = stranger
	if err := collection.Update(ctx, &foreignItem); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("owner update scope: %v", err)
	}
	item.PersonalRating, item.Review = nil, nil
	if err := collection.Update(ctx, item); err != nil {
		t.Fatal(err)
	}
	foundItem, err = collection.Get(ctx, user.ID, movie.ID)
	if err != nil || foundItem.PersonalRating != nil || foundItem.Review != nil {
		t.Fatalf("clear fields: %v", err)
	}
	if err := collection.Delete(ctx, user.ID, movie.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := collection.Get(ctx, user.ID, movie.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("deleted collection: %v", err)
	}
	if err := collection.Create(ctx, item); err != nil {
		t.Fatal(err)
	}
	if err := movies.Delete(ctx, movie.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := collection.Get(ctx, user.ID, movie.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("cascade: %v", err)
	}
	if err := movies.Update(ctx, movie); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("missing movie update: %v", err)
	}
	if err := users.Delete(ctx, user.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := users.GetByID(ctx, user.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("deleted user: %v", err)
	}
}
