package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	HTTPRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed.",
		},
	)

	HTTPErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors.",
		},
		[]string{"method", "path", "status"},
	)

	URLRedirectsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "url_redirects_total",
			Help: "Total number of URL redirects.",
		},
	)

	CacheHitsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits.",
		},
	)

	CacheMissesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses.",
		},
	)

	AnalyticsEventsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "analytics_events_total",
			Help: "Total number of analytics events received.",
		},
	)

	AnalyticsEventsProcessedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "analytics_events_processed_total",
			Help: "Total number of analytics events processed successfully.",
		},
	)

	AnalyticsEventsFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "analytics_events_failed_total",
			Help: "Total number of analytics events that failed to process.",
		},
	)

	AnalyticsQueueSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "analytics_queue_size",
			Help: "Current number of analytics events waiting in the queue.",
		},
	)

	AnalyticsProcessingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "analytics_processing_duration_seconds",
			Help:    "Time spent processing analytics events.",
			Buckets: prometheus.DefBuckets,
		},
	)

	CacheInvalidationErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_invalidation_errors_total",
			Help: "Total number of cache invalidation errors.",
		},
	)
)

func Init() {
	prometheus.MustRegister(
		HTTPRequestsTotal,
		HTTPRequestDuration,
		HTTPRequestsInFlight,
		HTTPErrorTotal,
		URLRedirectsTotal,
		CacheHitsTotal,
		CacheMissesTotal,
		AnalyticsEventsTotal,
		AnalyticsEventsProcessedTotal,
		AnalyticsEventsFailedTotal,
		AnalyticsQueueSize,
		AnalyticsProcessingDuration,
		CacheInvalidationErrorsTotal,
	)
}
