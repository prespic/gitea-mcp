package operation

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithNoBodyContentTypeAddsContentTypeForAcceptedAndNoContent(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "accepted",
			status: http.StatusAccepted,
		},
		{
			name:   "no_content",
			status: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := withNoBodyContentType(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
			}))

			req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.status {
				t.Fatalf("expected status %d, got %d", tc.status, rr.Code)
			}
			if got := rr.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("expected Content-Type application/json, got %q", got)
			}
		})
	}
}

func TestWithNoBodyContentTypePreservesExistingContentType(t *testing.T) {
	handler := withNoBodyContentType(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusAccepted)
	}))

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); got != "text/plain" {
		t.Fatalf("expected Content-Type text/plain, got %q", got)
	}
}

func TestWithNoBodyContentTypeDoesNotModifyOtherStatusCodes(t *testing.T) {
	handler := withNoBodyContentType(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); got != "" {
		t.Fatalf("expected empty Content-Type, got %q", got)
	}
}
