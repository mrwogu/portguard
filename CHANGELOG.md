# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- HTTP Basic Authentication support for securing endpoints
  - Optional authentication with `enabled`, `username`, and `password` configuration
  - All endpoints (`/health`, `/live`, `/`) protected when authentication is enabled
  - Constant-time comparison to prevent timing attacks
  - Returns 401 Unauthorized with `WWW-Authenticate` header for invalid credentials
  - Authentication is disabled by default for backward compatibility
- Comprehensive test suite for authentication middleware (7 test cases)
- Documentation for authentication in FAQ.md and EXAMPLES.md
- Per-check timeout configuration: Each check entry can now specify its own timeout value
- When a check-specific timeout is set, it overrides the global server timeout for that check
- Support for different timeout values per service (e.g., 500ms for fast local services, 10s for remote APIs)
- Added `timeout` field to `PortCheck` struct with `omitempty` YAML tag
- Updated configuration examples to demonstrate per-check timeout usage

### Changed
- Enhanced `ServerConfig` struct with `AuthConfig` for authentication settings
- Updated handlers to use `basicAuthMiddleware` wrapper
- Startup logging now indicates authentication status
- Health check logic now uses check-specific timeout when available, falling back to server timeout
- Enhanced example configurations with timeout demonstrations

### Security
- Uses `crypto/subtle.ConstantTimeCompare` to prevent timing attacks on password comparison
- Supports HTTP Basic Auth over HTTPS when deployed behind reverse proxy

### Documentation
- Added security section in FAQ.md with authentication best practices
- Added secured configuration example in EXAMPLES.md
- Updated README.md to list authentication as a feature
- Updated config.yaml.example with authentication configuration options
- Updated README.md with per-check timeout example
- Added new "Per-Check Timeout Configuration" section in EXAMPLES.md
- Updated config.yaml.example with comments about per-check timeouts
- Created config-timeout-demo.yaml demonstrating various timeout scenarios

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
