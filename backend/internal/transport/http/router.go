package httptransport

import (
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

type responseWriter struct {
	http.ResponseWriter
	status        int
	headerWritten bool
}

func (w *responseWriter) WriteHeader(status int) {
	if w.headerWritten {
		return
	}
	w.status = status
	w.headerWritten = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	if !w.headerWritten {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}

func (w *responseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func NewRouter(handler gen.ServerInterface, logger *slog.Logger, metrics ...*HTTPMetrics) (http.Handler, error) {
	var httpMetrics *HTTPMetrics
	if len(metrics) > 0 {
		httpMetrics = metrics[0]
	}
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
		writeJSON(
			w, r,
			http.StatusOK,
			map[string]string{"status": "ok"},
			logger,
		)
	})
	router := gen.HandlerWithOptions(handler, gen.StdHTTPServerOptions{
		BaseRouter: mux,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			writeError(
				w, r,
				http.StatusBadRequest,
				"BAD_REQUEST",
				"Invalid path or query parameters",
				logger,
			)
		},
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		routeName := "unmatched"
		if r.URL.Path == "/healthz" {
			routeName = "/healthz"
		}
		if httpMetrics != nil {
			httpMetrics.httpRequestsInFlight.Inc()
		}
		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.ErrorContext(
					r.Context(),
					"HTTP handler panic",
					"error.type", "panic",
					"exception.message", fmt.Sprint(recovered),
					"http.request.method", r.Method,
					"url.path", r.URL.Path,
				)
				if !rw.headerWritten {
					writeError(
						rw, r,
						http.StatusInternalServerError,
						"INTERNAL_ERROR",
						"Internal server error",
						logger,
					)
				}
			}
			if httpMetrics != nil {
				httpMetrics.httpRequestsInFlight.Dec()
				method := metricMethod(r.Method)
				httpMetrics.httpRequestsTotal.WithLabelValues(routeName, method, strconv.Itoa(rw.status)).Inc()
				httpMetrics.httpRequestDurationSeconds.WithLabelValues(routeName, method).Observe(time.Since(start).Seconds())
			}
			logger.InfoContext(
				r.Context(),
				"HTTP request completed",
				"http.request.method", r.Method,
				"url.path", r.URL.Path,
				"http.response.status_code", rw.status,
				"duration_ms", float64(time.Since(start).Microseconds())/1000,
			)
		}()
		if r.URL.Path != "/healthz" {
			route, params, err := validator.FindRoute(r)
			if err == nil {
				routeName = route.Path
				r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
				input := &openapi3filter.RequestValidationInput{
					Request:    r,
					PathParams: params,
					Route:      route,
					Options: &openapi3filter.Options{
						AuthenticationFunc:  openapi3filter.NoopAuthenticationFunc,
						SkipSettingDefaults: true,
					},
				}
				if err := openapi3filter.ValidateRequest(r.Context(), input); err != nil {
					writeError(
						rw, r,
						http.StatusBadRequest,
						"BAD_REQUEST",
						"Request does not match API schema",
						logger,
					)
					return
				}
				if offset := r.URL.Query().Get("offset"); offset != "" {
					number, err := strconv.ParseInt(offset, 10, 64)
					if err != nil || number > math.MaxInt32 {
						writeError(
							rw, r,
							http.StatusBadRequest,
							"BAD_REQUEST",
							"Offset is too large",
							logger,
						)
						return
					}
				}
			}
		}
		router.ServeHTTP(rw, r)
	}), nil
}

func metricMethod(method string) string {
	switch method {
	case
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions:
		return method
	default:
		return "OTHER"
	}
}
