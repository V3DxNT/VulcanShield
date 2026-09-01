package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vulcanshield/backend/internal/middleware"
)

const version = "0.1.0"


type Prober interface {
	Ping(ctx context.Context) error
}


type HealthHandlers struct {
	Postgres Prober
	Redis    Prober
	Kafka    Prober
}


type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}


type readyResponse struct {
	Status      string            `json:"status"`
	Checks      map[string]string `json:"checks"`
	LatenciesMS map[string]int64  `json:"latencies_ms"`
}


func (h *HealthHandlers) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Status:  "ok",
		Service: "backend",
		Version: version,
	})
}



func (h *HealthHandlers) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	checks := make(map[string]string, 3)
	latencies := make(map[string]int64, 3)

	probe := func(name string, p Prober) {
		if p == nil {
			checks[name] = "not configured"
			return
		}
		start := time.Now()
		if err := p.Ping(ctx); err != nil {
			checks[name] = "error: " + err.Error()
		} else {
			checks[name] = "ok"
		}
		latencies[name] = time.Since(start).Milliseconds()
	}

	probe("postgres", h.Postgres)
	probe("redis", h.Redis)
	probe("kafka", h.Kafka)

	status := "ready"
	httpStatus := http.StatusOK
	for _, v := range checks {
		if v != "ok" && v != "not configured" {
			status = "degraded"
			httpStatus = http.StatusServiceUnavailable
			break
		}
	}

	writeJSON(w, httpStatus, readyResponse{
		Status:      status,
		Checks:      checks,
		LatenciesMS: latencies,
	})
}


func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}


func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
		"code":  code,
	})
}


func requestID(r *http.Request) string {
	return middleware.GetRequestID(r.Context())
}



func decodeJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "malformed request body: "+err.Error())
		return false
	}
	return true
}
