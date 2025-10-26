package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestCheckPort(t *testing.T) {
	// Start a test TCP server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer func() { _ = listener.Close() }()

	// Get the actual port assigned
	_, portStr, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to parse listener address: %v", err)
	}
	var testPort int
	_, _ = fmt.Sscanf(portStr, "%d", &testPort)

	// Accept connections in background
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
		name    string
		host    string
		port    int
		timeout time.Duration
		wantErr bool
	}{
		{
			name:    "successful connection to listening port",
			host:    "127.0.0.1",
			port:    testPort,
			timeout: 2 * time.Second,
			wantErr: false,
		},
		{
			name:    "connection to closed port",
			host:    "127.0.0.1",
			port:    testPort + 1, // Adjacent port should be closed
			timeout: 1 * time.Second,
			wantErr: true,
		},
		{
			name:    "connection timeout",
			host:    "192.0.2.1", // TEST-NET-1, should timeout
			port:    80,
			timeout: 100 * time.Millisecond,
			wantErr: true,
		},
		{
			name:    "invalid host",
			host:    "invalid.host.that.does.not.exist.local",
			port:    80,
			timeout: 1 * time.Second,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkPort(tt.host, tt.port, tt.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkPort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckPortTimeout(t *testing.T) {
	// Use a non-routable IP to ensure timeout
	host := "192.0.2.1" // TEST-NET-1
	port := 80
	timeout := 500 * time.Millisecond

	start := time.Now()
	err := checkPort(host, port, timeout)
	elapsed := time.Since(start)

	if err == nil {
		t.Error("Expected error for non-routable address, got none")
	}

	// Timeout should occur within reasonable bounds
	if elapsed > timeout+1*time.Second {
		t.Errorf("Timeout took too long: %v (expected ~%v)", elapsed, timeout)
	}
}

func TestPerformHealthCheck(t *testing.T) {
	// Start multiple test servers
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

	// Get ports
	var port1, port2 int
	_, portStr1, _ := net.SplitHostPort(listener1.Addr().String())
	_, portStr2, _ := net.SplitHostPort(listener2.Addr().String())
	_, _ = fmt.Sscanf(portStr1, "%d", &port1)
	_, _ = fmt.Sscanf(portStr2, "%d", &port2)

	// Accept connections
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

	tests := []struct {
		name        string
		config      *Config
		wantStatus  string
		wantHealthy int
	}{
		{
			name: "all ports healthy",
			config: &Config{
				Server: ServerConfig{
					Port:    "8888",
					Timeout: 2 * time.Second,
				},
				Checks: []PortCheck{
					{
						Host:        "127.0.0.1",
						Port:        port1,
						Name:        "Service 1",
						Description: "Test service 1",
					},
					{
						Host:        "127.0.0.1",
						Port:        port2,
						Name:        "Service 2",
						Description: "Test service 2",
					},
				},
			},
			wantStatus:  "healthy",
			wantHealthy: 2,
		},
		{
			name: "one port unhealthy",
			config: &Config{
				Server: ServerConfig{
					Port:    "8888",
					Timeout: 1 * time.Second,
				},
				Checks: []PortCheck{
					{
						Host:        "127.0.0.1",
						Port:        port1,
						Name:        "Service 1",
						Description: "Test service 1",
					},
					{
						Host:        "127.0.0.1",
						Port:        99999, // Invalid port
						Name:        "Service 2",
						Description: "Test service 2",
					},
				},
			},
			wantStatus:  "unhealthy",
			wantHealthy: 1,
		},
		{
			name: "all ports unhealthy",
			config: &Config{
				Server: ServerConfig{
					Port:    "8888",
					Timeout: 1 * time.Second,
				},
				Checks: []PortCheck{
					{
						Host:        "192.0.2.1",
						Port:        80,
						Name:        "Unreachable 1",
						Description: "Should fail",
					},
					{
						Host:        "192.0.2.2",
						Port:        80,
						Name:        "Unreachable 2",
						Description: "Should fail",
					},
				},
			},
			wantStatus:  "unhealthy",
			wantHealthy: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := performHealthCheck(tt.config)

			if status.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q", status.Status, tt.wantStatus)
			}

			if len(status.Checks) != len(tt.config.Checks) {
				t.Errorf("Number of check results = %d, want %d", len(status.Checks), len(tt.config.Checks))
			}

			healthyCount := 0
			for _, check := range status.Checks {
				if check.Status == "healthy" {
					healthyCount++
				}
			}

			if healthyCount != tt.wantHealthy {
				t.Errorf("Healthy count = %d, want %d", healthyCount, tt.wantHealthy)
			}

			// Verify timestamp is set
			if status.Time == "" {
				t.Error("Timestamp is empty")
			}

			// Verify version is set
			if status.Version == "" {
				t.Error("Version is empty")
			}

			// Verify message is set
			if status.Message == "" {
				t.Error("Message is empty")
			}
		})
	}
}

func TestPerformHealthCheckResultDetails(t *testing.T) {
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

	cfg := &Config{
		Server: ServerConfig{
			Port:    "8888",
			Timeout: 2 * time.Second,
		},
		Checks: []PortCheck{
			{
				Host:        "127.0.0.1",
				Port:        testPort,
				Name:        "Test Service",
				Description: "A test service",
			},
		},
	}

	status := performHealthCheck(cfg)

	if len(status.Checks) != 1 {
		t.Fatalf("Expected 1 check result, got %d", len(status.Checks))
	}

	result := status.Checks[0]

	if result.Name != "Test Service" {
		t.Errorf("Result.Name = %q, want 'Test Service'", result.Name)
	}

	if result.Host != "127.0.0.1" {
		t.Errorf("Result.Host = %q, want '127.0.0.1'", result.Host)
	}

	if result.Port != testPort {
		t.Errorf("Result.Port = %d, want %d", result.Port, testPort)
	}

	if result.Description != "A test service" {
		t.Errorf("Result.Description = %q, want 'A test service'", result.Description)
	}

	if result.Status != "healthy" {
		t.Errorf("Result.Status = %q, want 'healthy'", result.Status)
	}

	if result.Error != "" {
		t.Errorf("Result.Error should be empty, got %q", result.Error)
	}
}

func TestPerformHealthCheckWithError(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port:    "8888",
			Timeout: 500 * time.Millisecond,
		},
		Checks: []PortCheck{
			{
				Host:        "127.0.0.1",
				Port:        1, // Typically not listening
				Name:        "Failed Service",
				Description: "Should fail",
			},
		},
	}

	status := performHealthCheck(cfg)

	if status.Status != "unhealthy" {
		t.Errorf("Status = %q, want 'unhealthy'", status.Status)
	}

	if len(status.Checks) != 1 {
		t.Fatalf("Expected 1 check result, got %d", len(status.Checks))
	}

	result := status.Checks[0]

	if result.Status != "unhealthy" {
		t.Errorf("Result.Status = %q, want 'unhealthy'", result.Status)
	}

	if result.Error == "" {
		t.Error("Result.Error should not be empty")
	}
}

