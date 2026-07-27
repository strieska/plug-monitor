package main

import (
	"log"
	"net/http"

	"plug-monitor/internal"
)

func main() {

	http.HandleFunc("/health", health)

	http.HandleFunc("/power", internal.PowerHandler)

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
