package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "runners_app_http_requests",
			Help: "Total number of HTTP requests",
		},
	)

	GetRunnerHttpResponsesCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "runners_app_get_runner_http_responses",
			Help: "Total number of HTTP responses for get runner API",
		},
		[]string{"status"},
	)

	GetAllRunnersTimer = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name: "runners_app_get_all_runners_duration",
			Help: "Duration of get all runners operation",
		},
	)
)
