package config

import (
	"os"
	"strconv"
	"time"
)

// GetEnv retrieves the environment variable for the given key.
// If the variable is unset or empty, it returns the provided default value.
func GetEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// GetEnvInt retrieves the environment variable for the given key as an int.
// If the variable is unset, empty, or cannot be parsed as an integer,
// the provided default value is returned.
func GetEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

// GetEnvDuration retrieves the environment variable for the given key as a time.Duration.
// If the variable is unset, empty, or cannot be parsed as a valid duration string,
// the provided default value is returned.
//
// Valid duration strings follow Go's time.ParseDuration format, e.g. "300ms", "1.5h", "2h45m".
func GetEnvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
