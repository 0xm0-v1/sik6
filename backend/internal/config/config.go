package config

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultCORSOrigin   = "http://localhost:4200"
	defaultRecipesURL   = "http://localhost:4200/api/recipes"
	defaultSmokeTimeout = 10 * time.Second
)

var (
	defaultCORSMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}
	defaultCORSHeaders = []string{"Authorization", "Content-Type", "Accept", "X-Requested-With"}
	defaultExpose      = []string{"Content-Length", "Content-Type"}
	defaultRecipes     = []string{"Classic Pancakes", "Spicy Ramen", "Veggie Tacos"}
)

// Config centralises application configuration derived from the environment.
type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	API         APIConfig
	CORS        CORSConfig
	Smoke       SmokeConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host              string
	Port              int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
}

// Addr returns host:port suitable for http.Server.Addr.
func (s ServerConfig) Addr() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
}

// DatabaseConfig contains connectivity and pool parameters for Postgres.
type DatabaseConfig struct {
	DSN               string
	MaxConns          int
	MinConns          int
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

// APIConfig groups authentication-related settings.
type APIConfig struct {
	Token string
}

// CORSConfig captures cross-origin access policy.
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
}

// SmokeConfig collects CLI smoke-test parameters.
type SmokeConfig struct {
	RecipesURL string
	Expected   []string
	Timeout    time.Duration
}

// LoadSmoke builds SmokeConfig from environment variables without validating
// unrelated configuration sections.
func LoadSmoke() (*SmokeConfig, error) {
	smoke := SmokeConfig{
		RecipesURL: GetEnv("SMOKE_RECIPES_URL", defaultRecipesURL),
		Expected:   append([]string(nil), defaultRecipes...),
		Timeout:    GetEnvDuration("SMOKE_TIMEOUT", defaultSmokeTimeout),
	}
	if expected := splitAndClean(GetEnv("SMOKE_EXPECTED_RECIPES", "")); len(expected) > 0 {
		smoke.Expected = expected
	}

	return &smoke, nil
}

// Load builds Config from environment variables and validates critical fields.
func Load() (*Config, error) {
	env := strings.TrimSpace(GetEnv("ENV", ""))
	if env == "" {
		env = "dev"
	}

	server := ServerConfig{
		Host:              GetEnv("HOST", "0.0.0.0"),
		Port:              GetEnvInt("PORT", 8080),
		ReadTimeout:       GetEnvDuration("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:      GetEnvDuration("WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:       GetEnvDuration("IDLE_TIMEOUT", 60*time.Second),
		ReadHeaderTimeout: GetEnvDuration("READ_HEADER_TIMEOUT", 5*time.Second),
		ShutdownTimeout:   GetEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}

	cors := CORSConfig{
		AllowedOrigins:   splitAndClean(GetEnv("CORS_ALLOWED_ORIGINS", defaultCORSOrigin)),
		AllowedMethods:   defaultCORSMethods,
		AllowedHeaders:   defaultCORSHeaders,
		ExposedHeaders:   defaultExpose,
		AllowCredentials: strings.EqualFold(GetEnv("CORS_ALLOW_CREDENTIALS", "false"), "true"),
	}
	if methods := splitAndClean(GetEnv("CORS_ALLOWED_METHODS", "")); len(methods) > 0 {
		cors.AllowedMethods = methods
	}
	if headers := splitAndClean(GetEnv("CORS_ALLOWED_HEADERS", "")); len(headers) > 0 {
		cors.AllowedHeaders = headers
	}
	if exposed := splitAndClean(GetEnv("CORS_EXPOSED_HEADERS", "")); len(exposed) > 0 {
		cors.ExposedHeaders = exposed
	}
	if len(cors.AllowedOrigins) == 0 {
		return nil, fmt.Errorf("config: no CORS allowed origins configured")
	}

	api := APIConfig{
		Token: strings.TrimSpace(GetEnv("API_TOKEN", "")),
	}

	db := DatabaseConfig{
		DSN:               strings.TrimSpace(GetEnv("DB_DSN", "")),
		MaxConns:          GetEnvInt("DB_MAX_CONNS", 0),
		MinConns:          GetEnvInt("DB_MIN_CONNS", 0),
		MaxConnLifetime:   GetEnvDuration("DB_MAX_CONN_LIFETIME", 0),
		MaxConnIdleTime:   GetEnvDuration("DB_MAX_CONN_IDLE_TIME", 0),
		HealthCheckPeriod: GetEnvDuration("DB_HEALTH_CHECK_PERIOD", 0),
	}
	if db.DSN == "" {
		host := strings.TrimSpace(GetEnv("DB_HOST", "localhost"))
		if host == "" {
			host = "localhost"
		}
		user := strings.TrimSpace(GetEnv("DB_USER", "postgres"))
		if user == "" {
			user = "postgres"
		}
		name := strings.TrimSpace(GetEnv("DB_NAME", "postgres"))
		if name == "" {
			name = "postgres"
		}
		password := strings.TrimSpace(GetEnv("DB_PASSWORD", ""))
		sslMode := strings.TrimSpace(GetEnv("DB_SSLMODE", "disable"))
		port := GetEnvInt("DB_PORT", 5432)

		if host == "" || user == "" || name == "" {
			return nil, fmt.Errorf("config: DB_DSN or DB_HOST/DB_USER/DB_NAME must be configured")
		}

		dsn, err := buildPostgresDSN(host, port, user, password, name, sslMode)
		if err != nil {
			return nil, fmt.Errorf("config: %w", err)
		}
		db.DSN = dsn
	}

	smoke, err := LoadSmoke()
	if err != nil {
		return nil, fmt.Errorf("config: load smoke: %w", err)
	}

	return &Config{
		Environment: env,
		Server:      server,
		Database:    db,
		API:         api,
		CORS:        cors,
		Smoke:       *smoke,
	}, nil
}

func splitAndClean(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func buildPostgresDSN(host string, port int, user, password, database, sslmode string) (string, error) {
	if host == "" {
		return "", fmt.Errorf("db host is required")
	}
	if user == "" {
		return "", fmt.Errorf("db user is required")
	}
	if database == "" {
		return "", fmt.Errorf("db name is required")
	}
	if port <= 0 {
		port = 5432
	}

	u := &url.URL{
		Scheme: "postgres",
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
		Path:   "/" + database,
	}

	if password != "" {
		u.User = url.UserPassword(user, password)
	} else {
		u.User = url.User(user)
	}

	values := url.Values{}
	if sslmode != "" {
		values.Set("sslmode", sslmode)
	}
	if encoded := values.Encode(); encoded != "" {
		u.RawQuery = encoded
	}

	return u.String(), nil
}
