package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHealthHandler(t *testing.T) {
	// Start a test server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer func() { _ = listener.Close() }()

	var testPort int
	_, portStr, _ := net.SplitHostPort(listener.Addr().String())
	_, _ = fmt.Sscanf(portStr, "%d", &testPort)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()

	tests := []struct {
		name           string
		config         *Config
		wantStatusCode int
		wantStatus     string
	}{
		{
			name: "healthy response",
			config: &Config{
				Server: ServerConfig{
					Port:    "8888",
					Timeout: 2 * time.Second,
				},
				Checks: []PortCheck{
					{
						Host:        "127.0.0.1",
						Port:        testPort,
						Name:        "Test Service",
						Description: "Test",
					},
				},
			},
			wantStatusCode: http.StatusOK,
			wantStatus:     "healthy",
		},
		{
			name: "unhealthy response",
			config: &Config{
				Server: ServerConfig{
					Port:    "8888",
					Timeout: 500 * time.Millisecond,
				},
				Checks: []PortCheck{
					{
						Host:        "192.0.2.1",
						Port:        80,
						Name:        "Unreachable",
						Description: "Should fail",
					},
				},
			},
			wantStatusCode: http.StatusServiceUnavailable,
			wantStatus:     "unhealthy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()

			handler := healthHandler(tt.config)
			handler(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("Status code = %d, want %d", rec.Code, tt.wantStatusCode)
			}

			contentType := rec.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Content-Type = %q, want 'application/json'", contentType)
			}

			var status HealthStatus
			if err := json.NewDecoder(rec.Body).Decode(&status); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if status.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q", status.Status, tt.wantStatus)
			}

			if len(status.Checks) != len(tt.config.Checks) {
				t.Errorf("Number of checks = %d, want %d", len(status.Checks), len(tt.config.Checks))
			}

			if status.Time == "" {
				t.Error("Timestamp is empty")
			}

			if status.Message == "" {
				t.Error("Message is empty")
			}
		})
	}
}

func TestHealthHandlerMultipleChecks(t *testing.T) {
	// Start two test servers
	listener1, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start test server 1: %v", err)
	}
	defer func() { _ = listener1.Close() }()

	listener2, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start test server 2: %v", err)
	}
	defer func() { _ = listener2.Close() }()

	var port1, port2 int
	_, portStr1, _ := net.SplitHostPort(listener1.Addr().String())
	_, portStr2, _ := net.SplitHostPort(listener2.Addr().String())
	_, _ = fmt.Sscanf(portStr1, "%d", &port1)
	_, _ = fmt.Sscanf(portStr2, "%d", &port2)

	go func() {
		for {
			conn, err := listener1.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()
	go func() {
		for {
			conn, err := listener2.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()

	cfg := &Config{
		Server: ServerConfig{
			Port:    "8888",
			Timeout: 2 * time.Second,
		},
		Checks: []PortCheck{
			{Host: "127.0.0.1", Port: port1, Name: "Service 1"},
			{Host: "127.0.0.1", Port: port2, Name: "Service 2"},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler := healthHandler(cfg)
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}

	var status HealthStatus
	if err := json.NewDecoder(rec.Body).Decode(&status); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(status.Checks) != 2 {
		t.Errorf("Number of checks = %d, want 2", len(status.Checks))
	}

	for i, check := range status.Checks {
		if check.Status != "healthy" {
			t.Errorf("Check[%d] status = %q, want 'healthy'", i, check.Status)
		}
	}
}

func TestLiveHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/live", nil)
	rec := httptest.NewRecorder()

	liveHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("Content-Type = %q, want 'text/plain'", contentType)
	}

	body := strings.TrimSpace(rec.Body.String())
	if body != "OK" {
		t.Errorf("Body = %q, want 'OK'", body)
	}
}

func TestLiveHandlerAlwaysHealthy(t *testing.T) {
	// Test that /live always returns OK regardless of config
	req := httptest.NewRequest(http.MethodGet, "/live", nil)
	rec := httptest.NewRecorder()

	liveHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}

	// Run multiple times to ensure consistency
	for i := 0; i < 5; i++ {
		req = httptest.NewRequest(http.MethodGet, "/live", nil)
		rec = httptest.NewRecorder()
		liveHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Iteration %d: Status code = %d, want %d", i, rec.Code, http.StatusOK)
		}
	}
}

