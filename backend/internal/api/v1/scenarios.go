package v1

import (
	"errors"
	"net/http"

	"github.com/vulcanshield/backend/internal/generator"
	"github.com/vulcanshield/backend/internal/models"
)


type ScenarioHandlers struct {
	Engine *generator.Engine
}




func (h *ScenarioHandlers) Start(w http.ResponseWriter, r *http.Request) {
	var req models.ScenarioStartRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	run, err := h.Engine.Start(r.Context(), req)
	if err != nil {
		if errors.Is(err, generator.ErrScenarioRunning) {
			writeError(w, http.StatusConflict, "SCENARIO_RUNNING",
				"a scenario is already running — stop it before starting a new one")
			return
		}
		writeError(w, http.StatusInternalServerError, "START_FAILED", err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, run)
}



func (h *ScenarioHandlers) Stop(w http.ResponseWriter, r *http.Request) {
	run, err := h.Engine.Stop(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "STOP_FAILED", err.Error())
		return
	}
	if run == nil {
		writeJSON(w, http.StatusOK, map[string]string{"status": "idle"})
		return
	}
	writeJSON(w, http.StatusOK, run)
}



func (h *ScenarioHandlers) Status(w http.ResponseWriter, r *http.Request) {
	run := h.Engine.GetStatus()
	if run == nil {
		writeJSON(w, http.StatusOK, map[string]string{"status": "idle"})
		return
	}
	writeJSON(w, http.StatusOK, run)
}
