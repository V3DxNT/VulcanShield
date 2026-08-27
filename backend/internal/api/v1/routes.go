package v1

import "net/http"

// Mount registers all /api/v1/* routes onto mux.
// Health routes (/health, /ready) are mounted directly on the root mux
// by the router — not here.
//
// Future route stubs (Phase 4+):
//
//	POST   /api/v1/transactions
//	GET    /api/v1/transactions/{id}
//	GET    /api/v1/transactions
//	POST   /api/v1/scenarios/{type}/start
//	POST   /api/v1/scenarios/{id}/stop
//	GET    /api/v1/scenarios/{id}/status
//	WS     /api/v1/ws
//	POST   /api/v1/challenges/{id}/verify
//	GET    /api/v1/investigations/{id}
func Mount(mux *http.ServeMux, h *HealthHandlers) {
	// Phase 3 only exposes health/readiness on the root mux (see server/router.go).
	// This function is the single registration point for future /api/v1/ routes.
	// Each subsequent phase will add its handler registrations here.
	_ = mux
	_ = h
}
