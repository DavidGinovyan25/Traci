package httptransport

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"traci/backend/internal/gen"
)

type catalogHandler struct {
	*Handler
	params gen.GetCatalogParams
	called bool
	serve  http.HandlerFunc
}

func (h *catalogHandler) GetCatalog(w http.ResponseWriter, r *http.Request, params gen.GetCatalogParams) {
	if h.serve != nil {
		h.serve(w, r)
		return
	}
	h.params = params
	h.called = true
	w.WriteHeader(http.StatusNoContent)
}

func TestGeneratedRouter(t *testing.T) {
	h := &catalogHandler{Handler: NewHandler(nil, slog.New(slog.NewTextHandler(io.Discard, nil)), time.Hour)}
	router, err := NewRouter(h, slog.New(slog.NewTextHandler(io.Discard, nil)))
	require.NoError(t, err)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/catalog?limit=5&offset=10&genre=DRAMA", nil))
	assert.Equal(t, http.StatusNoContent, response.Code, "handler not invoked: %d", response.Code)
	require.True(t, h.called, "handler not invoked: %d", response.Code)
	require.NotNil(t, h.params.Limit, "unexpected parameters: %+v", h.params)
	assert.EqualValues(t, 5, *h.params.Limit, "unexpected parameters: %+v", h.params)
	require.NotNil(t, h.params.Offset, "unexpected parameters: %+v", h.params)
	assert.EqualValues(t, 10, *h.params.Offset, "unexpected parameters: %+v", h.params)
	require.NotNil(t, h.params.Genre, "unexpected parameters: %+v", h.params)
	assert.Equal(t, gen.DRAMA, *h.params.Genre, "unexpected parameters: %+v", h.params)
	for _, path := range []string{"/catalog/not-a-uuid", "/catalog?limit=invalid"} {
		response = httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		var body gen.ErrorResponse
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.Equal(t, http.StatusBadRequest, response.Code, "unexpected error: %+v", body)
		assert.Equal(t, "BAD_REQUEST", body.Code, "unexpected error: %+v", body)
	}
	response = httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	assert.Equal(t, http.StatusOK, response.Code, "health status: %d", response.Code)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/auth/login", nil))
	assert.Equal(t, http.StatusBadRequest, response.Code, "invalid body status: %d", response.Code)
}

func TestRouterRequestLogging(t *testing.T) {
	for _, tc := range []struct {
		name, method, path, body string
		status                   int
	}{
		{"health", "GET", "/healthz?token=private-query", "", 200},
		{"unauthorized", "GET", "/catalog", "", 401},
		{"invalid query", "GET", "/catalog?limit=invalid", "", 400},
		{"large offset", "GET", "/catalog?offset=2147483648", "", 400},
		{"missing route", "GET", "/missing", "", 404},
		{"invalid body", "POST", "/auth/login", "{}", 400},
		{"body limit", "POST", "/auth/login", strings.Repeat(" ", 1<<20) + `{"email":"a@example.com","password":"private-password"}`, 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger, logs := newTestLogger()
			router, err := NewRouter(NewHandler(nil, logger, time.Hour), logger)
			require.NoError(t, err)
			request := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			if tc.body != "" {
				request.Header.Set("Content-Type", "application/json")
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			require.Equal(t, tc.status, response.Code)
			entries := logs.entries(t)
			require.NotEmpty(t, entries)
			entry := entries[len(entries)-1]
			assert.Equal(t, "HTTP request completed", entry["msg"])
			assert.Equal(t, "INFO", entry["level"])
			assert.Equal(t, tc.method, entry["http.request.method"])
			assert.Equal(t, request.URL.Path, entry["url.path"])
			assert.EqualValues(t, tc.status, entry["http.response.status_code"])
			duration, ok := entry["duration_ms"].(float64)
			require.True(t, ok)
			assert.GreaterOrEqual(t, duration, float64(0))
			assert.NotContains(t, logs.text(), "private-query")
			assert.NotContains(t, logs.text(), "private-password")
		})
	}
}

func TestRouterPanicRecovery(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{"before headers", 500, ""},
		{"after explicit headers", 202, "partial"},
		{"after implicit headers", 200, "partial"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger, logs := newTestLogger()
			handler := &catalogHandler{Handler: NewHandler(nil, logger, time.Hour)}
			handler.serve = func(w http.ResponseWriter, r *http.Request) {
				if tc.status == 202 {
					w.WriteHeader(tc.status)
					w.WriteHeader(http.StatusConflict)
				}
				if tc.body != "" {
					_, _ = w.Write([]byte(tc.body))
				}
				panic("test panic")
			}
			router, err := NewRouter(handler, logger)
			require.NoError(t, err)
			response := httptest.NewRecorder()
			require.NotPanics(t, func() { router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/catalog", nil)) })
			require.Equal(t, tc.status, response.Code)
			if tc.body != "" {
				assert.Equal(t, tc.body, response.Body.String())
			} else {
				var failure gen.ErrorResponse
				require.NoError(t, json.Unmarshal(response.Body.Bytes(), &failure))
				assert.Equal(t, "INTERNAL_ERROR", failure.Code)
				assert.NotContains(t, failure.Message, "test panic")
			}
			entries := logs.entries(t)
			require.Len(t, entries, 2)
			assert.Equal(t, "ERROR", entries[0]["level"])
			assert.Equal(t, "HTTP handler panic", entries[0]["msg"])
			assert.Equal(t, "panic", entries[0]["error.type"])
			assert.Equal(t, "test panic", entries[0]["exception.message"])
			assert.Equal(t, "GET", entries[0]["http.request.method"])
			assert.Equal(t, "/catalog", entries[0]["url.path"])
			assert.Equal(t, "HTTP request completed", entries[1]["msg"])
			assert.EqualValues(t, tc.status, entries[1]["http.response.status_code"])
		})
	}
}
