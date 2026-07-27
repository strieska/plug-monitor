package internal

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Threshold time.Duration
	Webhook   string
}

func LoadConfig() Config {

	minutes, err := strconv.Atoi(
		getEnv("THRESHOLD_MINUTES", "90"),
	)

	if err != nil {
		minutes = 90
	}

	return Config{
		Threshold: time.Duration(minutes) * time.Minute,
		Webhook:   getEnv("HA_WEBHOOK", ""),
	}
}

func getEnv(key string, fallback string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
