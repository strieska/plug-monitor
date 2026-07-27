package main

import (
	"log"
	"net/http"

	"plug-monitor/internal"
)

func main() {

	config := internal.LoadConfig()

	monitor := internal.NewMonitor()

	handler := internal.NewHandler(
		monitor,
		config,
	)

	http.HandleFunc("/health", health)

	http.HandleFunc("/power", handler.Power)

	http.HandleFunc("/status", handler.Status)

	log.Println("plug-monitor starting on :8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
