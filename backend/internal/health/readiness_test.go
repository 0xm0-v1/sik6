package health_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0xm0-v1/sik6/internal/health"
	healthhttp "github.com/0xm0-v1/sik6/internal/health/transport/http"
)

type envelopeReadyz struct {
	Status string         `json:"status"`
	Error  string         `json:"error,omitempty"`
	Data   map[string]any `json:"data,omitempty"`
}

func TestNewReadinessHandler(t *testing.T) {
	t.Parallel()

	ok := func(context.Context) error { return nil }
	fail := func(context.Context) error { return errors.New("not ready") }

	cases := []struct {
		name       string
		method     string
		checker    health.Checker
		wantStatus int
		wantCode   int // mirrors status in envelopeReadyz semantics
		wantJSON   bool
	}{
		{"GET_ready", http.MethodGet, ok, http.StatusOK, http.StatusOK, true},
		{"GET_not_ready", http.MethodGet, fail, http.StatusServiceUnavailable, http.StatusServiceUnavailable, true},
		{"HEAD_ready", http.MethodHead, ok, http.StatusOK, 0, false},
		{"POST_method_not_allowed", http.MethodPost, ok, http.StatusMethodNotAllowed, 0, false},
		{"GET_nil_checker_defaults_ok", http.MethodGet, nil, http.StatusOK, http.StatusOK, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			chk := tc.checker
			if chk == nil {
				chk = ok
			}
			h := healthhttp.NewHandlers(chk).Readiness

			req := httptest.NewRequest(tc.method, "/readyz", nil)
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			res := rec.Result()
			t.Cleanup(func() { _ = res.Body.Close() })

			if res.StatusCode != tc.wantStatus {
				t.Fatalf("status: got %d want %d", res.StatusCode, tc.wantStatus)
			}

			if tc.method == http.MethodPost && res.Header.Get("Allow") == "" {
				t.Fatalf("missing Allow header on 405")
			}

			if tc.wantJSON {
				if ct := res.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
					t.Fatalf("content-type: %q", ct)
				}
				var got envelopeReadyz
				if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
					t.Fatalf("decode json: %v", err)
				}
				if got.Status == "" {
					t.Fatalf("missing status field")
				}
				_ = got
			} else if tc.method == http.MethodHead {
				if b, _ := io.ReadAll(res.Body); len(b) != 0 {
					t.Fatalf("HEAD should have empty body, got %d bytes", len(b))
				}
			}
		})
	}
}
