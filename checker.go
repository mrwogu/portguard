// Package main implements the PortGuard HTTP health check service
package main

import (
	"fmt"
	"net"
	"time"
)

func checkPort(host string, port int, timeout time.Duration) error {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()
	return nil
}

func performHealthCheck(cfg *Config) HealthStatus {
	results := make([]PortCheckResult, 0, len(cfg.Checks))
	allHealthy := true
	failedPorts := []string{}

	for _, portCheck := range cfg.Checks {
		result := PortCheckResult{
			Name:        portCheck.Name,
			Host:        portCheck.Host,
			Port:        portCheck.Port,
			Description: portCheck.Description,
		}

		err := checkPort(portCheck.Host, portCheck.Port, cfg.Server.Timeout)
		if err != nil {
			result.Status = "unhealthy"
			result.Error = err.Error()
			allHealthy = false
			failedPorts = append(failedPorts, fmt.Sprintf("%s (%s:%d)", portCheck.Name, portCheck.Host, portCheck.Port))
		} else {
			result.Status = "healthy"
		}

		results = append(results, result)
	}

	status := HealthStatus{
		Checks:  results,
		Time:    time.Now().Format(time.RFC3339),
		Version: appVersion,
	}

	if allHealthy {
		status.Status = "healthy"
		status.Message = "All ports are listening and accessible"
	} else {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("Failed ports: %v", failedPorts)
	}

	return status
}
