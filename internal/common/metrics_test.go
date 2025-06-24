package common

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetricsMiddlewareAndHandler(t *testing.T) {
	// Reset metrics map for test isolation
	metricsMu.Lock()
	metricsMap = make(map[string]map[string]int)
	metricsMu.Unlock()

	// Create a test route and handler
	r := http.NewServeMux()
	r.Handle("/test", MetricsMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})))
	r.HandleFunc("/metrics", MetricsHandler)

	// Simulate a request to /test
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Simulate a request to /metrics
	metricsReq := httptest.NewRequest("GET", "/metrics", nil)
	metricsW := httptest.NewRecorder()
	r.ServeHTTP(metricsW, metricsReq)
	metricsBody, _ := io.ReadAll(metricsW.Body)
	output := string(metricsBody)

	if !strings.Contains(output, "petstore_http_requests_total") {
		t.Errorf("Expected metric name in output, got: %s", output)
	}
	if !strings.Contains(output, "method=\"GET\"") {
		t.Errorf("Expected method label in output, got: %s", output)
	}
	if !strings.Contains(output, "path=\"/test\"") && !strings.Contains(output, "path=\"unknown\"") {
		t.Errorf("Expected path label in output, got: %s", output)
	}
} 