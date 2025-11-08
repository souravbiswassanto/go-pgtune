package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	pgtune "github.com/souravbiswassanto/go-pgtune/pkg/pgtune"
)

// TuneRequest represents the API request body
type TuneRequest struct {
	DBVersion      float64 `json:"db_version"`
	OSType         string  `json:"os_type"`
	DBType         string  `json:"db_type"`
	TotalMemoryGB  float64 `json:"total_memory_gb"`
	CPUNum         int     `json:"cpu_num"`
	MaxConnections int     `json:"max_connections"`
	StorageType    string  `json:"storage_type"`
}

// TuneResponse represents the API response
type TuneResponse struct {
	Success bool               `json:"success"`
	Config  *pgtune.TuneOutput `json:"config,omitempty"`
	Error   string             `json:"error,omitempty"`
}

func tuneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TuneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TuneResponse{
			Success: false,
			Error:   fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	// Convert GB to bytes
	totalMemoryBytes := int64(req.TotalMemoryGB * float64(pgtune.GB))

	input := pgtune.TuneInput{
		DBVersion:      req.DBVersion,
		OSType:         req.OSType,
		DBType:         req.DBType,
		TotalMemory:    totalMemoryBytes,
		CPUNum:         req.CPUNum,
		MaxConnections: req.MaxConnections,
		StorageType:    req.StorageType,
	}

	output, err := pgtune.Tune(input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TuneResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TuneResponse{
		Success: true,
		Config:  output,
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func main() {
	http.HandleFunc("/tune", tuneHandler)
	http.HandleFunc("/health", healthHandler)

	port := "8080"
	if envPort := getEnv("PORT", ""); envPort != "" {
		port = envPort
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("POST /tune - Generate PostgreSQL tuning configuration")
	log.Printf("GET  /health - Health check")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultValue string) string {
	// Simple environment variable getter
	return defaultValue
}
