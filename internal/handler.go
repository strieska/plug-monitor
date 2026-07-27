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

	shouldNotify := false
	runningTime := time.Duration(0)

	//
	// Only state manipulation happens inside the lock
	//
	h.Monitor.mu.Lock()

	if request.Power <= 0 {

		// Appliance stopped
		h.Monitor.State.Running = false
		h.Monitor.State.StartTime = time.Time{}
		h.Monitor.State.Notified = false

	} else {

		// First time seeing power
		if !h.Monitor.State.Running {

			h.Monitor.State.Running = true
			h.Monitor.State.StartTime = now
			h.Monitor.State.Notified = false
		}

		// Check if running too long
		if h.Monitor.State.Running &&
			!h.Monitor.State.Notified {

			runningTime = now.Sub(h.Monitor.State.StartTime)

			if runningTime >= h.Config.Threshold {

				h.Monitor.State.Notified = true
				shouldNotify = true
			}
		}
	}

	currentState := h.Monitor.State

	h.Monitor.mu.Unlock()

	//
	// External calls happen AFTER unlocking
	//
	if shouldNotify {

		go NotifyHA(
			h.Config.Webhook,
			runningTime,
		)

	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(currentState)
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {

	h.Monitor.mu.Lock()
	defer h.Monitor.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(h.Monitor.State)
}
