package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		configData  string
		wantErr     bool
		wantPort    string
		wantTimeout time.Duration
		wantChecks  int
	}{
		{
			name: "valid config with custom values",
			configData: `
server:
  port: "9000"
  timeout: 5s
checks:
  - host: "localhost"
    port: 8080
    name: "Web Server"
    description: "Main web server"
  - host: "db.example.com"
    port: 5432
    name: "Database"
    description: "PostgreSQL database"
`,
			wantErr:     false,
			wantPort:    "9000",
			wantTimeout: 5 * time.Second,
			wantChecks:  2,
		},
		{
			name: "config with default values",
			configData: `
checks:
  - host: "localhost"
    port: 8080
    name: "Test Service"
`,
			wantErr:     false,
			wantPort:    defaultListenPort,
			wantTimeout: 2 * time.Second,
			wantChecks:  1,
		},
		{
			name: "minimal valid config",
			configData: `
checks:
  - host: "example.com"
    port: 80
    name: "Website"
`,
			wantErr:     false,
			wantPort:    defaultListenPort,
			wantTimeout: 2 * time.Second,
			wantChecks:  1,
		},
		{
			name:       "invalid YAML syntax",
			configData: `invalid: yaml: [syntax`,
			wantErr:    true,
		},
		{
			name:        "empty config file",
			configData:  ``,
			wantErr:     false,
			wantPort:    defaultListenPort,
			wantTimeout: 2 * time.Second,
			wantChecks:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			if err := os.WriteFile(configPath, []byte(tt.configData), 0644); err != nil {
				t.Fatalf("Failed to create test config file: %v", err)
			}

			// Load config
			cfg, err := loadConfig(configPath)

			// Check error expectation
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify config values
			if cfg.Server.Port != tt.wantPort {
				t.Errorf("Port = %q, want %q", cfg.Server.Port, tt.wantPort)
			}

			if cfg.Server.Timeout != tt.wantTimeout {
				t.Errorf("Timeout = %v, want %v", cfg.Server.Timeout, tt.wantTimeout)
			}

			if len(cfg.Checks) != tt.wantChecks {
				t.Errorf("Number of checks = %d, want %d", len(cfg.Checks), tt.wantChecks)
			}
		})
	}
}

func TestLoadConfigNonExistentFile(t *testing.T) {
	_, err := loadConfig("/nonexistent/path/to/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file, got none")
	}
}

func TestLoadConfigDetailedValues(t *testing.T) {
	configData := `
server:
  port: "3000"
  timeout: 10s
checks:
  - host: "api.example.com"
    port: 443
    name: "API Server"
    description: "Production API"
  - host: "cache.example.com"
    port: 6379
    name: "Redis Cache"
    description: "Redis instance"
  - host: "queue.example.com"
    port: 5672
    name: "RabbitMQ"
    description: "Message queue"
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	cfg, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify server config
	if cfg.Server.Port != "3000" {
		t.Errorf("Port = %q, want '3000'", cfg.Server.Port)
	}
	if cfg.Server.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want 10s", cfg.Server.Timeout)
	}

	// Verify checks
	if len(cfg.Checks) != 3 {
		t.Fatalf("Number of checks = %d, want 3", len(cfg.Checks))
	}

	expectedChecks := []struct {
		host        string
		port        int
		name        string
		description string
	}{
		{"api.example.com", 443, "API Server", "Production API"},
		{"cache.example.com", 6379, "Redis Cache", "Redis instance"},
		{"queue.example.com", 5672, "RabbitMQ", "Message queue"},
	}

	for i, expected := range expectedChecks {
		check := cfg.Checks[i]
		if check.Host != expected.host {
			t.Errorf("Check[%d].Host = %q, want %q", i, check.Host, expected.host)
		}
		if check.Port != expected.port {
			t.Errorf("Check[%d].Port = %d, want %d", i, check.Port, expected.port)
		}
		if check.Name != expected.name {
			t.Errorf("Check[%d].Name = %q, want %q", i, check.Name, expected.name)
		}
		if check.Description != expected.description {
			t.Errorf("Check[%d].Description = %q, want %q", i, check.Description, expected.description)
		}
	}
}

func TestLoadConfigWithPerCheckTimeout(t *testing.T) {
	configData := `
server:
  port: "8888"
  timeout: 2s
checks:
  - host: "localhost"
    port: 8080
    name: "Default Timeout Service"
    description: "Uses server timeout"
  - host: "remote.example.com"
    port: 443
    name: "Slow Service"
    description: "Needs custom timeout"
    timeout: 10s
  - host: "fast.local"
    port: 6379
    name: "Fast Service"
    description: "Quick check"
    timeout: 500ms
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	cfg, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify server config
	if cfg.Server.Timeout != 2*time.Second {
		t.Errorf("Server.Timeout = %v, want 2s", cfg.Server.Timeout)
	}

	// Verify checks
	if len(cfg.Checks) != 3 {
		t.Fatalf("Number of checks = %d, want 3", len(cfg.Checks))
	}

	// First check: no custom timeout (should be 0)
	if cfg.Checks[0].Timeout != 0 {
		t.Errorf("Check[0].Timeout = %v, want 0 (uses server default)", cfg.Checks[0].Timeout)
	}

	// Second check: custom timeout of 10s
	if cfg.Checks[1].Timeout != 10*time.Second {
		t.Errorf("Check[1].Timeout = %v, want 10s", cfg.Checks[1].Timeout)
	}

	// Third check: custom timeout of 500ms
	if cfg.Checks[2].Timeout != 500*time.Millisecond {
		t.Errorf("Check[2].Timeout = %v, want 500ms", cfg.Checks[2].Timeout)
	}
}
