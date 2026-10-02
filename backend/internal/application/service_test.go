package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"traci/backend/internal/domain"
)

type userStub struct {
	domain.UserRepository
	create  func(context.Context, *domain.User) error
	byID    func(context.Context, uuid.UUID) (*domain.User, error)
	byEmail func(context.Context, string) (*domain.User, error)
}

func (s userStub) Create(ctx context.Context, user *domain.User) error { return s.create(ctx, user) }
func (s userStub) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.byID(ctx, id)
}
func (s userStub) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.byEmail(ctx, email)
}

type movieStub struct {
	domain.MovieRepository
	get func(context.Context, uuid.UUID) (*domain.Movie, error)
}

func (s movieStub) GetByID(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	return s.get(ctx, id)
}

type collectionStub struct {
	domain.CollectionRepository
	create func(context.Context, *domain.CollectionItem) error
	list   func(context.Context, uuid.UUID, domain.CollectionFilter) ([]domain.CollectionItem, int, error)
}

func (s collectionStub) Create(ctx context.Context, item *domain.CollectionItem) error {
	return s.create(ctx, item)
}
func (s collectionStub) List(ctx context.Context, id uuid.UUID, filter domain.CollectionFilter) ([]domain.CollectionItem, int, error) {
	return s.list(ctx, id, filter)
}

type credentialStub struct {
	hash   func(string) (string, error)
	verify func(string, string) bool
	issue  func(uuid.UUID) (string, error)
	parse  func(string) (uuid.UUID, error)
}

func (s credentialStub) Hash(password string) (string, error)  { return s.hash(password) }
func (s credentialStub) Verify(password, hash string) bool     { return s.verify(password, hash) }
func (s credentialStub) Issue(id uuid.UUID) (string, error)    { return s.issue(id) }
func (s credentialStub) Parse(token string) (uuid.UUID, error) { return s.parse(token) }

func TestRegisterStoresHashAndNormalizesEmail(t *testing.T) {
	input := Register{Username: "alice", FirstName: "Alice", SecondName: "User", Email: "ALICE@EXAMPLE.COM", Password: "password123"}
	called := false
	users := userStub{create: func(ctx context.Context, user *domain.User) error {
		called = true
		require.NotEqual(t, uuid.Nil, user.ID, "unexpected user: %+v", user)
		assert.Equal(t, domain.RoleUser, user.Role, "unexpected user: %+v", user)
		require.False(t, user.IsBlocked, "unexpected user: %+v", user)
		assert.Equal(t, "alice@example.com", user.Email, "unexpected user: %+v", user)
		assert.Equal(t, "stored-hash", user.PasswordHash, "unexpected user: %+v", user)
		return nil
	}}
	credentials := credentialStub{hash: func(password string) (string, error) {
		assert.Equal(t, input.Password, password, "wrong password sent to hasher")
		return "stored-hash", nil
	}}
	service := NewService(users, nil, nil, credentials)
	user, err := service.Register(context.Background(), input)
	require.NoError(t, err, "register: %v", err)
	require.NotNil(t, user, "register: %v", err)
	require.True(t, called, "register: %v", err)
}

func TestRegisterFailure(t *testing.T) {
	failure := errors.New("hasher unavailable")
	service := NewService(nil, nil, nil, credentialStub{hash: func(string) (string, error) { return "", failure }})
	{
		user, err := service.Register(context.Background(), Register{})
		require.Nil(t, user, "hash error: %v", err)
		require.ErrorIs(t, err, failure, "hash error: %v", err)
	}
	service = NewService(userStub{create: func(context.Context, *domain.User) error { return domain.ErrAlreadyExists }}, nil, nil,
		credentialStub{hash: func(string) (string, error) { return "hash", nil }})
	{
		user, err := service.Register(context.Background(), Register{})
		require.Nil(t, user, "duplicate error: %v", err)
		require.ErrorIs(t, err, domain.ErrAlreadyExists, "duplicate error: %v", err)
	}
}

