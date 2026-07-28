package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"plug-monitor/internal"
	"syscall"
	"time"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func main() {

	config := internal.LoadConfig()

	monitor := internal.NewMonitor()

	handler := internal.NewHandler(
		monitor,
		config,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", health)
	mux.HandleFunc("/status", handler.Status)
	mux.HandleFunc("/power", handler.Power)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           internal.Logging(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	go func() {

		log.Println("Server listening on :8080")

		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}

	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("Shutdown requested")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped")
}
