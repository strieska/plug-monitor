package main

import (
	"io"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {

		body, _ := io.ReadAll(r.Body)

		log.Println("received:")
		log.Println(string(body))

		w.WriteHeader(http.StatusOK)
	})

	log.Println("fake webhook listening on :9999")

	log.Fatal(
		http.ListenAndServe(":9999", nil),
	)
}