func TestLoginPolicies(t *testing.T) {
	unavailable := errors.New("database unavailable")
	signingFailure := errors.New("signing failed")
	for _, tc := range []struct {
		name             string
		lookupErr        error
		correct, blocked bool
		issueErr, want   error
	}{
		{name: "missing account", lookupErr: domain.ErrNotFound, want: ErrUnauthorized},
		{name: "repository failure", lookupErr: unavailable, want: unavailable},
		{name: "wrong password", want: ErrUnauthorized},
		{name: "blocked user", correct: true, blocked: true, want: ErrBlocked},
		{name: "blocked user wrong password", blocked: true, want: ErrUnauthorized},
		{name: "signing failure", correct: true, issueErr: signingFailure, want: signingFailure},
		{name: "success", correct: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			user := domain.User{ID: uuid.New(), PasswordHash: "stored", IsBlocked: tc.blocked}
			issued := false
			service := NewService(userStub{byEmail: func(ctx context.Context, email string) (*domain.User, error) {
				assert.Equal(t, "alice@example.com", email, "email not normalized")
				return &user, tc.lookupErr
			}}, nil, nil, credentialStub{
				verify: func(password, hash string) bool { return tc.correct },
				issue: func(id uuid.UUID) (string, error) {
					issued = true
					assert.Equal(t, user.ID, id, "wrong token subject")
					return "token", tc.issueErr
				},
			})
			token, err := service.Login(context.Background(), Login{Email: "ALICE@EXAMPLE.COM", Password: "password"})
			require.ErrorIs(t, err, tc.want, "got %v, want %v", err, tc.want)
			if tc.want == nil {
				assert.Equal(t, "token", token)
			}
			assert.Equal(t, (tc.correct && !tc.blocked && tc.lookupErr == nil), issued, "token issued before policy checks")
		})
	}
}

func TestAuthenticateUsesCurrentUserState(t *testing.T) {
	id := uuid.New()
	current := domain.User{ID: id, Role: domain.RoleUser}
	calls := 0
	service := NewService(userStub{byID: func(ctx context.Context, requested uuid.UUID) (*domain.User, error) {
		calls++
		assert.Equal(t, id, requested, "wrong token subject")
		copy := current
		return &copy, nil
	}}, nil, nil, credentialStub{parse: func(token string) (uuid.UUID, error) { return id, nil }})
	actor, err := service.Authenticate(context.Background(), "same-token")
	require.NoError(t, err, "first authentication: %v", err)
	assert.Equal(t, domain.RoleUser, actor.Role, "first authentication: %v", err)
	current.Role = domain.RoleAdmin
	actor, err = service.Authenticate(context.Background(), "same-token")
	require.NoError(t, err, "role change: %v", err)
	assert.Equal(t, domain.RoleAdmin, actor.Role, "role change: %v", err)
	current.IsBlocked = true
	{
		_, err := service.Authenticate(context.Background(), "same-token")
		require.ErrorIs(t, err, ErrBlocked, "block not honored: %v", err)
	}
	assert.EqualValues(t, 3, calls, "user state was cached")
}

func TestAuthenticateErrors(t *testing.T) {
	unavailable := errors.New("database unavailable")
	for _, tc := range []struct {
		name                      string
		parseErr, lookupErr, want error
	}{
		{"invalid token", errors.New("invalid signature"), nil, ErrUnauthorized},
		{"deleted account", nil, domain.ErrNotFound, ErrUnauthorized},
		{"repository failure", nil, unavailable, unavailable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lookedUp := false
			service := NewService(userStub{byID: func(context.Context, uuid.UUID) (*domain.User, error) { lookedUp = true; return nil, tc.lookupErr }}, nil, nil,
				credentialStub{parse: func(string) (uuid.UUID, error) { return uuid.New(), tc.parseErr }})
			{
				_, err := service.Authenticate(context.Background(), "token")
				require.ErrorIs(t, err, tc.want, "got %v, want %v", err, tc.want)
			}
			assert.Equal(t, (tc.parseErr == nil), lookedUp, "invalid token queried the database")
		})
	}
}

func TestAdminOperationsDenyBeforeRepositoryAccess(t *testing.T) {
	service := NewService(nil, nil, nil, nil)
	target := uuid.New()
	operations := []struct {
		name string
		run  func(*domain.User) error
	}{
		{"create movie", func(actor *domain.User) error {
			return service.CreateMovie(context.Background(), actor, &domain.Movie{})
		}},
		{"update movie", func(actor *domain.User) error {
			_, err := service.UpdateMovie(context.Background(), actor, &domain.Movie{})
			return err
		}},
		{"delete movie", func(actor *domain.User) error { return service.DeleteMovie(context.Background(), actor, target) }},
		{"list users", func(actor *domain.User) error {
			_, err := service.ListUsers(context.Background(), actor, 20, 0)
			return err
		}},
		{"get user", func(actor *domain.User) error {
			_, err := service.GetUser(context.Background(), actor, target)
			return err
		}},
		{"block user", func(actor *domain.User) error {
			_, err := service.SetUserBlocked(context.Background(), actor, target, true)
			return err
		}},
		{"delete user", func(actor *domain.User) error { return service.DeleteUser(context.Background(), actor, target) }},
	}
	for _, operation := range operations {
		for _, tc := range []struct {
			name  string
			actor *domain.User
			want  error
		}{
			{"anonymous", nil, ErrUnauthorized},
			{"ordinary user", &domain.User{Role: domain.RoleUser}, ErrForbidden},
			{"blocked admin", &domain.User{Role: domain.RoleAdmin, IsBlocked: true}, ErrBlocked},
		} {
			t.Run(operation.name+"/"+tc.name, func(t *testing.T) {
				{
					err := operation.run(tc.actor)
					require.ErrorIs(t, err, tc.want, "got %v, want %v", err, tc.want)
				}
			})
		}
	}
}

