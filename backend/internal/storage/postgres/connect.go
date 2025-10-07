package postgres

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/0xm0-v1/sik6/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect creates a pgx connection pool using environment variables and validates it with Ping.
func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := config.GetEnv("DB_DSN", "")
	if dsn == "" {
		var err error
		dsn, err = buildDSN(
			config.GetEnv("DB_HOST", "localhost"),
			config.GetEnvInt("DB_PORT", 5432),
			config.GetEnv("DB_USER", "postgres"),
			config.GetEnv("DB_PASSWORD", ""),
			config.GetEnv("DB_NAME", "postgres"),
			config.GetEnv("DB_SSLMODE", "disable"),
		)
		if err != nil {
			return nil, fmt.Errorf("build DSN: %w", err)
		}
	}

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	if maxConns := config.GetEnvInt("DB_MAX_CONNS", 0); maxConns > 0 {
		poolConfig.MaxConns = int32(maxConns)
	}
	if minConns := config.GetEnvInt("DB_MIN_CONNS", 0); minConns > 0 {
		poolConfig.MinConns = int32(minConns)
	}
	if lifetime := config.GetEnvDuration("DB_MAX_CONN_LIFETIME", 0); lifetime > 0 {
		poolConfig.MaxConnLifetime = lifetime
	}
	if idle := config.GetEnvDuration("DB_MAX_CONN_IDLE_TIME", 0); idle > 0 {
		poolConfig.MaxConnIdleTime = idle
	}
	if period := config.GetEnvDuration("DB_HEALTH_CHECK_PERIOD", 0); period > 0 {
		poolConfig.HealthCheckPeriod = period
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func buildDSN(host string, port int, user, password, database, sslmode string) (string, error) {
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
