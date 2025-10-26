package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestMainVersionFlag tests the --version flag
func TestMainVersionFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "portguard-test")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	// Run with --version flag
	cmd = exec.Command(binaryPath, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run binary with --version: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "PortGuard version") {
		t.Errorf("Version output does not contain expected text. Got: %s", outputStr)
	}
}

// TestMainMissingConfig tests behavior when config file is missing
func TestMainMissingConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "portguard-test")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	// Run with non-existent config
	cmd = exec.Command(binaryPath, "--config", "/nonexistent/config.yaml")
	output, err := cmd.CombinedOutput()

	// Should exit with error
	if err == nil {
		t.Error("Expected error for missing config, got none")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Error loading configuration") {
		t.Errorf("Error message does not mention config loading. Got: %s", outputStr)
	}
}

// TestMainNoChecksConfigured tests behavior when no checks are configured
func TestMainNoChecksConfigured(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "portguard-test")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	// Create empty config
	configPath := filepath.Join(tmpDir, "empty-config.yaml")
	emptyConfig := `server:
  port: "9999"
checks: []
`
	if err := os.WriteFile(configPath, []byte(emptyConfig), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Run with empty checks
	cmd = exec.Command(binaryPath, "--config", configPath)
	output, err := cmd.CombinedOutput()

	// Should exit with error
	if err == nil {
		t.Error("Expected error for no checks configured, got none")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "No port checks configured") {
		t.Errorf("Error message does not mention no checks. Got: %s", outputStr)
	}
}

// TestMainServerStart tests that the server starts successfully
func TestMainServerStart(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "portguard-test")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	// Create valid config with unique port
	configPath := filepath.Join(tmpDir, "test-config.yaml")
	testConfig := `server:
  port: "18888"
  timeout: 2s
checks:
  - host: "localhost"
    port: 22
    name: "SSH"
    description: "SSH service"
`
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Start server in background
	cmd = exec.Command(binaryPath, "--config", configPath)

	// Capture stdout/stderr
	outputBuilder := &strings.Builder{}
	cmd.Stdout = outputBuilder
	cmd.Stderr = outputBuilder

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Ensure cleanup
	defer func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}()

	// Give the server time to start
	time.Sleep(500 * time.Millisecond)

	// Check if process is still running
	if cmd.Process == nil {
		t.Fatal("Server process is nil")
	}

	// Try to signal the process to verify it's running
	if err := cmd.Process.Signal(os.Signal(os.Interrupt)); err != nil {
		// On some systems, the process might have already exited
		output := outputBuilder.String()
		if !strings.Contains(output, "HTTP server listening") {
			t.Errorf("Server did not start properly. Output: %s", output)
		}
	}
}

// TestMainDefaultConfigPath tests default config path behavior
func TestMainDefaultConfigPath(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "portguard-test")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	// Run without config flag (will try default path)
	cmd = exec.Command(binaryPath)
	output, err := cmd.CombinedOutput()

	// Should fail because default config doesn't exist
	if err == nil {
		t.Error("Expected error for missing default config, got none")
	}

	outputStr := string(output)
	// Should mention the default config path
	if !strings.Contains(outputStr, defaultConfigPath) && !strings.Contains(outputStr, "Error loading configuration") {
		t.Errorf("Error should mention config issue. Got: %s", outputStr)
	}
}

// TestAppVersion tests that appVersion variable is set
func TestAppVersion(t *testing.T) {
	if appVersion == "" {
		t.Error("appVersion should not be empty")
	}
}

// TestMainHelpFlag tests the --help flag
func TestMainHelpFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "portguard-test")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	// Run with --help flag
	cmd = exec.Command(binaryPath, "--help")
	output, _ := cmd.CombinedOutput()

	// Help should exit with code 0 or non-zero (depends on flag package)
	outputStr := string(output)

	// Should show usage information
	expectedStrings := []string{
		"config",
		"version",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(strings.ToLower(outputStr), strings.ToLower(expected)) {
			t.Errorf("Help output should contain %q. Got: %s", expected, outputStr)
		}
	}
}

// TestMainInvalidFlag tests behavior with invalid flag
func TestMainInvalidFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "portguard-test")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	// Run with invalid flag
	cmd = exec.Command(binaryPath, "--invalid-flag")
	output, _ := cmd.CombinedOutput()

	// The flag package will handle this and show usage
	outputStr := string(output)
	if !strings.Contains(outputStr, "flag provided but not defined") {
		t.Logf("Output for invalid flag: %s", outputStr)
	}
}