func TestAdminCannotModifyOwnAccount(t *testing.T) {
	service := NewService(nil, nil, nil, nil)
	actor := &domain.User{ID: uuid.New(), Role: domain.RoleAdmin}
	for _, blocked := range []bool{true, false} {
		{
			_, err := service.SetUserBlocked(context.Background(), actor, actor.ID, blocked)
			require.ErrorIs(t, err, ErrSelfModification, "self block: %v", err)
		}
	}
	{
		err := service.DeleteUser(context.Background(), actor, actor.ID)
		require.ErrorIs(t, err, ErrSelfModification, "self delete: %v", err)
	}
}

func TestAddCollectionItemKeepsOwner(t *testing.T) {
	userID, movieID := uuid.New(), uuid.New()
	movie := domain.Movie{ID: movieID, Name: "Film"}
	service := NewService(nil, movieStub{get: func(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
		assert.Equal(t, movieID, id, "wrong movie lookup")
		return &movie, nil
	}}, collectionStub{create: func(ctx context.Context, item *domain.CollectionItem) error {
		assert.Equal(t, userID, item.UserID, "wrong collection item: %+v", item)
		assert.Equal(t, movieID, item.MovieID, "wrong collection item: %+v", item)
		assert.Equal(t, domain.MovieStatusWatched, item.Status, "wrong collection item: %+v", item)
		require.NotEqual(t, uuid.Nil, item.ID, "wrong collection item: %+v", item)
		return nil
	}}, nil)
	entry, err := service.AddCollectionItem(context.Background(), userID, movieID, domain.MovieStatusWatched)
	require.NoError(t, err, "add collection: %v", err)
	assert.Equal(t, movieID, entry.Movie.ID, "add collection: %v", err)
	service = NewService(nil, movieStub{get: func(context.Context, uuid.UUID) (*domain.Movie, error) { return nil, domain.ErrNotFound }}, nil, nil)
	{
		_, err := service.AddCollectionItem(context.Background(), userID, movieID, domain.MovieStatusPlanned)
		require.ErrorIs(t, err, domain.ErrNotFound, "missing movie: %v", err)
	}
}

func TestListCollectionJoinsMoviesAndKeepsPagination(t *testing.T) {
	userID, movieID := uuid.New(), uuid.New()
	status := domain.MovieStatusWatched
	filter := domain.CollectionFilter{Status: &status, Limit: 5, Offset: 10}
	service := NewService(nil, movieStub{get: func(context.Context, uuid.UUID) (*domain.Movie, error) { return &domain.Movie{ID: movieID}, nil }},
		collectionStub{list: func(ctx context.Context, id uuid.UUID, got domain.CollectionFilter) ([]domain.CollectionItem, int, error) {
			assert.Equal(t, userID, id, "wrong collection filter")
			assert.EqualValues(t, 5, got.Limit, "wrong collection filter")
			assert.EqualValues(t, 10, got.Offset, "wrong collection filter")
			require.NotNil(t, got.Status, "wrong collection filter")
			assert.Equal(t, status, *got.Status, "wrong collection filter")
			return []domain.CollectionItem{{UserID: userID, MovieID: movieID}}, 12, nil
		}}, nil)
	page, err := service.ListCollection(context.Background(), userID, filter)
	require.NoError(t, err, "joined page: %+v, %v", page, err)
	assert.EqualValues(t, 12, page.Total, "joined page: %+v, %v", page, err)
	assert.EqualValues(t, 5, page.Limit, "joined page: %+v, %v", page, err)
	assert.EqualValues(t, 10, page.Offset, "joined page: %+v, %v", page, err)
	require.Len(t, page.Items, 1, "joined page: %+v, %v", page, err)
	assert.Equal(t, movieID, page.Items[0].Movie.ID, "joined page: %+v, %v", page, err)
	service.movies = movieStub{get: func(context.Context, uuid.UUID) (*domain.Movie, error) { return nil, domain.ErrNotFound }}
	{
		_, err := service.ListCollection(context.Background(), userID, filter)
		require.ErrorIs(t, err, domain.ErrNotFound, "movie join error: %v", err)
	}
}

type failingMovies struct {
	domain.MovieRepository
	err error
}

func (s failingMovies) Create(context.Context, *domain.Movie) error { return s.err }
func (s failingMovies) Update(context.Context, *domain.Movie) error { return s.err }
func (s failingMovies) Delete(context.Context, uuid.UUID) error     { return s.err }
func (s failingMovies) List(context.Context, domain.MovieFilter) ([]domain.Movie, int, error) {
	return nil, 0, s.err
}

