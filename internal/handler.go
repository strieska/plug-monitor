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
}

func NewHandler(monitor *Monitor) *Handler {
	return &Handler{
		Monitor: monitor,
	}
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {

	h.Monitor.mu.Lock()
	defer h.Monitor.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(h.Monitor.State)
}

func (h *Handler) Power(w http.ResponseWriter, r *http.Request) {

	var request PowerRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	h.Monitor.mu.Lock()
	defer h.Monitor.mu.Unlock()

	now := time.Now()

	if request.Power <= 0 {

		h.Monitor.State.Running = false
		h.Monitor.State.StartTime = time.Time{}
		h.Monitor.State.Notified = false

	} else {

		if !h.Monitor.State.Running {

			h.Monitor.State.Running = true
			h.Monitor.State.StartTime = now
			h.Monitor.State.Notified = false
		}
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(h.Monitor.State)
}
