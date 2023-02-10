package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9000", nil)
}
