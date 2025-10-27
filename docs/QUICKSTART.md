# PortGuard Quick Start Guide

Get up and running with PortGuard in 5 minutes!

## Option 1: Download Pre-built Binary (Recommended)

### Step 1: Download

Go to [Releases](https://github.com/mrwogu/portguard/releases/latest) and download the appropriate binary for your system:

- **Linux (64-bit)**: `portguard-linux-amd64.tar.gz`
- **macOS (Intel)**: `portguard-darwin-amd64.tar.gz`
- **macOS (Apple Silicon)**: `portguard-darwin-arm64.tar.gz`
- **Windows**: `portguard-windows-amd64.zip`

### Step 2: Extract and Install

**Linux/macOS:**
```bash
tar -xzf portguard-*.tar.gz
chmod +x portguard-*
sudo mv portguard-* /usr/local/bin/portguard
```

**Windows:**
Extract the ZIP file and add the directory to your PATH.

### Step 3: Create Configuration

```bash
# Create config directory
sudo mkdir -p /etc/portguard

# Download example config
curl -o /tmp/config.yaml.example https://raw.githubusercontent.com/mrwogu/portguard/main/config.yaml.example

# Copy and edit
sudo cp /tmp/config.yaml.example /etc/portguard/config.yaml
sudo nano /etc/portguard/config.yaml
```

Edit the configuration to monitor your ports:

```yaml
server:
  port: "8888"
  timeout: 2s

checks:
  - host: "your-server.com"
    port: 80
    name: "Web Server"
    description: "HTTP"
```

### Step 4: Run

```bash
# Run directly
portguard

# Or specify config location
portguard --config /etc/portguard/config.yaml
```

### Step 5: Test

Open your browser to `http://localhost:8888` or use curl:

```bash
curl http://localhost:8888/health
```

## Option 2: Using Docker

### Quick Test

```bash
# Create a config file
cat > config.yaml << 'EOF'
server:
  port: "8888"
  timeout: 2s
checks:
  - host: "google.com"
    port: 443
    name: "Google"
    description: "Test connection"
EOF

# Run with Docker
docker run -d \
  --name portguard \
  -p 8888:8888 \
  -v $(pwd)/config.yaml:/etc/portguard/config.yaml:ro \
  ghcr.io/mrwogu/portguard:latest

# Check logs
docker logs portguard

# Test
curl http://localhost:8888/health
```

### Using Docker Compose

```bash
# Create docker-compose.yml
cat > docker-compose.yml << 'EOF'
version: '3.8'
services:
  portguard:
    image: ghcr.io/mrwogu/portguard:latest
    ports:
      - "8888:8888"
    volumes:
      - ./config.yaml:/etc/portguard/config.yaml:ro
    restart: unless-stopped
EOF

# Start
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

## Option 3: Build from Source

### Prerequisites

- Go 1.23 or later
- Git

### Steps

```bash
# Clone repository
git clone https://github.com/mrwogu/portguard.git
cd portguard

# Build
go build -o portguard main.go

# Or use Make
make build

# Create config
cp config.yaml.example config.yaml
nano config.yaml

# Run
./portguard --config config.yaml
```

## Option 4: Install as System Service (Linux)

### Using Make (Recommended)

The `Makefile` includes an `install` target that builds the binary, installs it to `/opt/portguard`, installs a default config to `/etc/portguard/config.yaml` if one does not already exist, and copies the systemd service file.

```bash
# Build and install (requires root for system locations)
sudo make install

# Edit configuration (only if you need changes beyond the example)
sudo nano /etc/portguard/config.yaml

# Start service
sudo systemctl start portguard

# Check status
sudo systemctl status portguard

# Enable auto-start on boot
sudo systemctl enable portguard

# View logs
sudo journalctl -u portguard -f

# Uninstall (optional)
sudo make uninstall
```

### Manual Installation (Alternative)

```bash
# Copy binary
sudo mkdir -p /opt/portguard
sudo cp portguard /opt/portguard/
sudo chmod +x /opt/portguard/portguard

# Copy config
sudo mkdir -p /etc/portguard
sudo cp config.yaml.example /etc/portguard/config.yaml

# Copy service file
sudo cp portguard.service /etc/systemd/system/

# Reload and start
sudo systemctl daemon-reload
sudo systemctl start portguard
sudo systemctl enable portguard
```

## Testing Your Installation

### 1. Check if Service is Running

```bash
curl http://localhost:8888/live
# Should return: OK
```

### 2. Check Health Status

```bash
curl http://localhost:8888/health | jq
```

Example healthy response:
```json
{
  "status": "healthy",
  "message": "All ports are listening and accessible",
  "checks": [
    {
      "name": "Web Server",
      "host": "example.com",
      "port": 80,
      "description": "HTTP",
      "status": "healthy"
    }
  ],
  "timestamp": "2025-10-26T10:30:00Z",
  "version": "1.0.0"
}
```

### 3. Run Built-in Tests

PortGuard includes standard Go tests. Run them with:

```bash
make test       # Runs tests with race detector and coverage
make test-coverage  # (Optional) Generates HTML coverage report
```

## Common Issues

### Port Already in Use

If port 8888 is already in use, change it in `config.yaml`:

```yaml
server:
  port: "9999"  # Use a different port
  timeout: 2s
```

### Permission Denied

If you get permission errors:

```bash
# Give executable permissions
chmod +x portguard

# Or run with sudo if binding to ports < 1024
sudo ./portguard
```

### Connection Refused

If checks are failing:

1. Verify the host is reachable: `ping your-host`
2. Check if the port is open: `telnet your-host port`
3. Verify firewall rules
4. Increase timeout in config if network is slow

## Next Steps

- 📖 Read the full [README](../README.md)
- 🔧 See [EXAMPLES.md](EXAMPLES.md) for configuration examples
- 🤝 Check [CONTRIBUTING.md](../.github/CONTRIBUTING.md) to contribute
- 🔒 Review [SECURITY.md](../.github/SECURITY.md) for security best practices

## Getting Help

- 📝 [File an Issue](https://github.com/mrwogu/portguard/issues/new)
- 📚 [Read the Documentation](https://github.com/mrwogu/portguard)

---

**That's it!** You now have PortGuard monitoring your services. 🎉
