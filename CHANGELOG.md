# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-26

### Added
- Initial release of PortGuard - lightweight HTTP service for monitoring port availability
- Core health checking functionality for TCP ports
- YAML-based configuration system with support for multiple port checks
- HTTP API endpoints:
  - `/health` - Detailed JSON status of all monitored ports
  - `/live` - Simple liveness probe endpoint
  - `/` - HTML info page with service information
- Response codes: 200 OK (all healthy) / 503 Service Unavailable (one or more unhealthy)
- Configurable timeouts for port checks
- Command-line flags: `--config` for custom config path, `--version` for version info
- Per-check configuration: host, port, name, and description
- Detailed JSON responses with per-port status information
- Version information included in all responses

### Infrastructure
- Systemd service file (`portguard.service`) for Linux service management
- Docker support with Dockerfile and docker-compose.yml
- Multi-platform build scripts for Linux, macOS, Windows (amd64, arm64)
- Release automation scripts (`build.sh`, `release.sh`)
- Installation script for automated setup

### Configuration
- Example configuration file (`config.yaml.example`)
- Support for mail server monitoring (SMTP, IMAP, POP3, Submission)
- Flexible configuration structure for any TCP service monitoring
- Default config path: `/etc/portguard/config.yaml`

### Features
- Single binary with no external dependencies
- Minimal resource footprint
- Production-ready with proper error handling
- Structured logging for operations and debugging
- Cross-platform compatibility (Linux, macOS, Windows)
- Container-friendly design for Docker and Kubernetes deployments
