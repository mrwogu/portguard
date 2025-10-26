package main

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testAuthHeaderBasic              = "Basic "
	testExpectedStatusUnauthorizedFmt = "Expected status Unauthorized, got %d"
)

func TestBasicAuthMiddlewareDisabled(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Auth: AuthConfig{
				Enabled: false,
			},
		},
	}

	handler := basicAuthMiddleware(cfg, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}
	if w.Body.String() != "success" {
		t.Errorf("Expected 'success', got %s", w.Body.String())
	}
}

func TestBasicAuthMiddlewareNoCredentialsConfigured(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Auth: AuthConfig{
				Enabled:  true,
				Username: "",
				Password: "",
			},
		},
	}

	handler := basicAuthMiddleware(cfg, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	// Should pass through if no credentials configured
	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK when no credentials configured, got %d", w.Code)
	}
}

func TestBasicAuthMiddlewareMissingAuth(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Auth: AuthConfig{
				Enabled:  true,
				Username: "admin",
				Password: "password",
			},
		},
	}

	handler := basicAuthMiddleware(cfg, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf(testExpectedStatusUnauthorizedFmt, w.Code)
	}

	if w.Header().Get("WWW-Authenticate") == "" {
		t.Error("Expected WWW-Authenticate header")
	}
}

func TestBasicAuthMiddlewareInvalidCredentials(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Auth: AuthConfig{
				Enabled:  true,
				Username: "admin",
				Password: "password",
			},
		},
	}

	handler := basicAuthMiddleware(cfg, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", testAuthHeaderBasic+base64.StdEncoding.EncodeToString([]byte("wrong:credentials")))
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf(testExpectedStatusUnauthorizedFmt, w.Code)
	}
}

func TestBasicAuthMiddlewareValidCredentials(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Auth: AuthConfig{
				Enabled:  true,
				Username: "admin",
				Password: "password",
			},
		},
	}

	handler := basicAuthMiddleware(cfg, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", testAuthHeaderBasic+base64.StdEncoding.EncodeToString([]byte("admin:password")))
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}
	if w.Body.String() != "success" {
		t.Errorf("Expected 'success', got %s", w.Body.String())
	}
}

func TestBasicAuthMiddlewareWrongUsername(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Auth: AuthConfig{
				Enabled:  true,
				Username: "admin",
				Password: "password",
			},
		},
	}

	handler := basicAuthMiddleware(cfg, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", testAuthHeaderBasic+base64.StdEncoding.EncodeToString([]byte("wronguser:password")))
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf(testExpectedStatusUnauthorizedFmt, w.Code)
	}
}

func TestBasicAuthMiddlewareWrongPassword(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Auth: AuthConfig{
				Enabled:  true,
				Username: "admin",
				Password: "password",
			},
		},
	}

	handler := basicAuthMiddleware(cfg, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", testAuthHeaderBasic+base64.StdEncoding.EncodeToString([]byte("admin:wrongpassword")))
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf(testExpectedStatusUnauthorizedFmt, w.Code)
	}
}
