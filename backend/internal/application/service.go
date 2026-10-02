package application

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"traci/backend/internal/domain"
)

type Credentials interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
	Issue(userID uuid.UUID) (string, error)
	Parse(token string) (uuid.UUID, error)
}

type Service struct {
	users       domain.UserRepository
	movies      domain.MovieRepository
	collection  domain.CollectionRepository
	credentials Credentials
}

func NewService(users domain.UserRepository, movies domain.MovieRepository, collection domain.CollectionRepository, credentials Credentials) *Service {
	return &Service{users: users, movies: movies, collection: collection, credentials: credentials}
}

func (s *Service) Register(ctx context.Context, input Register) (*domain.User, error) {
	hash, err := s.credentials.Hash(input.Password)
	if err != nil {
		return nil, err
	}
	user := domain.NewUser(input.Username, input.FirstName, input.SecondName, strings.ToLower(input.Email), hash)
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, input Login) (string, error) {
	user, err := s.users.GetByEmail(ctx, strings.ToLower(input.Email))
	if errors.Is(err, domain.ErrNotFound) {
		return "", ErrUnauthorized
	}
	if err != nil {
		return "", err
	}
	if !s.credentials.Verify(input.Password, user.PasswordHash) {
		return "", ErrUnauthorized
	}
	if user.IsBlocked {
		return "", ErrBlocked
	}
	return s.credentials.Issue(user.ID)
}

func (s *Service) Authenticate(ctx context.Context, token string) (*domain.User, error) {
	id, err := s.credentials.Parse(token)
	if err != nil {
		return nil, ErrUnauthorized
	}
	user, err := s.users.GetByID(ctx, id)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}
	if user.IsBlocked {
		return nil, ErrBlocked
	}
	return user, nil
}

func requireAdmin(actor *domain.User) error {
	if actor == nil {
		return ErrUnauthorized
	}
	if actor.IsBlocked {
		return ErrBlocked
	}
	if actor.Role != domain.RoleAdmin {
		return ErrForbidden
	}
	return nil
}

func (s *Service) GetMovie(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	return s.movies.GetByID(ctx, id)
}

func (s *Service) ListMovies(ctx context.Context, filter domain.MovieFilter) (MoviePage, error) {
	items, total, err := s.movies.List(ctx, filter)
	return MoviePage{Items: items, Total: total, Limit: filter.Limit, Offset: filter.Offset}, err
}

func (s *Service) CreateMovie(ctx context.Context, actor *domain.User, movie *domain.Movie) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	return s.movies.Create(ctx, movie)
}

func (s *Service) UpdateMovie(ctx context.Context, actor *domain.User, movie *domain.Movie) (*domain.Movie, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	if err := s.movies.Update(ctx, movie); err != nil {
		return nil, err
	}
	return s.movies.GetByID(ctx, movie.ID)
}

func (s *Service) DeleteMovie(ctx context.Context, actor *domain.User, id uuid.UUID) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	return s.movies.Delete(ctx, id)
}

func (s *Service) GetCollectionItem(ctx context.Context, userID, movieID uuid.UUID) (CollectionEntry, error) {
	item, err := s.collection.Get(ctx, userID, movieID)
	if err != nil {
		return CollectionEntry{}, err
	}
	movie, err := s.movies.GetByID(ctx, movieID)
	if err != nil {
		return CollectionEntry{}, err
	}
	return CollectionEntry{Item: *item, Movie: *movie}, nil
}

func (s *Service) ListCollection(ctx context.Context, userID uuid.UUID, filter domain.CollectionFilter) (CollectionPage, error) {
	items, total, err := s.collection.List(ctx, userID, filter)
	if err != nil {
		return CollectionPage{}, err
	}
	entries := make([]CollectionEntry, 0, len(items))
	for _, item := range items {
		movie, err := s.movies.GetByID(ctx, item.MovieID)
		if err != nil {
			return CollectionPage{}, err
		}
		entries = append(entries, CollectionEntry{Item: item, Movie: *movie})
	}
	return CollectionPage{Items: entries, Total: total, Limit: filter.Limit, Offset: filter.Offset}, nil
}

func (s *Service) AddCollectionItem(ctx context.Context, userID, movieID uuid.UUID, status domain.MovieStatus) (CollectionEntry, error) {
	movie, err := s.movies.GetByID(ctx, movieID)
	if err != nil {
		return CollectionEntry{}, err
	}
	item := domain.NewCollectionItem(userID, movieID, status)
	if err := s.collection.Create(ctx, item); err != nil {
		return CollectionEntry{}, err
	}
	return CollectionEntry{Item: *item, Movie: *movie}, nil
}

func (s *Service) UpdateCollectionItem(ctx context.Context, item *domain.CollectionItem) (CollectionEntry, error) {
	if err := s.collection.Update(ctx, item); err != nil {
		return CollectionEntry{}, err
	}
	return s.GetCollectionItem(ctx, item.UserID, item.MovieID)
}

func (s *Service) DeleteCollectionItem(ctx context.Context, userID, movieID uuid.UUID) error {
	return s.collection.Delete(ctx, userID, movieID)
}

func (s *Service) ListUsers(ctx context.Context, actor *domain.User, limit, offset int) (UserPage, error) {
	if err := requireAdmin(actor); err != nil {
		return UserPage{}, err
	}
	users, total, err := s.users.List(ctx, limit, offset)
	return UserPage{Items: users, Total: total, Limit: limit, Offset: offset}, err
}

func (s *Service) GetUser(ctx context.Context, actor *domain.User, id uuid.UUID) (*domain.User, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	return s.users.GetByID(ctx, id)
}

func (s *Service) SetUserBlocked(ctx context.Context, actor *domain.User, id uuid.UUID, blocked bool) (*domain.User, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	if actor.ID == id {
		return nil, ErrSelfModification
	}
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.IsBlocked = blocked
	if err := s.users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, actor *domain.User, id uuid.UUID) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	if actor.ID == id {
		return ErrSelfModification
	}
	return s.users.Delete(ctx, id)
}
