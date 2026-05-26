package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EventsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "search_events_processed_total",
		Help: "Общее количество поисковых событий, обработанных из NATS",
	})

	EventsDropped = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "search_events_dropped_total",
		Help: "Общее количество пропущенных событий по причинам",
	}, []string{"reason"})

	TopRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "search_top_requests_total",
		Help: "Общее количество запросов GET /top",
	})

	RequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "search_http_duration_seconds",
		Help:    "Длительность HTTP-запроса в секундах",
		Buckets: prometheus.DefBuckets,
	})

	ActiveQueries = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "search_active_queries_total",
		Help: "Количество уникальных запросов в текущем скользящем окне",
	})
)
