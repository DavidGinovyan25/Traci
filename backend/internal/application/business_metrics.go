package application

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type BusinessMetrics struct {
	authenticationsTotal     *prometheus.CounterVec
	collectionAdditionsTotal prometheus.Counter
	adminMoviesAddedTotal    prometheus.Counter
	moviesCurrent            prometheus.Gauge
}

func NewBusinessMetrics(reg prometheus.Registerer) *BusinessMetrics {
	metricFactory := promauto.With(reg)

	return &BusinessMetrics{
		authenticationsTotal: metricFactory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "traci_authentications_total",
				Help: "Total number of authentication attempts, partitioned by result.",
			},
			[]string{"result"},
		),
		collectionAdditionsTotal: metricFactory.NewCounter(
			prometheus.CounterOpts{
				Name: "traci_collection_additions_total",
				Help: "Total number of movies successfully added to user collections.",
			},
		),
		adminMoviesAddedTotal: metricFactory.NewCounter(
			prometheus.CounterOpts{
				Name: "traci_admin_movies_added_total",
				Help: "Total number of movies added to the catalog by administrators.",
			},
		),
		moviesCurrent: metricFactory.NewGauge(
			prometheus.GaugeOpts{
				Name: "traci_movies_current",
				Help: "Current number of movies in the application catalog.",
			},
		),
	}
}
