package sql

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"traci/backend/internal/domain"
	"traci/backend/internal/gen/db"
)

type UserRepository struct{ queries db.Querier }

var _ domain.UserRepository = (*UserRepository)(nil)

func NewUserRepository(queries db.Querier) *UserRepository { return &UserRepository{queries: queries} }

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	err := r.queries.CreateUser(ctx, db.CreateUserParams{
		ID: databaseUUID(user.ID), Role: db.Role(user.Role), Username: user.Username,
		FirstName: user.FirstName, SecondName: user.SecondName, Email: user.Email,
		PasswordHash: user.PasswordHash, IsBlocked: user.IsBlocked,
	})
	return repositoryError("create user", err)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row, err := r.queries.GetUserByID(ctx, databaseUUID(id))
	if err != nil {
		return nil, repositoryError("get user by id", err)
	}
	user := userFromRow(row)
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, repositoryError("get user by email", err)
	}
	user := userFromRow(row)
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]domain.User, int, error) {
	l, o, err := pageParams(limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	rows, err := r.queries.ListUsers(ctx, db.ListUsersParams{PageLimit: l, PageOffset: o})
	if err != nil {
		return nil, 0, repositoryError("list users", err)
	}
	total, err := r.queries.CountUsers(ctx)
	if err != nil {
		return nil, 0, repositoryError("count users", err)
	}
	users := make([]domain.User, len(rows))
	for i, row := range rows {
		users[i] = userFromRow(row)
	}
	return users, int(total), nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	rows, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID: databaseUUID(user.ID), Role: db.Role(user.Role), Username: user.Username,
		FirstName: user.FirstName, SecondName: user.SecondName, Email: user.Email,
		PasswordHash: user.PasswordHash, IsBlocked: user.IsBlocked,
	})
	return affectedRows("update user", rows, err)
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	rows, err := r.queries.DeleteUser(ctx, databaseUUID(id))
	return affectedRows("delete user", rows, err)
}

func userFromRow(row db.User) domain.User {
	return domain.User{ID: uuid.UUID(row.ID.Bytes), Role: domain.Role(row.Role),
		Username: row.Username, FirstName: row.FirstName, SecondName: row.SecondName,
		Email: row.Email, PasswordHash: row.PasswordHash, IsBlocked: row.IsBlocked}
}
