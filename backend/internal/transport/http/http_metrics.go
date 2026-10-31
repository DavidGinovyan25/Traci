package httptransport

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type HTTPMetrics struct {
	httpRequestsTotal          *prometheus.CounterVec
	httpRequestDurationSeconds *prometheus.HistogramVec
	httpRequestsInFlight       prometheus.Gauge
}

func NewHTTPMetrics(reg prometheus.Registerer) *HTTPMetrics {
	metricFactory := promauto.With(reg)

	return &HTTPMetrics{
		httpRequestsTotal: metricFactory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "traci_http_requests_total",
				Help: "Total number of HTTP requests handled, partitioned by route, method, and status code.",
			},
			[]string{"route", "method", "status"},
		),
		httpRequestDurationSeconds: metricFactory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "traci_http_request_duration_seconds",
				Help: "Duration of HTTP requests in seconds, partitioned by route and method.",
			},
			[]string{"route", "method"},
		),
		httpRequestsInFlight: metricFactory.NewGauge(
			prometheus.GaugeOpts{
				Name: "traci_http_requests_in_flight",
				Help: "Number of HTTP requests currently being processed.",
			},
		),
	}
}