package internal

import (
	"encoding/json"
	"log"
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

	log.Printf("Received power: %.2f W", request.Power)

	now := time.Now()

	shouldNotify := false
	runningTime := time.Duration(0)

	h.Monitor.mu.Lock()

	if request.Power <= 0 {

		log.Println("Appliance stopped")

		h.Monitor.State.Running = false
		h.Monitor.State.StartTime = time.Time{}
		h.Monitor.State.Notified = false

	} else {

		if !h.Monitor.State.Running {

			log.Println("Appliance started")

			h.Monitor.State.Running = true
			h.Monitor.State.StartTime = now
			h.Monitor.State.Notified = false
		}

		if h.Monitor.State.Running &&
			!h.Monitor.State.Notified {

			runningTime = now.Sub(h.Monitor.State.StartTime)

			if runningTime >= h.Config.Threshold {

				log.Printf(
					"Threshold exceeded after %v",
					runningTime.Round(time.Second),
				)

				h.Monitor.State.Notified = true
				shouldNotify = true
			}
		}
	}

	currentState := h.Monitor.State

	h.Monitor.mu.Unlock()

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
