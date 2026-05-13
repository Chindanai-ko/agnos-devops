package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Response structure for API output
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Env       string `json:"environment"`
}

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "code"},
	)
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows exposing the HTTP status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		path := r.URL.Path
		method := r.Method
		code := strconv.Itoa(rw.statusCode)

		httpRequestsTotal.WithLabelValues(path, method, code).Inc()
		httpRequestDuration.WithLabelValues(path, method).Observe(duration)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "DEV"
	}

	res := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
		Env:       env,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)

	slog.Info("Health check accessed",
		"path", r.URL.Path,
		"env", env,
		"time", res.Timestamp,
	)
}

func main() {
	// Configure slog for JSON output
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Setup handlers with metrics middleware
	http.Handle("/health", metricsMiddleware(http.HandlerFunc(healthHandler)))
	http.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Starting API Service", "port", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
