# Changelog

## [1.1.0] - 2025-10-26

### Added
- HTTP Basic Authentication support
  - Optional authentication with constant-time password comparison
  - Protects all endpoints (`/health`, `/live`, `/`)
  - Returns 401 with `WWW-Authenticate` header
  - Disabled by default for backward compatibility
- Per-check timeout configuration
  - Each check can override the global timeout
  - Supports different timeouts per service (e.g., 500ms for local, 10s for remote APIs)

## [1.0.0] - 2025-10-26

### Added
- Initial release - lightweight HTTP service for TCP port monitoring
- HTTP endpoints: `/health` (JSON), `/live` (liveness), `/` (HTML info)
- YAML configuration with multiple port checks
- Configurable timeouts and per-check settings
- Systemd service, Docker support, multi-platform builds
- Single binary with no external dependencies

[1.1.0]: https://github.com/mrwogu/portguard/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/mrwogu/portguard/tree/v1.0.0
