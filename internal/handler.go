package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

type PowerRequest struct {
	Power float64 `json:"power"`
}

type Handler struct {
	Monitor *Monitor
	Config  Config
}

func NewHandler(
	monitor *Monitor,
	config Config,
) *Handler {

	return &Handler{
		Monitor: monitor,
		Config:  config,
	}
}

func (h *Handler) Power(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request PowerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid json",
			http.StatusBadRequest,
		)
		return
	}

	log.Printf("Power: %.2f W", request.Power)

	state, notify, runtime := h.Monitor.Update(
		request.Power,
		h.Config.Threshold,
	)

	if notify {

		log.Printf(
			"Threshold exceeded after %v",
			runtime.Round(0),
		)

		go NotifyHA(
			h.Config.Webhook,
			runtime,
		)
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(state)
}

func (h *Handler) Status(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		h.Monitor.GetState(),
	)
}
