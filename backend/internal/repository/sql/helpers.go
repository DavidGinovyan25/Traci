package sql

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"traci/backend/internal/domain"
)

func repositoryError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%s: %w", operation, domain.ErrNotFound)
	}
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) {
		switch postgresError.Code {
		case "23505":
			return fmt.Errorf("%s: %w: %w", operation, domain.ErrAlreadyExists, err)
		case "23503":
			return fmt.Errorf("%s: %w: %w", operation, domain.ErrNotFound, err)
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}

func affectedRows(operation string, rows int64, err error) error {
	if err != nil {
		return repositoryError(operation, err)
	}
	if rows == 0 {
		return fmt.Errorf("%s: %w", operation, domain.ErrNotFound)
	}
	return nil
}

func databaseInt(value int) (int32, error) {
	if value < math.MinInt32 || value > math.MaxInt32 {
		return 0, fmt.Errorf("integer %d is outside PostgreSQL integer range", value)
	}
	return int32(value), nil
}

func pageParams(limit, offset int) (int32, int32, error) {
	l, err := databaseInt(limit)
	if err != nil {
		return 0, 0, err
	}
	o, err := databaseInt(offset)
	if err != nil {
		return 0, 0, err
	}
	return l, o, nil
}

func databaseUUID(id uuid.UUID) pgtype.UUID { return pgtype.UUID{Bytes: id, Valid: true} }
func databaseTime(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: value, Valid: true}
}
func databaseText(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *value, Valid: true}
}
func textPointer(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
func databaseRating(value *int) (pgtype.Int4, error) {
	if value == nil {
		return pgtype.Int4{}, nil
	}
	n, err := databaseInt(*value)
	if err != nil {
		return pgtype.Int4{}, err
	}
	return pgtype.Int4{Int32: n, Valid: true}, nil
}
func ratingPointer(value pgtype.Int4) *int {
	if !value.Valid {
		return nil
	}
	n := int(value.Int32)
	return &n
}
