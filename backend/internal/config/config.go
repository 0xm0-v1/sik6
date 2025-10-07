package config

import (
	"net"
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
		Host:              GetEnv("HOST", "0.0.0.0"),
		Port:              GetEnvInt("PORT", 8080),
		ReadTimeout:       GetEnvDuration("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:      GetEnvDuration("WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:       GetEnvDuration("IDLE_TIMEOUT", 60*time.Second),
		ReadHeaderTimeout: GetEnvDuration("READ_HEADER_TIMEOUT", 5*time.Second),
		ShutdownTimeout:   GetEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

// Addr returns host:port suitable for http.Server.Addr.
func (c *Config) Addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}
