package sql

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"traci/backend/internal/domain"
	"traci/backend/internal/gen/db"
)

type CollectionRepository struct{ queries db.Querier }

var _ domain.CollectionRepository = (*CollectionRepository)(nil)

func NewCollectionRepository(queries db.Querier) *CollectionRepository {
	return &CollectionRepository{queries: queries}
}

func (r *CollectionRepository) Create(ctx context.Context, item *domain.CollectionItem) error {
	rating, err := databaseRating(item.PersonalRating)
	if err != nil {
		return fmt.Errorf("create collection item: %w", err)
	}
	err = r.queries.CreateCollectionItem(ctx, db.CreateCollectionItemParams{
		ID: databaseUUID(item.ID), UserID: databaseUUID(item.UserID), MovieID: databaseUUID(item.MovieID),
		PersonalRating: rating, Review: databaseText(item.Review), Status: db.MovieStatus(item.Status),
		CreatedAt: databaseTime(item.CreatedAt), UpdatedAt: databaseTime(item.UpdatedAt),
	})
	return repositoryError("create collection item", err)
}

func (r *CollectionRepository) Get(ctx context.Context, userID, movieID uuid.UUID) (*domain.CollectionItem, error) {
	row, err := r.queries.GetCollectionItem(ctx, db.GetCollectionItemParams{
		UserID: databaseUUID(userID), MovieID: databaseUUID(movieID),
	})
	if err != nil {
		return nil, repositoryError("get collection item", err)
	}
	item := collectionFromRow(row)
	return &item, nil
}

func (r *CollectionRepository) List(ctx context.Context, userID uuid.UUID, filter domain.CollectionFilter) ([]domain.CollectionItem, int, error) {
	l, o, err := pageParams(filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list collection items: %w", err)
	}
	status := db.NullMovieStatus{}
	if filter.Status != nil {
		status = db.NullMovieStatus{MovieStatus: db.MovieStatus(*filter.Status), Valid: true}
	}
	rows, err := r.queries.ListCollectionItems(ctx, db.ListCollectionItemsParams{
		UserID: databaseUUID(userID), Status: status, PageLimit: l, PageOffset: o,
	})
	if err != nil {
		return nil, 0, repositoryError("list collection items", err)
	}
	total, err := r.queries.CountCollectionItems(ctx, db.CountCollectionItemsParams{UserID: databaseUUID(userID), Status: status})
	if err != nil {
		return nil, 0, repositoryError("count collection items", err)
	}
	items := make([]domain.CollectionItem, len(rows))
	for i, row := range rows {
		items[i] = collectionFromRow(db.GetCollectionItemRow(row))
	}
	return items, int(total), nil
}

func (r *CollectionRepository) Update(ctx context.Context, item *domain.CollectionItem) error {
	rating, err := databaseRating(item.PersonalRating)
	if err != nil {
		return fmt.Errorf("update collection item: %w", err)
	}
	rows, err := r.queries.UpdateCollectionItem(ctx, db.UpdateCollectionItemParams{
		UserID: databaseUUID(item.UserID), MovieID: databaseUUID(item.MovieID),
		PersonalRating: rating, Review: databaseText(item.Review), Status: db.MovieStatus(item.Status),
	})
	return affectedRows("update collection item", rows, err)
}

func (r *CollectionRepository) Delete(ctx context.Context, userID, movieID uuid.UUID) error {
	rows, err := r.queries.DeleteCollectionItem(ctx, db.DeleteCollectionItemParams{
		UserID: databaseUUID(userID), MovieID: databaseUUID(movieID),
	})
	return affectedRows("delete collection item", rows, err)
}

func collectionFromRow(row db.GetCollectionItemRow) domain.CollectionItem {
	return domain.CollectionItem{ID: uuid.UUID(row.ID.Bytes), UserID: uuid.UUID(row.UserID.Bytes),
		MovieID: uuid.UUID(row.MovieID.Bytes), PersonalRating: ratingPointer(row.PersonalRating),
		Review: textPointer(row.Review), Status: domain.MovieStatus(row.Status),
		CreatedAt: row.CreatedAt.Time, UpdatedAt: row.UpdatedAt.Time}
}
