package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Response structure for JSON logging and API output
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Env       string `json:"environment"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get environment variable (DEV, UAT, PROD) default to DEV
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "DEV"
	}

	// 2. Prepare the response data
	res := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
		Env:       env,
	}

	// 3. Set content type to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 4. Send the JSON response
	json.NewEncoder(w).Encode(res)

	// 5. Output structured JSON log (Requirement 5a)
	logData, _ := json.Marshal(map[string]interface{}{
		"level":   "info",
		"message": "Health check accessed",
		"path":    r.URL.Path,
		"env":     env,
		"time":    res.Timestamp,
	})
	fmt.Println(string(logData))
}

func main() {
	// Configure standard logger to only print the message (we handle JSON formatting ourselves)
	log.SetFlags(0)

	http.HandleFunc("/health", healthHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	startupLog, _ := json.Marshal(map[string]interface{}{
		"level":   "info",
		"message": "Starting API Service",
		"port":    port,
	})
	fmt.Println(string(startupLog))

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		errorLog, _ := json.Marshal(map[string]interface{}{
			"level":   "fatal",
			"message": "Failed to start server",
			"error":   err.Error(),
		})
		fmt.Println(string(errorLog))
		os.Exit(1)
	}
}
