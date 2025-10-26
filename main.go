package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var appVersion = "1.0.0"

func main() {
	configPath := flag.String("config", defaultConfigPath, "Path to configuration file")
	showVersion := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("PortGuard version %s\n", appVersion)
		os.Exit(0)
	}

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	if len(cfg.Checks) == 0 {
		log.Fatal("No port checks configured. Please add checks to the configuration file.")
	}

	// Wrap handlers with authentication middleware
	http.HandleFunc("/health", basicAuthMiddleware(cfg, healthHandler(cfg)))
	http.HandleFunc("/live", basicAuthMiddleware(cfg, liveHandler))
	http.HandleFunc("/", basicAuthMiddleware(cfg, rootHandler(cfg)))

	listenAddr := ":" + cfg.Server.Port

	log.Printf("PortGuard v%s starting...", appVersion)
	log.Printf("Configuration loaded from: %s", *configPath)
	log.Printf("Monitoring %d ports with %s timeout", len(cfg.Checks), cfg.Server.Timeout)
	if cfg.Server.Auth.Enabled && cfg.Server.Auth.Username != "" && cfg.Server.Auth.Password != "" {
		log.Printf("HTTP Basic Authentication: ENABLED (username: %s)", cfg.Server.Auth.Username)
	} else {
		log.Printf("HTTP Basic Authentication: DISABLED")
	}
	log.Printf("HTTP server listening on %s", listenAddr)
	log.Printf("Endpoints:")
	log.Printf("  - http://localhost%s/health (detailed JSON status)", listenAddr)
	log.Printf("  - http://localhost%s/live (simple OK response)", listenAddr)

	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}
