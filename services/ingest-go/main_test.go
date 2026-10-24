package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestEventIngestion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := newRouter()

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid login event",
			body:       `{"event_type":"login_failed","source_ip":"192.0.2.10","username":"alice","timestamp":"2026-10-24T02:00:00Z"}`,
			wantStatus: http.StatusAccepted,
		},
		{
			name:       "unsupported event type",
			body:       `{"event_type":"unknown","source_ip":"192.0.2.10","username":"alice","timestamp":"2026-10-24T02:00:00Z"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid timestamp",
			body:       `{"event_type":"login_failed","source_ip":"192.0.2.10","username":"alice","timestamp":"yesterday"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing username",
			body:       `{"event_type":"login_failed","source_ip":"192.0.2.10","timestamp":"2026-10-24T02:00:00Z"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/events",
				strings.NewReader(tt.body),
			)
			request.Header.Set("Content-Type", "application/json")

			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Errorf(
					"got status %d, want %d; body: %s",
					response.Code,
					tt.wantStatus,
					response.Body.String(),
				)
			}
		})
	}
}
