package config

import (
	"net"
	"os"
	"strconv"
	"time"
)

// Config holds server and process settings.
type Config struct {
	Host              string
	Port              int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
}

// New returns a Config initialized from environment variables with defaults.
func New() *Config {
	return &Config{
		Host:              getenv("HOST", "0.0.0.0"),
		Port:              getenvInt("PORT", 8080),
		ReadTimeout:       getenvDuration("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:      getenvDuration("WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:       getenvDuration("IDLE_TIMEOUT", 60*time.Second),
		ReadHeaderTimeout: getenvDuration("READ_HEADER_TIMEOUT", 5*time.Second),
		ShutdownTimeout:   getenvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

// Addr returns host:port suitable for http.Server.Addr.
func (c *Config) Addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getenvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
