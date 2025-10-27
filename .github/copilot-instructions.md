# PortGuard AI Coding Agent Instructions

## Project Overview

PortGuard is a **single-binary HTTP health check service** written in Go that monitors TCP port availability. The architecture is intentionally flat—all code lives in the root directory as individual `.go` files, not packages. The binary is stateless and reads all configuration from a YAML file at startup.

**Core flow**: HTTP request → handler → checker → TCP dial → JSON response

## Architecture & File Structure

- **`main.go`**: Entry point with flag parsing, config loading, HTTP server setup
- **`config.go`**: YAML config loading with defaults (`/etc/portguard/config.yaml`)
- **`types.go`**: All struct definitions (Config, PortCheck, HealthStatus, etc.)
- **`checker.go`**: TCP port checking logic (`net.DialTimeout`)
- **`handlers.go`**: Three HTTP handlers: `/health` (JSON), `/live` (text), `/` (HTML)

**Key pattern**: Handlers receive `*Config` via closure from `main.go`, avoiding global state.

## Critical Development Workflows

### Building & Testing
```bash
make build              # Single binary with version from git tags
make build-all          # Cross-compile for 5 platforms (linux/darwin/windows)
make test               # Tests with race detector
make lint               # Requires golangci-lint installed locally
```

**Version injection**: `appVersion` variable in `main.go` is overridden at build time via `-ldflags "-X main.appVersion=$(VERSION)"`. The Makefile auto-detects version from `git describe --tags`.

### Running Locally
```bash
make run                # Uses config.yaml.example
./portguard --config /path/to/config.yaml
./portguard --version
```

## Project-Specific Conventions

### Error Handling Pattern
- Config loading: Fatal on error (service cannot run without config)
- Port checks: Non-fatal, capture error in `PortCheckResult.Error` field
- HTTP handlers: Always return 200 for `/live`, 200/503 for `/health` based on check results

### HTTP Response Contract
- **`/health`**: Returns JSON with `status: "healthy"|"unhealthy"`, array of per-port results, and fails with HTTP 503 if ANY port check fails
- **`/live`**: Always returns HTTP 200 "OK" (used for Kubernetes liveness—app is alive even if downstream ports are down)
- **`/`**: HTML info page, returns 404 for any path other than exact `/`

### Configuration
- Timeout defaults to 2 seconds if not specified
- Port defaults to "8888" (stored as string because it becomes `":8888"` for `http.ListenAndServe`)
- No validation of check uniqueness—duplicate host:port entries are allowed

## Code Style & Patterns

### Go Idioms in Use
- Closures for handler injection: `healthHandler(cfg)` returns `http.HandlerFunc`
- Error wrapping: `fmt.Errorf("msg: %w", err)` for context preservation
- Minimal dependencies: Only `gopkg.in/yaml.v3` beyond standard library
- No goroutines for checks—sequential execution is intentional (simple, predictable)

### Constants
- `defaultConfigPath = "/etc/portguard/config.yaml"`
- `defaultListenPort = "8888"`
- `headerContentType = "Content-Type"` (single use, but extracted for consistency)

## Integration Points

### Deployment Targets
- **Systemd**: `portguard.service` expects binary at `/opt/portguard/portguard`
- **Docker**: Multi-stage build, Alpine-based, runs as root (standard for port binding)
- **Kubernetes**: Examples in `docs/EXAMPLES.md` use `/live` for liveness, `/health` for readiness

### External Dependencies
- **At runtime**: Only requires network access to monitored hosts
- **At build time**: Go 1.23+, optional golangci-lint for linting

## Common Pitfalls & Gotchas

1. **Config file location**: Default is `/etc/portguard/config.yaml`. Local testing requires `--config` flag or copying config to `/etc/portguard/`.

2. **Port type**: Stored as `int` in structs but server port is `string` (historical reason—Go's http package expects `":8888"` format).

3. **No graceful shutdown**: HTTP server blocks forever; systemd handles SIGTERM. Future enhancement opportunity.

4. **Tests**: Project now includes standard Go `*_test.go` files. Use `make test` (race + coverage) and `make test-coverage` for HTML report. Manual script-based testing is no longer required.

## Documentation Structure

- **README.md**: Quick start, links to detailed docs
- **docs/QUICKSTART.md**: Installation for all platforms (binary, Docker, systemd, Kubernetes)
- **docs/EXAMPLES.md**: Config examples (web apps, databases, microservices, k8s probes)
- **docs/FAQ.md**: Load balancer integration (HAProxy, Nginx), troubleshooting

## Release Process

Managed in `release/` directory:
- `build.sh`: Cross-compiles binaries
- `release.sh`: Creates GitHub release with tar.gz/zip archives
- Versioning: Git tags drive both build version and release artifacts
