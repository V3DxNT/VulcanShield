package v1

import (
	"net/http"

	appws "github.com/vulcanshield/backend/internal/websocket"
)

// Handlers holds all Phase 4+ HTTP handler groups.
type Handlers struct {
	Scenarios      *ScenarioHandlers
	Transactions   *TransactionHandlers
	Challenges     *ChallengeHandlers
	Graph          *GraphHandlers
	Investigations *InvestigationHandlers
	WSHub          *appws.Hub
}

// Mount registers all /api/v1/* routes onto mux.
func Mount(mux *http.ServeMux, h *Handlers) {
	// ── Phase 13: WebSocket Endpoint ──────────────────────────────────────────
	if h.WSHub != nil {
		mux.Handle("/api/v1/ws", h.WSHub)
	}

	// ── Phase 4: Scenario Control ─────────────────────────────────────────────
	if h.Scenarios != nil {
		mux.HandleFunc("POST /api/v1/scenarios/start", h.Scenarios.Start)
		mux.HandleFunc("POST /api/v1/scenarios/stop", h.Scenarios.Stop)
		mux.HandleFunc("GET /api/v1/scenarios/status", h.Scenarios.Status)
	}

	// ── Phase 4: Transaction Queries ──────────────────────────────────────────
	if h.Transactions != nil {
		mux.HandleFunc("GET /api/v1/transactions/{id}", h.Transactions.GetByID)
		mux.HandleFunc("GET /api/v1/transactions", h.Transactions.List)
	}

	// ── Phase 10: Step-Up OTP Verification ───────────────────────────────────
	if h.Challenges != nil {
		mux.HandleFunc("GET /api/v1/challenges/{id}", h.Challenges.GetByID)
		mux.HandleFunc("POST /api/v1/challenges/{id}/verify", h.Challenges.Verify)
	}

	// ── Phase 11: Fraud Network Graph ─────────────────────────────────────────
	if h.Graph != nil {
		mux.HandleFunc("GET /api/v1/graph/relationships", h.Graph.ListRelationships)
		mux.HandleFunc("GET /api/v1/graph/neighbors/{user_id}", h.Graph.GetNeighbors)
	}

	// ── Phase 18: AI Investigations ───────────────────────────────────────────
	if h.Investigations != nil {
		mux.HandleFunc("GET /api/v1/investigations/{transaction_id}", h.Investigations.Investigate)
	}
}
