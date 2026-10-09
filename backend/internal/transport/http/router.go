package httptransport

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/legacy"

	"traci/backend/internal/gen"
)

func NewRouter(handler gen.ServerInterface, logger *slog.Logger) (http.Handler, error) {
	specification, err := gen.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("load OpenAPI schema: %w", err)
	}
	specification.Servers = nil
	validator, err := legacy.NewRouter(specification)
	if err != nil {
		return nil, fmt.Errorf("create OpenAPI validator: %w", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	router := gen.HandlerWithOptions(handler, gen.StdHTTPServerOptions{
		BaseRouter: mux,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Invalid path or query parameters")
		},
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			logger.Info("http request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
			if recovered := recover(); recovered != nil {
				logger.Error("http handler panic", "panic", recovered)
				writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
			}
		}()
		if r.URL.Path != "/healthz" {
			route, params, err := validator.FindRoute(r)
			if err == nil {
				r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
				input := &openapi3filter.RequestValidationInput{
					Request: r, PathParams: params, Route: route,
					Options: &openapi3filter.Options{AuthenticationFunc: openapi3filter.NoopAuthenticationFunc, SkipSettingDefaults: true},
				}
				if err := openapi3filter.ValidateRequest(r.Context(), input); err != nil {
					writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Request does not match API schema")
					return
				}
				if offset := r.URL.Query().Get("offset"); offset != "" {
					number, err := strconv.ParseInt(offset, 10, 64)
					if err != nil || number > math.MaxInt32 {
						writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Offset is too large")
						return
					}
				}
			}
		}
		router.ServeHTTP(w, r)
	}), nil
}
