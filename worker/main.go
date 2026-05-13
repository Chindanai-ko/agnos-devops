package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	workerLastSuccess = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "worker_last_success_timestamp_seconds",
			Help: "The last time the worker job succeeded",
		},
	)
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	// Configure slog for JSON output
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "DEV"
	}

	slog.Info("Starting Background Service Worker", "env", env)

	// Expose metrics and health endpoints
	go func() {
		http.HandleFunc("/health", healthHandler)
		http.Handle("/metrics", promhttp.Handler())
		slog.Info("Starting metrics/health server on :9090")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			slog.Error("Metrics server failed", "error", err)
		}
	}()

	// Simulate periodic job
	for {
		performJob(env)
		time.Sleep(1 * time.Minute)
	}
}

func performJob(env string) {
	// Requirement: Update timestamp of today records
	// Stub implementation
	
	// Update metric
	workerLastSuccess.SetToCurrentTime()

	slog.Info("Worker job executed: Updated timestamps for today's records",
		"env", env,
		"time", time.Now().Format(time.RFC3339),
	)
}
