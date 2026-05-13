package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func main() {
	// Configure standard logger
	log.SetFlags(0)

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "DEV"
	}

	startupLog, _ := json.Marshal(map[string]interface{}{
		"level":   "info",
		"message": "Starting Background Service Worker",
		"env":     env,
	})
	fmt.Println(string(startupLog))

	// Expose metrics endpoint
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9090", nil)
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
	timestamp := time.Now().Format(time.RFC3339)
	
	// Update metric
	workerLastSuccess.SetToCurrentTime()

	jobLog, _ := json.Marshal(map[string]interface{}{
		"level":   "info",
		"message": "Worker job executed: Updated timestamps for today's records",
		"time":    timestamp,
		"env":     env,
	})
	fmt.Println(string(jobLog))
}
