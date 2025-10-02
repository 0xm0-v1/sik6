package health

import (
	"errors"
	"time"
)

// --- Types --- //

// Response is the JSON payload returned by health endpoints.
type Response struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
	Code   int       `json:"code"`
}

// ErrNotReady is a possible sentinel error for readiness checks.
var ErrNotReady = errors.New("not ready")
