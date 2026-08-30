package v1

import "net/http"

// Handlers holds all Phase 4+ HTTP handler groups.
type Handlers struct {
	Scenarios    *ScenarioHandlers
	Transactions *TransactionHandlers
}

// Mount registers all /api/v1/* routes onto mux.
// Health routes (/health, /ready) are mounted directly on the root mux
// by the router — not here.
func Mount(mux *http.ServeMux, h *Handlers) {
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

	// Future route stubs (Phase 5+):
	//   WS     /api/v1/ws
	//   POST   /api/v1/challenges/{id}/verify
	//   GET    /api/v1/investigations/{id}
}
