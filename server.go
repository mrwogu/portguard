package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

// serverStarter is a function type that starts the HTTP server
type serverStarter func(addr string, handler http.Handler) error

// run is the main application logic, separated from main() for testability
func run(args []string, exit func(int), startServer serverStarter) error {
	fs := flag.NewFlagSet("portguard", flag.ContinueOnError)
	configPath := fs.String("config", defaultConfigPath, "Path to configuration file")
	showVersion := fs.Bool("version", false, "Show version and exit")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *showVersion {
		fmt.Printf("PortGuard version %s\n", appVersion)
		exit(0)
		return nil
	}

	cfg, err := loadConfig(*configPath)
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	if len(cfg.Checks) == 0 {
		return fmt.Errorf("no port checks configured. Please add checks to the configuration file")
	}

	return setupAndStartServer(cfg, *configPath, startServer)
}

// setupAndStartServer configures HTTP handlers and starts the server
func setupAndStartServer(cfg *Config, configPath string, startServer serverStarter) error {
	// Create a new ServeMux for this server instance
	mux := http.NewServeMux()
	
	// Wrap handlers with authentication middleware
	mux.HandleFunc("/health", basicAuthMiddleware(cfg, healthHandler(cfg)))
	mux.HandleFunc("/live", basicAuthMiddleware(cfg, liveHandler))
	mux.HandleFunc("/", basicAuthMiddleware(cfg, rootHandler(cfg)))

	listenAddr := ":" + cfg.Server.Port

	log.Printf("PortGuard v%s starting...", appVersion)
	log.Printf("Configuration loaded from: %s", configPath)
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

	if err := startServer(listenAddr, mux); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
