package httptransport

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
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
}

func (h *catalogHandler) GetCatalog(w http.ResponseWriter, r *http.Request, params gen.GetCatalogParams) {
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
