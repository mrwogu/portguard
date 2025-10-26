package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const headerContentType = "Content-Type"

func healthHandler(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		status := performHealthCheck(cfg)

		w.Header().Set(headerContentType, "application/json")

		if status.Status == "healthy" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		_ = json.NewEncoder(w).Encode(status)
	}
}

func liveHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(headerContentType, "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, "OK")
}

func rootHandler(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set(headerContentType, "text/html")
		_, _ = fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>PortGuard - Health Check Service</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; }
        .version { color: #666; font-size: 14px; }
        ul { line-height: 1.8; }
        a { color: #0066cc; text-decoration: none; }
        a:hover { text-decoration: underline; }
        code { background: #f0f0f0; padding: 2px 6px; border-radius: 3px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🛡️ PortGuard - Health Check Service</h1>
        <p class="version">Version: %s</p>
        <p>A lightweight HTTP service for monitoring port availability.</p>
        <h2>Available Endpoints:</h2>
        <ul>
            <li><a href="/health"><code>/health</code></a> - Detailed health status with all port checks (JSON)</li>
            <li><a href="/live"><code>/live</code></a> - Simple liveness check (returns OK)</li>
        </ul>
        <h2>Configuration:</h2>
        <p>Monitoring <strong>%d ports</strong></p>
    </div>
</body>
</html>
`, appVersion, len(cfg.Checks))
	}
}
