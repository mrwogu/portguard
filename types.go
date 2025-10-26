package main

import "time"

// Config represents the main configuration structure for PortGuard.
// It contains server settings and a list of ports to check.
type Config struct {
	Server ServerConfig `yaml:"server"`
	Checks []PortCheck  `yaml:"checks"`
}

// ServerConfig holds the HTTP server configuration.
// Port specifies which port the HTTP server listens on.
// Timeout sets the maximum duration for port check operations.
// Auth contains optional HTTP Basic Authentication settings.
type ServerConfig struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
	Auth    AuthConfig    `yaml:"auth,omitempty"`
}

// AuthConfig holds HTTP Basic Authentication configuration.
// When both Username and Password are empty, authentication is disabled.
type AuthConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// PortCheck defines a single port to monitor.
// It includes the target host, port number, and descriptive information.
// An optional Timeout can be specified per check, otherwise the server timeout is used.
type PortCheck struct {
	Host        string        `yaml:"host" json:"host"`
	Port        int           `yaml:"port" json:"port"`
	Name        string        `yaml:"name" json:"name"`
	Description string        `yaml:"description" json:"description"`
	Timeout     time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

// HealthStatus represents the overall health check response.
// It contains the aggregated status and results from all port checks.
type HealthStatus struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Checks  []PortCheckResult `json:"checks"`
	Time    string            `json:"timestamp"`
	Version string            `json:"version"`
}

// PortCheckResult holds the result of checking a single port.
// It includes the check details and whether the port is reachable.
type PortCheckResult struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Error       string `json:"error,omitempty"`
}