func TestPerformHealthCheckWithPerCheckTimeout(t *testing.T) {
	// Start test server
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

	cfg := &Config{
		Server: ServerConfig{
			Port:    "8888",
			Timeout: 2 * time.Second, // Default timeout
		},
		Checks: []PortCheck{
			{
				Host:        "127.0.0.1",
				Port:        testPort,
				Name:        "Service with default timeout",
				Description: "Uses server timeout",
			},
			{
				Host:        "127.0.0.1",
				Port:        testPort,
				Name:        "Service with custom timeout",
				Description: "Uses custom timeout",
				Timeout:     5 * time.Second, // Override with custom timeout
			},
			{
				Host:        "192.0.2.1", // Non-routable, will timeout
				Port:        80,
				Name:        "Slow service with short timeout",
				Description: "Should timeout quickly",
				Timeout:     100 * time.Millisecond, // Very short timeout
			},
		},
	}

	status := performHealthCheck(cfg)

	// First two should be healthy, third should timeout
	if len(status.Checks) != 3 {
		t.Fatalf("Expected 3 check results, got %d", len(status.Checks))
	}

	// Check that first service is healthy (uses default timeout)
	if status.Checks[0].Status != "healthy" {
		t.Errorf("First check should be healthy, got %q", status.Checks[0].Status)
	}

	// Check that second service is healthy (uses custom timeout)
	if status.Checks[1].Status != "healthy" {
		t.Errorf("Second check should be healthy, got %q", status.Checks[1].Status)
	}

	// Check that third service is unhealthy (custom short timeout)
	if status.Checks[2].Status != "unhealthy" {
		t.Errorf("Third check should be unhealthy, got %q", status.Checks[2].Status)
	}

	// Overall status should be unhealthy because one check failed
	if status.Status != "unhealthy" {
		t.Errorf("Overall status should be unhealthy, got %q", status.Status)
	}
}

func TestPerformHealthCheckTimeoutBehavior(t *testing.T) {
	tests := []struct {
		name           string
		serverTimeout  time.Duration
		checkTimeout   time.Duration
		expectedFaster time.Duration // Which timeout should be effective
	}{
		{
			name:           "uses server timeout when check timeout is 0",
			serverTimeout:  2 * time.Second,
			checkTimeout:   0,
			expectedFaster: 2 * time.Second,
		},
		{
			name:           "uses check timeout when specified",
			serverTimeout:  5 * time.Second,
			checkTimeout:   1 * time.Second,
			expectedFaster: 1 * time.Second,
		},
		{
			name:           "check timeout can be longer than server timeout",
			serverTimeout:  500 * time.Millisecond,
			checkTimeout:   3 * time.Second,
			expectedFaster: 3 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{
					Port:    "8888",
					Timeout: tt.serverTimeout,
				},
				Checks: []PortCheck{
					{
						Host:        "192.0.2.1", // Non-routable IP for timeout test
						Port:        80,
						Name:        "Test Service",
						Description: "Timeout test",
						Timeout:     tt.checkTimeout,
					},
				},
			}

			start := time.Now()
			status := performHealthCheck(cfg)
			elapsed := time.Since(start)

			// Check should fail due to timeout
			if status.Checks[0].Status != "unhealthy" {
				t.Errorf("Check should be unhealthy, got %q", status.Checks[0].Status)
			}

			// Verify timeout occurred within reasonable bounds
			// Allow 1 second margin for system delays
			maxExpected := tt.expectedFaster + 1*time.Second
			if elapsed > maxExpected {
				t.Errorf("Check took too long: %v (expected ~%v, max %v)", elapsed, tt.expectedFaster, maxExpected)
			}
		})
	}
}