type failingUsers struct {
	domain.UserRepository
	err  error
	user *domain.User
}

func (s failingUsers) GetByID(context.Context, uuid.UUID) (*domain.User, error) {
	if s.user != nil {
		return s.user, nil
	}
	return nil, s.err
}
func (s failingUsers) Update(context.Context, *domain.User) error { return s.err }
func (s failingUsers) Delete(context.Context, uuid.UUID) error    { return s.err }
func (s failingUsers) List(context.Context, int, int) ([]domain.User, int, error) {
	return nil, 0, s.err
}

type failingCollection struct {
	domain.CollectionRepository
	err error
}

func (s failingCollection) Get(context.Context, uuid.UUID, uuid.UUID) (*domain.CollectionItem, error) {
	return nil, s.err
}
func (s failingCollection) Create(context.Context, *domain.CollectionItem) error { return s.err }
func (s failingCollection) Update(context.Context, *domain.CollectionItem) error { return s.err }
func (s failingCollection) Delete(context.Context, uuid.UUID, uuid.UUID) error   { return s.err }
func (s failingCollection) List(context.Context, uuid.UUID, domain.CollectionFilter) ([]domain.CollectionItem, int, error) {
	return nil, 0, s.err
}

func TestRepositoryFailuresPropagate(t *testing.T) {
	ctx := context.Background()
	failure := errors.New("database unavailable")
	admin := &domain.User{ID: uuid.New(), Role: domain.RoleAdmin}
	id := uuid.New()
	movie := domain.NewMovie("Film", 90, domain.MovieGenreDrama, 2020, nil)
	item := domain.NewCollectionItem(id, movie.ID, domain.MovieStatusPlanned)
	service := NewService(failingUsers{err: failure}, failingMovies{err: failure}, failingCollection{err: failure}, nil)
	cases := []struct {
		name string
		run  func() error
	}{
		{"create movie", func() error { return service.CreateMovie(ctx, admin, movie) }},
		{"update movie stops before reload", func() error { _, err := service.UpdateMovie(ctx, admin, movie); return err }},
		{"delete movie", func() error { return service.DeleteMovie(ctx, admin, movie.ID) }},
		{"list movies", func() error { _, err := service.ListMovies(ctx, domain.MovieFilter{Limit: 20}); return err }},
		{"get collection stops before movie lookup", func() error { _, err := service.GetCollectionItem(ctx, id, movie.ID); return err }},
		{"update collection stops before reload", func() error { _, err := service.UpdateCollectionItem(ctx, item); return err }},
		{"delete collection", func() error { return service.DeleteCollectionItem(ctx, id, movie.ID) }},
		{"list collection stops before movie lookup", func() error {
			_, err := service.ListCollection(ctx, id, domain.CollectionFilter{Limit: 20})
			return err
		}},
		{"list users", func() error { _, err := service.ListUsers(ctx, admin, 20, 0); return err }},
		{"get user", func() error { _, err := service.GetUser(ctx, admin, id); return err }},
		{"block lookup", func() error { _, err := service.SetUserBlocked(ctx, admin, id, true); return err }},
		{"delete user", func() error { return service.DeleteUser(ctx, admin, id) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			{
				err := tc.run()
				require.ErrorIs(t, err, failure, "got %v, want %v", err, failure)
			}
		})
	}
	t.Run("add collection insert failure", func(t *testing.T) {
		service := NewService(nil, movieStub{get: func(context.Context, uuid.UUID) (*domain.Movie, error) { return movie, nil }}, failingCollection{err: failure}, nil)
		{
			_, err := service.AddCollectionItem(ctx, id, movie.ID, domain.MovieStatusPlanned)
			require.ErrorIs(t, err, failure)
		}
	})
	t.Run("block update failure", func(t *testing.T) {
		service := NewService(failingUsers{err: failure, user: &domain.User{ID: id}}, nil, nil, nil)
		{
			_, err := service.SetUserBlocked(ctx, admin, id, true)
			require.ErrorIs(t, err, failure)
		}
	})
	t.Run("get collection movie failure", func(t *testing.T) {
		service := NewService(nil, movieStub{get: func(context.Context, uuid.UUID) (*domain.Movie, error) { return nil, failure }}, collectionLookup{item: item}, nil)
		{
			_, err := service.GetCollectionItem(ctx, id, movie.ID)
			require.ErrorIs(t, err, failure)
		}
	})
}

type collectionLookup struct {
	domain.CollectionRepository
	item *domain.CollectionItem
}

func (s collectionLookup) Get(context.Context, uuid.UUID, uuid.UUID) (*domain.CollectionItem, error) {
	return s.item, nil
}
