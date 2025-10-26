package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestConfigStructure(t *testing.T) {
	cfg := Config{
		Server: ServerConfig{
			Port:    "9000",
			Timeout: 5 * time.Second,
		},
		Checks: []PortCheck{
			{
				Host:        "example.com",
				Port:        443,
				Name:        "Example Service",
				Description: "A test service",
			},
		},
	}

	if cfg.Server.Port != "9000" {
		t.Errorf("Server.Port = %q, want '9000'", cfg.Server.Port)
	}

	if cfg.Server.Timeout != 5*time.Second {
		t.Errorf("Server.Timeout = %v, want 5s", cfg.Server.Timeout)
	}

	if len(cfg.Checks) != 1 {
		t.Fatalf("len(Checks) = %d, want 1", len(cfg.Checks))
	}

	check := cfg.Checks[0]
	if check.Host != "example.com" {
		t.Errorf("Check.Host = %q, want 'example.com'", check.Host)
	}
	if check.Port != 443 {
		t.Errorf("Check.Port = %d, want 443", check.Port)
	}
	if check.Name != "Example Service" {
		t.Errorf("Check.Name = %q, want 'Example Service'", check.Name)
	}
	if check.Description != "A test service" {
		t.Errorf("Check.Description = %q, want 'A test service'", check.Description)
	}
}

func TestHealthStatusJSON(t *testing.T) {
	status := HealthStatus{
		Status:  "healthy",
		Message: "All services running",
		Checks: []PortCheckResult{
			{
				Name:        "Web Server",
				Host:        "localhost",
				Port:        8080,
				Description: "Main web server",
				Status:      "healthy",
			},
		},
		Time:    "2024-01-15T10:30:00Z",
		Version: "1.0.0",
	}

	jsonData, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var decoded HealthStatus
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if decoded.Status != status.Status {
		t.Errorf("Status = %q, want %q", decoded.Status, status.Status)
	}

	if decoded.Message != status.Message {
		t.Errorf("Message = %q, want %q", decoded.Message, status.Message)
	}

	if decoded.Time != status.Time {
		t.Errorf("Time = %q, want %q", decoded.Time, status.Time)
	}

	if decoded.Version != status.Version {
		t.Errorf("Version = %q, want %q", decoded.Version, status.Version)
	}

	if len(decoded.Checks) != 1 {
		t.Fatalf("len(Checks) = %d, want 1", len(decoded.Checks))
	}

	check := decoded.Checks[0]
	if check.Name != "Web Server" {
		t.Errorf("Check.Name = %q, want 'Web Server'", check.Name)
	}
	if check.Status != "healthy" {
		t.Errorf("Check.Status = %q, want 'healthy'", check.Status)
	}
}

func TestPortCheckResultWithError(t *testing.T) {
	result := PortCheckResult{
		Name:        "Database",
		Host:        "db.example.com",
		Port:        5432,
		Description: "PostgreSQL database",
		Status:      "unhealthy",
		Error:       "connection refused",
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var decoded PortCheckResult
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if decoded.Error != "connection refused" {
		t.Errorf("Error = %q, want 'connection refused'", decoded.Error)
	}

	if decoded.Status != "unhealthy" {
		t.Errorf("Status = %q, want 'unhealthy'", decoded.Status)
	}
}

func TestPortCheckResultOmitEmptyError(t *testing.T) {
	result := PortCheckResult{
		Name:   "Service",
		Host:   "localhost",
		Port:   8080,
		Status: "healthy",
		Error:  "", // Empty error should be omitted
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	jsonStr := string(jsonData)

	// Error field should not be present in JSON when empty
	if containsErrorField(jsonStr) {
		t.Error("JSON should not contain 'error' field when it's empty")
	}
}

func containsErrorField(jsonStr string) bool {
	// Check if "error" key exists in JSON string
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return false
	}
	_, exists := data["error"]
	return exists
}

func TestPortCheckJSONTags(t *testing.T) {
	check := PortCheck{
		Host:        "api.example.com",
		Port:        443,
		Name:        "API",
		Description: "REST API",
	}

	jsonData, err := json.Marshal(check)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify JSON field names match the struct tags
	expectedFields := []string{
		`"host":"api.example.com"`,
		`"port":443`,
		`"name":"API"`,
		`"description":"REST API"`,
	}

	for _, expected := range expectedFields {
		if !contains(jsonStr, expected) {
			t.Errorf("JSON does not contain %q. Got: %s", expected, jsonStr)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestHealthStatusUnhealthyWithMultipleFailures(t *testing.T) {
	status := HealthStatus{
		Status:  "unhealthy",
		Message: "Failed ports: [Service1, Service2]",
		Checks: []PortCheckResult{
			{
				Name:   "Service1",
				Host:   "host1",
				Port:   8080,
				Status: "unhealthy",
				Error:  "connection timeout",
			},
			{
				Name:   "Service2",
				Host:   "host2",
				Port:   8081,
				Status: "unhealthy",
				Error:  "connection refused",
			},
			{
				Name:   "Service3",
				Host:   "host3",
				Port:   8082,
				Status: "healthy",
			},
		},
		Time:    time.Now().Format(time.RFC3339),
		Version: "1.0.0",
	}

	jsonData, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var decoded HealthStatus
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(decoded.Checks) != 3 {
		t.Fatalf("len(Checks) = %d, want 3", len(decoded.Checks))
	}

	unhealthyCount := 0
	for _, check := range decoded.Checks {
		if check.Status == "unhealthy" {
			unhealthyCount++
			if check.Error == "" {
				t.Errorf("Unhealthy check %q should have error message", check.Name)
			}
		}
	}

	if unhealthyCount != 2 {
		t.Errorf("Unhealthy count = %d, want 2", unhealthyCount)
	}
}

func TestEmptyConfig(t *testing.T) {
	cfg := Config{}

	if cfg.Server.Port != "" {
		t.Errorf("Empty config should have empty port, got %q", cfg.Server.Port)
	}

	if cfg.Server.Timeout != 0 {
		t.Errorf("Empty config should have zero timeout, got %v", cfg.Server.Timeout)
	}

	if len(cfg.Checks) != 0 {
		t.Errorf("Empty config should have no checks, got %d", len(cfg.Checks))
	}
}
