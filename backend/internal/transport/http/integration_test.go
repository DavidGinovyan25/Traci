package httptransport

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"traci/backend/internal/application"
	"traci/backend/internal/gen"
	"traci/backend/internal/gen/db"
	"traci/backend/internal/infrastructure/auth"
	sqlrepository "traci/backend/internal/repository/sql"
	"traci/backend/internal/testutil"
)

func TestHTTPPostgres(t *testing.T) {
	url := testutil.Postgres(t)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, url)
	require.NoError(t, err)
	defer conn.Close(ctx)
	queries := db.New(conn)
	credentials, err := auth.NewCredentials(strings.Repeat("a", 32), time.Hour)
	require.NoError(t, err)
	service := application.NewService(sqlrepository.NewUserRepository(queries), sqlrepository.NewMovieRepository(queries), sqlrepository.NewCollectionRepository(queries), credentials)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router, err := NewRouter(NewHandler(service, logger, time.Hour), logger)
	require.NoError(t, err)
	request := func(method, path, token, body string, want int) []byte {
		t.Helper()
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)
		require.Equal(t, want, response.Code, "%s %s: got %d, want %d: %s", method, path, response.Code, want, response.Body.String())
		require.NotContains(t, response.Body.String(), "password_hash", "password hash leaked")
		return response.Body.Bytes()
	}
	decode := func(data []byte, target any) {
		t.Helper()
		require.NoError(t, json.Unmarshal(data, target))
	}
	var user, admin gen.User
	decode(request("POST", "/auth/register", "", `{"username":"alice","first_name":"Alice","second_name":"User","email":"ALICE@example.com","password":"password123"}`, 201), &user)
	decode(request("POST", "/auth/register", "", `{"username":"bob","first_name":"Bob","second_name":"Admin","email":"bob@example.com","password":"password123"}`, 201), &admin)
	require.Equal(t, gen.RoleUser, user.Role, "registration defaults not applied")
	require.Equal(t, "alice@example.com", string(user.Email), "registration defaults not applied")
	{
		_, err := conn.Exec(ctx, "UPDATE users SET role = 'admin' WHERE id = $1", admin.Id)
		require.NoError(t, err)
	}
	request("POST", "/auth/register", "", `{"username":"alice","first_name":"Alice","second_name":"User","email":"alice@example.com","password":"password123"}`, 409)
	request("POST", "/auth/register", "", `{"username":"mallory","first_name":"M","second_name":"M","email":"m@example.com","password":"password123","role":"admin"}`, 400)
	request("POST", "/auth/login", "", `{"email":"alice@example.com","password":"wrongpassword"}`, 401)
	var userLogin, adminLogin gen.LoginResponse
	decode(request("POST", "/auth/login", "", `{"email":"alice@example.com","password":"password123"}`, 200), &userLogin)
	decode(request("POST", "/auth/login", "", `{"email":"bob@example.com","password":"password123"}`, 200), &adminLogin)
	userToken, adminToken := userLogin.AccessToken, adminLogin.AccessToken
	request("GET", "/catalog", "", "", 401)
	request("GET", "/catalog", "invalid", "", 401)
	request("GET", "/users", userToken, "", 403)
	movieBody := `{"name":"Film","duration_min":90,"genre":"DRAMA","release_year":2020,"description":"original"}`
	request("POST", "/catalog", userToken, movieBody, 403)
	request("POST", "/catalog", adminToken, `{"name":"Film","duration_min":0,"genre":"DRAMA","release_year":2020}`, 400)
	request("POST", "/catalog", adminToken, `{"name":"Film","duration_min":90,"genre":"INVALID","release_year":2020}`, 400)
	request("POST", "/catalog", adminToken, `{"name":"Film","duration_min":90,"genre":"DRAMA"}`, 400)
	var movie gen.Movie
	decode(request("POST", "/catalog", adminToken, movieBody, 201), &movie)
	moviePath := "/catalog/" + movie.Id.String()
	request("GET", moviePath, userToken, "", 200)
	request("GET", "/catalog/not-a-uuid", userToken, "", 400)
	request("GET", "/catalog?limit=0", userToken, "", 400)
	request("GET", "/catalog?genre=INVALID", userToken, "", 400)
	var page gen.MoviePage
	decode(request("GET", "/catalog?search=film&genre=DRAMA", userToken, "", 200), &page)
	require.EqualValues(t, 1, page.Total, "invalid page: %+v", page)
	require.Len(t, page.Items, 1, "invalid page: %+v", page)
	require.EqualValues(t, 20, page.Limit, "invalid page: %+v", page)
	request("PATCH", moviePath, userToken, `{"name":"No"}`, 403)
	request("PATCH", moviePath, adminToken, `{}`, 400)
	request("PATCH", moviePath, adminToken, `{"name":null}`, 400)
	decode(request("PATCH", moviePath, adminToken, `{"description":null}`, 200), &movie)
	require.True(t, movie.Description.IsNull(), "movie patch semantics broken")
	require.Equal(t, "Film", movie.Name, "movie patch semantics broken")
	collectionBody := `{"movie_id":"` + movie.Id.String() + `"}`
	var item gen.CollectionItem
	decode(request("POST", "/collection", userToken, collectionBody, 201), &item)
	require.Equal(t, gen.PLANNED, item.Status, "collection default not applied")
	request("POST", "/collection", userToken, collectionBody, 409)
	collectionPath := "/collection/" + movie.Id.String()
	request("GET", collectionPath, adminToken, "", 404)
	request("PATCH", collectionPath, userToken, `{"personal_rating":11}`, 400)
	request("PATCH", collectionPath, userToken, `{"status":null}`, 400)
	decode(request("PATCH", collectionPath, userToken, `{"personal_rating":8,"review":"nice","status":"WATCHED"}`, 200), &item)
	require.EqualValues(t, 8, item.PersonalRating.MustGet(), "collection update failed")
	require.Equal(t, "nice", item.Review.MustGet(), "collection update failed")
	decode(request("PATCH", collectionPath, userToken, `{"personal_rating":null}`, 200), &item)
	require.True(t, item.PersonalRating.IsNull(), "collection patch semantics broken")
	require.Equal(t, "nice", item.Review.MustGet(), "collection patch semantics broken")
	var collectionPage gen.CollectionPage
	decode(request("GET", "/collection?status=WATCHED", userToken, "", 200), &collectionPage)
	require.EqualValues(t, 1, collectionPage.Total, "collection filter failed")
	require.Len(t, collectionPage.Items, 1, "collection filter failed")
	userPath, adminPath := "/users/"+user.Id.String(), "/users/"+admin.Id.String()
	request("PATCH", adminPath, adminToken, `{"is_blocked":true}`, 409)
	request("DELETE", adminPath, adminToken, "", 409)
	request("PATCH", userPath, adminToken, `{}`, 400)
	request("PATCH", userPath, adminToken, `{"is_blocked":true}`, 200)
	request("GET", "/catalog", userToken, "", 403)
	request("POST", "/auth/login", "", `{"email":"alice@example.com","password":"password123"}`, 403)
	request("PATCH", userPath, adminToken, `{"is_blocked":false}`, 200)
	request("GET", userPath, adminToken, "", 200)
	request("GET", "/users", adminToken, "", 200)
	request("DELETE", collectionPath, userToken, "", 204)
	request("GET", collectionPath, userToken, "", 404)
	request("POST", "/collection", userToken, collectionBody, 201)
	request("DELETE", moviePath, adminToken, "", 204)
	request("GET", collectionPath, userToken, "", 404)
	request("GET", moviePath, userToken, "", 404)
	request("DELETE", userPath, adminToken, "", 204)
	request("GET", "/catalog", userToken, "", 401)
}
