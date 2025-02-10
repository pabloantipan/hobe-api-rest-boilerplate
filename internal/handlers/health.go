package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pabloantipan/hobe-api-boilerplate/config"
)

type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
	Timestamp string            `json:"timestamp"`
}

type HealthHandler struct {
	cfg *config.Config
}

func NewHealthHandler(cfg *config.Config) HealthHandlerInterface {
	return &HealthHandler{
		cfg: cfg,
	}
}

type HealthHandlerInterface interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status:    "UP",
		Version:   h.cfg.Version,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Services: map[string]string{
			"players": "UP",
			// Add more services as they come
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
