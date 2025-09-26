package health_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/0xm0-v1/sik6/internal/health"
)

type envelopeLivez struct {
	Status string                 `json:"status"`
	Error  string                 `json:"error,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

func TestNewLivenessHandler(t *testing.T) {
	t.Parallel()

	h := health.NewLivenessHandler() // http.Handler

	cases := []struct {
		name       string
		method     string
		wantStatus int
		wantJSON   bool
	}{
		{"GET_ok", http.MethodGet, http.StatusOK, true},
		{"HEAD_ok", http.MethodHead, http.StatusOK, false},
		{"POST_method_not_allowed", http.MethodPost, http.StatusMethodNotAllowed, false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, "/livez", nil)
			rec := httptest.NewRecorder()

			// http.Handler: call ServeHTTP, not h(...)
			h.ServeHTTP(rec, req)

			res := rec.Result()
			t.Cleanup(func() { _ = res.Body.Close() })

			if res.StatusCode != tc.wantStatus {
				t.Fatalf("status: got %d want %d", res.StatusCode, tc.wantStatus)
			}

			// For 405, the Allow header should be present.
			if tc.method == http.MethodPost && res.Header.Get("Allow") == "" {
				t.Fatalf("missing Allow header on 405")
			}

			if tc.wantJSON {
				if ct := res.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
					t.Fatalf("content-type: %q", ct)
				}
				var got envelopeLivez
				if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
					t.Fatalf("decode json: %v", err)
				}
				if got.Status != "ok" {
					t.Fatalf("status field: %+v", got)
				}
				// time is encoded as RFC3339Nano string under data.time
				ts, _ := got.Data["time"].(string)
				if ts == "" {
					t.Fatalf("missing data.time")
				}
				if parsed, err := time.Parse(time.RFC3339Nano, ts); err != nil {
					t.Fatalf("bad time format: %v", err)
				} else if parsed.After(time.Now().UTC().Add(5 * time.Second)) {
					t.Fatalf("time too new: %v", parsed)
				}
			} else if tc.method == http.MethodHead {
				// HEAD must have no body.
				if b, _ := io.ReadAll(res.Body); len(b) != 0 {
					t.Fatalf("HEAD should have empty body, got %d bytes", len(b))
				}
			}
		})
	}
}
