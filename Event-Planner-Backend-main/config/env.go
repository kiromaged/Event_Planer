package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads variables from a .env file if present.
// It is safe to call multiple times; subsequent calls are no-ops.
func LoadEnv() {
	_ = godotenv.Load()
}

// GetEnv returns the value for a key or a default if not set.
func GetEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

// MustGetEnv returns the value for a key or logs a warning and returns fallback.
func MustGetEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Printf("env %s not set, using default", key)
		return fallback
	}
	return v
}