func TestRootHandler(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Port: "8888"},
		Checks: []PortCheck{
			{Host: "localhost", Port: 8080, Name: "Test"},
			{Host: "localhost", Port: 8081, Name: "Test2"},
		},
	}

	tests := []struct {
		name            string
		path            string
		wantStatusCode  int
		wantContentType string
		checkBody       bool
	}{
		{
			name:            "root path returns HTML",
			path:            "/",
			wantStatusCode:  http.StatusOK,
			wantContentType: "text/html",
			checkBody:       true,
		},
		{
			name:            "non-root path returns 404",
			path:            "/notfound",
			wantStatusCode:  http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
			checkBody:       false,
		},
		{
			name:            "nested path returns 404",
			path:            "/some/nested/path",
			wantStatusCode:  http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
			checkBody:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			handler := rootHandler(cfg)
			handler(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("Status code = %d, want %d", rec.Code, tt.wantStatusCode)
			}

			contentType := rec.Header().Get("Content-Type")
			if contentType != tt.wantContentType {
				t.Errorf("Content-Type = %q, want %q", contentType, tt.wantContentType)
			}

			if tt.checkBody {
				body := rec.Body.String()

				expectedStrings := []string{
					"PortGuard",
					"Health Check Service",
					"/health",
					"/live",
					"2 ports", // Based on config above
					appVersion,
				}

				for _, expected := range expectedStrings {
					if !strings.Contains(body, expected) {
						t.Errorf("Body does not contain %q", expected)
					}
				}

				// Check HTML structure
				if !strings.Contains(body, "<!DOCTYPE html>") {
					t.Error("Body does not contain HTML doctype")
				}
				if !strings.Contains(body, "<html>") {
					t.Error("Body does not contain <html> tag")
				}
			}
		})
	}
}

func TestRootHandlerPortCount(t *testing.T) {
	tests := []struct {
		name       string
		checkCount int
		wantText   string
	}{
		{
			name:       "zero ports",
			checkCount: 0,
			wantText:   "0 ports",
		},
		{
			name:       "one port",
			checkCount: 1,
			wantText:   "1 ports",
		},
		{
			name:       "multiple ports",
			checkCount: 5,
			wantText:   "5 ports",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checks := make([]PortCheck, tt.checkCount)
			for i := 0; i < tt.checkCount; i++ {
				checks[i] = PortCheck{
					Host: "localhost",
					Port: 8080 + i,
					Name: fmt.Sprintf("Service %d", i),
				}
			}

			cfg := &Config{
				Server: ServerConfig{Port: "8888"},
				Checks: checks,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			handler := rootHandler(cfg)
			handler(rec, req)

			body := rec.Body.String()
			if !strings.Contains(body, tt.wantText) {
				t.Errorf("Body does not contain %q", tt.wantText)
			}
		})
	}
}

func TestHandlersHTTPMethods(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Port: "8888", Timeout: 2 * time.Second},
		Checks: []PortCheck{
			{Host: "localhost", Port: 8080, Name: "Test"},
		},
	}

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	t.Run("health handler accepts all methods", func(t *testing.T) {
		for _, method := range methods {
			req := httptest.NewRequest(method, "/health", nil)
			rec := httptest.NewRecorder()

			handler := healthHandler(cfg)
			handler(rec, req)

			// Handler should accept all methods (it doesn't check method)
			if rec.Code != http.StatusOK && rec.Code != http.StatusServiceUnavailable {
				t.Errorf("Method %s: unexpected status code %d", method, rec.Code)
			}
		}
	})

	t.Run("live handler accepts all methods", func(t *testing.T) {
		for _, method := range methods {
			req := httptest.NewRequest(method, "/live", nil)
			rec := httptest.NewRecorder()

			liveHandler(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("Method %s: status code = %d, want %d", method, rec.Code, http.StatusOK)
			}
		}
	})

	t.Run("root handler accepts all methods", func(t *testing.T) {
		for _, method := range methods {
			req := httptest.NewRequest(method, "/", nil)
			rec := httptest.NewRecorder()

			handler := rootHandler(cfg)
			handler(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("Method %s: status code = %d, want %d", method, rec.Code, http.StatusOK)
			}
		}
	})
}
