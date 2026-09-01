package v1

import (
	"net/http"

	appws "github.com/vulcanshield/backend/internal/websocket"
)


type Handlers struct {
	Scenarios      *ScenarioHandlers
	Transactions   *TransactionHandlers
	Challenges     *ChallengeHandlers
	Graph          *GraphHandlers
	Investigations *InvestigationHandlers
	WSHub          *appws.Hub
}


func Mount(mux *http.ServeMux, h *Handlers) {
	
	if h.WSHub != nil {
		mux.Handle("/api/v1/ws", h.WSHub)
	}

	
	if h.Scenarios != nil {
		mux.HandleFunc("POST /api/v1/scenarios/start", h.Scenarios.Start)
		mux.HandleFunc("POST /api/v1/scenarios/stop", h.Scenarios.Stop)
		mux.HandleFunc("GET /api/v1/scenarios/status", h.Scenarios.Status)
	}

	
	if h.Transactions != nil {
		mux.HandleFunc("GET /api/v1/transactions/{id}", h.Transactions.GetByID)
		mux.HandleFunc("GET /api/v1/transactions", h.Transactions.List)
	}

	
	if h.Challenges != nil {
		mux.HandleFunc("GET /api/v1/challenges/{id}", h.Challenges.GetByID)
		mux.HandleFunc("POST /api/v1/challenges/{id}/verify", h.Challenges.Verify)
	}

	
	if h.Graph != nil {
		mux.HandleFunc("GET /api/v1/graph/relationships", h.Graph.ListRelationships)
		mux.HandleFunc("GET /api/v1/graph/neighbors/{user_id}", h.Graph.GetNeighbors)
	}

	
	if h.Investigations != nil {
		mux.HandleFunc("GET /api/v1/investigations/{transaction_id}", h.Investigations.Investigate)
	}
}
