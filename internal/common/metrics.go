package common

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// metricsMap stores counts: route -> method -> count
var (
	metricsMap = make(map[string]map[string]int)
	metricsMu  sync.RWMutex
)

// MetricsMiddleware increments the counter for each route and method
func MetricsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get route pattern
			route := "unknown"
			if mux.CurrentRoute(r) != nil {
				if p, err := mux.CurrentRoute(r).GetPathTemplate(); err == nil {
					route = p
				}
			}
			method := r.Method

			metricsMu.Lock()
			if _, ok := metricsMap[route]; !ok {
				metricsMap[route] = make(map[string]int)
			}
			metricsMap[route][method]++
			metricsMu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

// MetricsHandler serves the metrics in Prometheus format
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	metricsMu.RLock()
	defer metricsMu.RUnlock()
	for route, methods := range metricsMap {
		for method, count := range methods {
			fmt.Fprintf(w, "petstore_http_requests_total{path=\"%s\",method=\"%s\"} %d\n", route, method, count)
		}
	}
}
