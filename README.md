# PortGuard 🛡️

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.23-blue)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![CI](https://github.com/mrwogu/portguard/actions/workflows/ci.yml/badge.svg)](https://github.com/mrwogu/portguard/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/mrwogu/portguard/graph/badge.svg?token=8KKR2UY0TS)](https://codecov.io/gh/mrwogu/portguard)
[![Release](https://img.shields.io/github/v/release/mrwogu/portguard?sort=semver)](https://github.com/mrwogu/portguard/releases/latest)

A lightweight, configurable HTTP service for monitoring port availability and health checking. Perfect for load balancers, Kubernetes health probes, and service monitoring.

![PortGuard](assets/portguard.jpg)

📚 **[Quick Start](docs/QUICKSTART.md)** | 📖 **[Examples](docs/EXAMPLES.md)** | ❓ **[FAQ](docs/FAQ.md)**

## Features

✨ **Simple & Lightweight** - Single binary, minimal dependencies  
⚙️ **Configurable** - YAML-based configuration  
🔍 **Detailed Health Checks** - JSON responses with per-port status  
🔒 **Secure** - Optional HTTP Basic Authentication  
🚀 **Production Ready** - Systemd, Docker, Kubernetes support  

## Quick Start

```bash
# Download latest release
wget https://github.com/mrwogu/portguard/releases/latest/download/portguard-linux-amd64.tar.gz
tar -xzf portguard-*.tar.gz

# Create config
mkdir -p /etc/portguard
cp config.yaml.example /etc/portguard/config.yaml
nano /etc/portguard/config.yaml

# Run
./portguard
```

**📖 See [Quick Start Guide](docs/QUICKSTART.md) for detailed installation instructions including Docker, systemd, and other platforms.**

## Configuration

```yaml
server:
  port: "8888"
  timeout: 2s  # Default timeout for all checks
  
  # Optional: HTTP Basic Authentication
  auth:
    enabled: false  # Set to true to enable
    username: "admin"
    password: "secure-password"

checks:
  - host: "mail.example.com"
    port: 25
    name: "SMTP"
    description: "Mail Transfer Protocol"
    # Uses default timeout (2s)
  
  - host: "remote-api.example.com"
    port: 443
    name: "Remote API"
    description: "Remote API endpoint"
    timeout: 10s  # Custom timeout for this check
```

**📖 See [Examples](docs/EXAMPLES.md) for more configuration examples (web apps, databases, Kubernetes, microservices).**

## API Endpoints

- **`/health`** - Detailed JSON status (200 OK = healthy, 503 = unhealthy)
- **`/live`** - Simple liveness probe (always returns 200 OK)
- **`/`** - HTML info page

**📖 See [FAQ](docs/FAQ.md) for integration examples with HAProxy, Nginx, and Kubernetes.**

## Development

```bash
make build          # Build binary
make test           # Run tests with coverage
make test-coverage  # Generate HTML coverage report
make lint           # Run linter
make build-all      # Build for all platforms
```

See `Makefile` for all available commands.

**👨‍💻 For contributors:** See [Release Process](docs/RELEASING.md) for creating new releases.

## Support & Documentation

- 📚 [Quick Start Guide](docs/QUICKSTART.md) - Installation & setup
- 📖 [Configuration Examples](docs/EXAMPLES.md) - Various use cases
- ❓ [FAQ](docs/FAQ.md) - Common questions & troubleshooting
- 🚀 [Release Process](docs/RELEASING.md) - For maintainers

---

Made with ❤️ for reliable service monitoring
