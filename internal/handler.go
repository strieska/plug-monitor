package internal

import (
	"encoding/json"
	"net/http"
)

type PowerRequest struct {
	Power float64 `json:"power"`
}

func PowerHandler(w http.ResponseWriter, r *http.Request) {

	var request PowerRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"received_power": request.Power,
	})
}
