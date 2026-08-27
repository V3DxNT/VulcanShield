package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vulcanshield/backend/internal/middleware"
)

const version = "0.1.0"

// Prober is satisfied by any dependency that can report its health.
type Prober interface {
	Ping(ctx context.Context) error
}

// HealthHandlers holds the dependencies needed for health/readiness endpoints.
type HealthHandlers struct {
	Postgres Prober
	Redis    Prober
	Kafka    Prober
}

// healthResponse is the liveness endpoint payload.
type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// readyResponse is the readiness endpoint payload.
type readyResponse struct {
	Status      string            `json:"status"`
	Checks      map[string]string `json:"checks"`
	LatenciesMS map[string]int64  `json:"latencies_ms"`
}

// Health handles GET /health — always 200 if the server is running.
func (h *HealthHandlers) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Status:  "ok",
		Service: "backend",
		Version: version,
	})
}

// Ready handles GET /ready — probes PostgreSQL, Redis, and Kafka.
// Returns 200 if all dependencies are healthy, 503 if any are unavailable.
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

// writeJSON writes a JSON-encoded value with the given HTTP status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes a structured JSON error response.
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
		"code":  code,
	})
}

// requestID is a convenience wrapper that reads the ID from context.
func requestID(r *http.Request) string {
	return middleware.GetRequestID(r.Context())
}
