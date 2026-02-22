package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		method          string
		wantStatus      int
		wantContentType string
		wantJSONStatus  string
	}{
		{
			name:            "GET returns ok",
			method:          http.MethodGet,
			wantStatus:      http.StatusOK,
			wantContentType: "application/json",
			wantJSONStatus:  "ok",
		},
		{
			name:       "POST is not allowed",
			method:     http.MethodPost,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(tt.method, "/health", nil)
			rr := httptest.NewRecorder()

			healthHandler(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rr.Code, tt.wantStatus)
			}

			if tt.wantContentType != "" {
				if got := rr.Header().Get("Content-Type"); got != tt.wantContentType {
					t.Fatalf("content-type = %q, want %q", got, tt.wantContentType)
				}
			}

			if tt.wantJSONStatus != "" {
				var body map[string]string
				if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
					t.Fatalf("decode body: %v", err)
				}

				if got := body["status"]; got != tt.wantJSONStatus {
					t.Fatalf("body.status = %q, want %q", got, tt.wantJSONStatus)
				}
			}
		})
	}
}
