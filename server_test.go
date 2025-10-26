package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mockExit is a test helper that captures exit codes
type mockExit struct {
	called   bool
	exitCode int
}

func (m *mockExit) exit(code int) {
	m.called = true
	m.exitCode = code
}

// mockServerStarter is a test helper that simulates server startup
type mockServerStarter struct {
	called     bool
	addr       string
	handler    http.Handler
	returnErr  error
}

func (m *mockServerStarter) start(addr string, handler http.Handler) error {
	m.called = true
	m.addr = addr
	m.handler = handler
	return m.returnErr
}

func TestRunVersionFlag(t *testing.T) {
	// Capture output by temporarily redirecting stdout
	oldAppVersion := appVersion
	appVersion = "test-version"
	defer func() { appVersion = oldAppVersion }()

	mockExit := &mockExit{}
	mockServer := &mockServerStarter{}
	args := []string{"--version"}

	err := run(args, mockExit.exit, mockServer.start)
	if err != nil {
		t.Errorf("run() returned error: %v", err)
	}

	if !mockExit.called {
		t.Error("Expected exit to be called for --version flag")
	}

	if mockExit.exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", mockExit.exitCode)
	}

	if mockServer.called {
		t.Error("Server should not be started for --version flag")
	}
}

func TestRunMissingConfig(t *testing.T) {
	mockExit := &mockExit{}
	mockServer := &mockServerStarter{}
	args := []string{"--config", "/nonexistent/path/to/config.yaml"}

	err := run(args, mockExit.exit, mockServer.start)
	if err == nil {
		t.Error("Expected error for missing config file")
	}

	if !strings.Contains(err.Error(), "error loading configuration") {
		t.Errorf("Expected config loading error, got: %v", err)
	}

	if mockServer.called {
		t.Error("Server should not be started when config is missing")
	}
}

func TestRunNoChecks(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "empty-config.yaml")

	emptyConfig := `server:
  port: "9999"
  timeout: 2s
checks: []
`
	if err := os.WriteFile(configPath, []byte(emptyConfig), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	mockExit := &mockExit{}
	mockServer := &mockServerStarter{}
	args := []string{"--config", configPath}

	err := run(args, mockExit.exit, mockServer.start)
	if err == nil {
		t.Error("Expected error when no checks are configured")
	}

	if !strings.Contains(err.Error(), "no port checks configured") {
		t.Errorf("Expected 'no port checks configured' error, got: %v", err)
	}

	if mockServer.called {
		t.Error("Server should not be started when no checks configured")
	}
}

func TestRunInvalidFlags(t *testing.T) {
	mockExit := &mockExit{}
	mockServer := &mockServerStarter{}
	args := []string{"--invalid-flag"}

	err := run(args, mockExit.exit, mockServer.start)
	if err == nil {
		t.Error("Expected error for invalid flag")
	}

	if mockServer.called {
		t.Error("Server should not be started with invalid flags")
	}
}

func TestRunWithAuthEnabled(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "auth-config.yaml")

	authConfig := `server:
  port: "8888"
  timeout: 1s
  auth:
    enabled: true
    username: "testuser"
    password: "testpass"
checks:
  - host: "localhost"
    port: 80
    name: "test-check"
`
	if err := os.WriteFile(configPath, []byte(authConfig), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	mockExit := &mockExit{}
	mockServer := &mockServerStarter{}
	args := []string{"--config", configPath}

	err := run(args, mockExit.exit, mockServer.start)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !mockServer.called {
		t.Error("Expected server to be started")
	}

	if mockServer.addr != ":8888" {
		t.Errorf("Expected server to listen on :8888, got %s", mockServer.addr)
	}
}

func TestRunWithAuthDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "no-auth-config.yaml")

	noAuthConfig := `server:
  port: "9999"
  timeout: 1s
  auth:
    enabled: false
checks:
  - host: "localhost"
    port: 80
    name: "test-check"
`
	if err := os.WriteFile(configPath, []byte(noAuthConfig), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	mockExit := &mockExit{}
	mockServer := &mockServerStarter{}
	args := []string{"--config", configPath}

	err := run(args, mockExit.exit, mockServer.start)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !mockServer.called {
		t.Error("Expected server to be started")
	}

	if mockServer.addr != ":9999" {
		t.Errorf("Expected server to listen on :9999, got %s", mockServer.addr)
	}
}

func TestRunServerStartError(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "valid-config.yaml")

	validConfig := `server:
  port: "8080"
  timeout: 1s
checks:
  - host: "localhost"
    port: 80
    name: "test-check"
`
	if err := os.WriteFile(configPath, []byte(validConfig), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	mockExit := &mockExit{}
	mockServer := &mockServerStarter{
		returnErr: fmt.Errorf("bind: address already in use"),
	}
	args := []string{"--config", configPath}

	err := run(args, mockExit.exit, mockServer.start)
	if err == nil {
		t.Error("Expected error when server fails to start")
	}

	if !strings.Contains(err.Error(), "server error") {
		t.Errorf("Expected 'server error', got: %v", err)
	}

	if !mockServer.called {
		t.Error("Expected server start to be attempted")
	}
}

func TestSetupAndStartServerSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	testConfig := `server:
  port: "7777"
  timeout: 2s
  auth:
    enabled: true
    username: "admin"
    password: "secret"
checks:
  - host: "localhost"
    port: 443
    name: "HTTPS"
`
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	cfg, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	mockServer := &mockServerStarter{}

	err = setupAndStartServer(cfg, configPath, mockServer.start)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !mockServer.called {
		t.Error("Expected server to be started")
	}

	if mockServer.addr != ":7777" {
		t.Errorf("Expected server to listen on :7777, got %s", mockServer.addr)
	}
}

func TestSetupAndStartServerAuthDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	testConfig := `server:
  port: "6666"
  timeout: 2s
  auth:
    enabled: false
checks:
  - host: "localhost"
    port: 22
    name: "SSH"
`
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	cfg, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	mockServer := &mockServerStarter{}

	err = setupAndStartServer(cfg, configPath, mockServer.start)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !mockServer.called {
		t.Error("Expected server to be started")
	}
}

func TestSetupAndStartServerError(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	testConfig := `server:
  port: "5555"
  timeout: 1s
checks:
  - host: "localhost"
    port: 80
    name: "HTTP"
`
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	cfg, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	mockServer := &mockServerStarter{
		returnErr: fmt.Errorf("cannot bind to port"),
	}

	err = setupAndStartServer(cfg, configPath, mockServer.start)
	if err == nil {
		t.Error("Expected error when server fails to start")
	}

	if !strings.Contains(err.Error(), "server error") {
		t.Errorf("Expected 'server error', got: %v", err)
	}
}
