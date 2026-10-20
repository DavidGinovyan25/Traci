package httptransport

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"traci/backend/internal/gen"
)

type failingResponseWriter struct {
	*httptest.ResponseRecorder
	err error
}

func (w *failingResponseWriter) Write([]byte) (int, error) { return 0, w.err }

func TestWriteJSON(t *testing.T) {
	logger, logs := newTestLogger()
	response := httptest.NewRecorder()
	writeJSON(response, httptest.NewRequest(http.MethodGet, "/healthz", nil), http.StatusOK, map[string]string{"status": "ok"}, logger)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"status":"ok"}`, response.Body.String())
	assert.Empty(t, logs.entries(t))
}

func TestWriteJSONSerializationFailure(t *testing.T) {
	logger, logs := newTestLogger()
	response := httptest.NewRecorder()
	writeJSON(response, httptest.NewRequest(http.MethodGet, "/catalog", nil), http.StatusCreated, make(chan int), logger)
	require.Equal(t, http.StatusInternalServerError, response.Code)
	var body gen.ErrorResponse
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	assert.Equal(t, "INTERNAL_ERROR", body.Code)
	assert.Equal(t, "/catalog", body.Path)
	assert.NotContains(t, body.Message, "chan")
	entries := logs.entries(t)
	require.Len(t, entries, 1)
	assert.Equal(t, "ERROR", entries[0]["level"])
	assert.Equal(t, "JSON serialization failed", entries[0]["msg"])
	assert.Contains(t, entries[0]["error"], "unsupported type")
}

func TestResponseWriteFailuresAreLogged(t *testing.T) {
	for _, tc := range []struct{ name, level, message string }{
		{"json", "WARN", "HTTP response write failed"},
		{"error", "ERROR", "JSON response encoding failed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger, logs := newTestLogger()
			response := &failingResponseWriter{httptest.NewRecorder(), errors.New("connection closed")}
			request := httptest.NewRequest(http.MethodGet, "/catalog", nil)
			if tc.name == "json" {
				writeJSON(response, request, http.StatusOK, map[string]string{"status": "ok"}, logger)
			} else {
				writeError(response, request, http.StatusBadRequest, "BAD_REQUEST", "Invalid request", logger)
			}
			entries := logs.entries(t)
			require.Len(t, entries, 1)
			assert.Equal(t, tc.level, entries[0]["level"])
			assert.Equal(t, tc.message, entries[0]["msg"])
			assert.Equal(t, "connection closed", entries[0]["error"])
		})
	}
}

func TestWriteError(t *testing.T) {
	logger, logs := newTestLogger()
	response := httptest.NewRecorder()
	writeError(response, httptest.NewRequest(http.MethodGet, "/catalog?token=secret", nil), http.StatusBadRequest, "BAD_REQUEST", "Invalid request", logger)
	var body gen.ErrorResponse
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	assert.Equal(t, "/catalog", body.Path)
	assert.Equal(t, "BAD_REQUEST", body.Code)
	assert.Equal(t, "Bad Request", body.Error)
	assert.Equal(t, "Invalid request", body.Message)
	assert.EqualValues(t, http.StatusBadRequest, body.Status)
	assert.WithinDuration(t, time.Now(), body.Timestamp, time.Second)
	assert.Empty(t, logs.entries(t))
}

func TestDecodeBody(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		ok         bool
		message    string
	}{
		{"valid", `{"email":"alice@example.com","password":"secret"}`, true, ""},
		{"empty", "", false, "Invalid JSON body"},
		{"invalid", `{`, false, "Invalid JSON body"},
		{"unknown field", `{"extra":true}`, false, "Invalid JSON body"},
		{"multiple objects", `{} {}`, false, "Expected a single JSON object"},
		{"trailing garbage", `{} garbage`, false, "Expected a single JSON object"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger, _ := newTestLogger()
			response := httptest.NewRecorder()
			body, ok := decodeBody[gen.LoginRequest](response, httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(tc.body)), logger)
			require.Equal(t, tc.ok, ok)
			if ok {
				assert.Equal(t, "secret", body.Password)
				assert.Empty(t, response.Body.String())
				return
			}
			require.Equal(t, http.StatusBadRequest, response.Code)
			var failure gen.ErrorResponse
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &failure))
			assert.Equal(t, "BAD_REQUEST", failure.Code)
			assert.Equal(t, tc.message, failure.Message)
		})
	}
}
