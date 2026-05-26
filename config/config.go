package config

import (
	"os"
	"strconv"
)

type Config struct {
	NATSURL     string
	NATSSubject string
	NATSStream  string
	HTTPAddr    string
	WindowSecs  int
	TopCacheTTL int
}

func Load() Config {
	return Config{
		NATSURL:     getEnv("NATS_URL", "nats://localhost:4222"),
		NATSSubject: getEnv("NATS_SUBJECT", "search.events"),
		NATSStream:  getEnv("NATS_STREAM", "SEARCH"),
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
		WindowSecs:  getEnvInt("WINDOW_SECS", 300), // sliding window size
		TopCacheTTL: getEnvInt("TOP_CACHE_TTL", 1), // how often do we recalculate top
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
