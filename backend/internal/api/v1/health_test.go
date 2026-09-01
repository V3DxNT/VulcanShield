package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)


type mockProber struct {
	err error
}

func (m *mockProber) Ping(_ context.Context) error {
	return m.err
}

func TestHealth_AlwaysOK(t *testing.T) {
	h := &HealthHandlers{}
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	h.Health(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var body healthResponse
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if body.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", body.Status)
	}
	if body.Service != "backend" {
		t.Errorf("expected service 'backend', got %q", body.Service)
	}
}

func TestReady_AllHealthy(t *testing.T) {
	h := &HealthHandlers{
		Postgres: &mockProber{},
		Redis:    &mockProber{},
		Kafka:    &mockProber{},
	}

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rr := httptest.NewRecorder()

	h.Ready(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var body readyResponse
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if body.Status != "ready" {
		t.Errorf("expected status 'ready', got %q", body.Status)
	}
}

func TestReady_DegradedWhenDependencyFails(t *testing.T) {
	h := &HealthHandlers{
		Postgres: &mockProber{},
		Redis:    &mockProber{err: context.DeadlineExceeded},
		Kafka:    &mockProber{},
	}

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rr := httptest.NewRecorder()

	h.Ready(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rr.Code)
	}
	var body readyResponse
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if body.Status != "degraded" {
		t.Errorf("expected status 'degraded', got %q", body.Status)
	}
	if body.Checks["redis"] == "ok" {
		t.Error("expected redis check to report error")
	}
}
