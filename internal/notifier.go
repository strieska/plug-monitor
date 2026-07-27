package internal

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func NotifyHA(webhookURL string, duration time.Duration) {

	if webhookURL == "" {
		log.Println("HA webhook not configured")
		return
	}

	payload := map[string]interface{}{
		"event":           "plug_running_too_long",
		"running_minutes": int(duration.Minutes()),
	}

	body, err := json.Marshal(payload)

	if err != nil {
		log.Println("failed to create webhook payload:", err)
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		webhookURL,
		bytes.NewBuffer(body),
	)

	if err != nil {
		log.Println("failed to create webhook request:", err)
		return
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Println("HA webhook failed:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf(
			"HA webhook returned status %d",
			resp.StatusCode,
		)

		return
	}

	log.Println("HA notification sent")
}
