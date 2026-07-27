package internal

import (
	"encoding/json"
	"net/http"
	"time"
)

type PowerRequest struct {
	Power float64 `json:"power"`
}

type Handler struct {
	Monitor *Monitor
	Config  Config
}

func NewHandler(monitor *Monitor, config Config) *Handler {
	return &Handler{
		Monitor: monitor,
		Config:  config,
	}
}

func (h *Handler) Power(w http.ResponseWriter, r *http.Request) {

	var request PowerRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	now := time.Now()

	h.Monitor.mu.Lock()
	defer h.Monitor.mu.Unlock()

	//
	// Power is zero -> appliance stopped
	//
	if request.Power <= 0 {

		h.Monitor.State.Running = false
		h.Monitor.State.StartTime = time.Time{}
		h.Monitor.State.Notified = false

	} else {

		//
		// First time seeing power -> start tracking
		//
		if !h.Monitor.State.Running {

			h.Monitor.State.Running = true
			h.Monitor.State.StartTime = now
			h.Monitor.State.Notified = false

		}

		//
		// Still running -> check timeout
		//
		if h.Monitor.State.Running &&
			!h.Monitor.State.Notified &&
			now.Sub(h.Monitor.State.StartTime) >= h.Config.Threshold {

			h.Monitor.State.Notified = true

			// TODO:
			// Send Home Assistant notification webhook here

		}
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(h.Monitor.State)
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {

	h.Monitor.mu.Lock()
	defer h.Monitor.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(h.Monitor.State)
}
