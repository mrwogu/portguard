# Frequently Asked Questions (FAQ)

## General Questions

### What is PortGuard?

PortGuard is a lightweight HTTP service that monitors TCP port availability. It provides health check endpoints that return JSON status information, making it perfect for load balancers, monitoring systems, and Kubernetes health probes.

### Why should I use PortGuard instead of [other tool]?

PortGuard is:
- **Simple**: Single binary, no dependencies
- **Fast**: Written in Go for performance
- **Flexible**: YAML configuration for any ports
- **Standard**: Returns HTTP 200/503 status codes
- **Production-ready**: Includes systemd service, Docker support

### Is PortGuard production-ready?

Yes! PortGuard is designed for production use with:
- Proper error handling
- Configurable timeouts
- Systemd service support
- Docker support
- Comprehensive logging

## Configuration

### Where should I put the config file?

Default location: `/etc/portguard/config.yaml`

You can also specify a custom location:
```bash
portguard --config /path/to/config.yaml
```

### How do I monitor multiple servers?

Add multiple checks in your config:

```yaml
checks:
  - host: "server1.example.com"
    port: 80
    name: "Server 1"
  - host: "server2.example.com"
    port: 80
    name: "Server 2"
```

### Can I monitor localhost ports?

Yes! Just use `localhost` or `127.0.0.1`:

```yaml
checks:
  - host: "localhost"
    port: 8080
    name: "Local App"
```

### What timeout should I use?

- **Local services**: 1-2s
- **Same datacenter**: 2-3s
- **Remote services**: 5-10s
- **Slow connections**: 10-30s

```yaml
server:
  timeout: 5s  # Adjust based on your needs
```

### Can I monitor UDP ports?

Not yet. PortGuard currently only supports TCP ports. UDP support is on the roadmap.

## Security

### How do I secure PortGuard endpoints?

PortGuard supports HTTP Basic Authentication. Enable it in your config:

```yaml
server:
  port: "8888"
  timeout: 2s
  auth:
    enabled: true
    username: "admin"
    password: "your-secure-password"
```

When authentication is enabled, all endpoints (`/health`, `/live`, `/`) require valid credentials.

### Is authentication required?

No, authentication is **disabled by default**. Enable it only when needed:

```yaml
auth:
  enabled: false  # Default - no authentication required
```

### How secure is HTTP Basic Auth?

HTTP Basic Auth sends credentials as base64-encoded strings. For production:

1. **Always use HTTPS** - Deploy behind a reverse proxy (nginx, Traefik, Caddy) with TLS
2. **Use strong passwords** - Generate random, long passwords
3. **Rotate credentials** - Change passwords periodically
4. **Use secret management** - Store credentials in environment variables or secret managers

Example with environment variables:

```bash
# Store password securely
export PORTGUARD_PASSWORD="$(openssl rand -base64 32)"

# Use in config (requires templating or script to inject)
```

### Can I use different authentication for different endpoints?

No, currently all endpoints share the same authentication settings. This is intentional to keep the configuration simple.

### Should I enable authentication for Kubernetes health probes?

It depends on your security requirements:

- **Without auth**: Simpler, works out-of-the-box with Kubernetes
- **With auth**: More secure, requires creating a Secret and configuring probes

Example Kubernetes probe with authentication:

```yaml
livenessProbe:
  httpGet:
    path: /live
    port: 8888
    httpHeaders:
    - name: Authorization
      value: Basic YWRtaW46cGFzc3dvcmQ=  # base64 encoded "admin:password"
```

Better approach - use a Secret:

```bash
# Create secret
kubectl create secret generic portguard-auth \
  --from-literal=username=admin \
  --from-literal=password=secure-password
```

## Deployment

### How do I install PortGuard?

See the [Quick Start Guide](QUICKSTART.md) for detailed instructions. Quick summary:

```bash
# Download binary
wget https://github.com/mrwogu/portguard/releases/latest/download/portguard-linux-amd64.tar.gz
tar -xzf portguard-*.tar.gz
sudo mv portguard-* /usr/local/bin/portguard

# Create config
sudo mkdir -p /etc/portguard
sudo nano /etc/portguard/config.yaml

# Run
portguard
```

### How do I run PortGuard in Docker?

```bash
docker run -d \
  -p 8888:8888 \
  -v /path/to/config.yaml:/etc/portguard/config.yaml:ro \
  ghcr.io/mrwogu/portguard:latest
```

### How do I install as a systemd service?

```bash
sudo ./scripts/install.sh
```

Or manually:
```bash
sudo cp portguard /opt/portguard/
sudo cp portguard.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl start portguard
sudo systemctl enable portguard
```

### Can I run multiple instances?

Yes! Just use different config files and ports:

```bash
portguard --config /etc/portguard/config1.yaml &
portguard --config /etc/portguard/config2.yaml &
```

## Usage

### What endpoints are available?

- `/health` - Detailed health status (JSON)
- `/live` - Simple liveness check (text)
- `/` - Information page (HTML)

### How do I integrate with my load balancer?

**HAProxy:**
```haproxy
option httpchk GET /health
http-check expect status 200
```

**Nginx:**
```nginx
health_check uri=/health;
```

See [README.md](README.md) for more examples.

### How do I use with Kubernetes?

```yaml
livenessProbe:
  httpGet:
    path: /live
    port: 8888
readinessProbe:
  httpGet:
    path: /health
    port: 8888
```

### What do the status codes mean?

- **200 OK**: All monitored ports are healthy
- **503 Service Unavailable**: One or more ports are unhealthy
- **404 Not Found**: Invalid endpoint

### How do I check only specific services?

Configure only the ports you want to monitor in `config.yaml`. PortGuard only checks what you configure.

## Troubleshooting

### PortGuard says a port is unhealthy, but it's working

Possible causes:

1. **Firewall**: Check firewall rules
   ```bash
   sudo iptables -L
   ```

2. **Network**: Test connectivity
   ```bash
   telnet hostname port
   ```

3. **Timeout**: Increase timeout in config
   ```yaml
   server:
     timeout: 10s
   ```

4. **DNS**: Verify hostname resolution
   ```bash
   nslookup hostname
   ```

### Port 8888 is already in use

Change the port in config:

```yaml
server:
  port: "9999"
```

### Permission denied errors

Either:
- Run with sudo: `sudo portguard`
- Use a port > 1024 in config
- Set capabilities: `sudo setcap 'cap_net_bind_service=+ep' /path/to/portguard`

### Config file not found

Specify the config file explicitly:

```bash
portguard --config /path/to/config.yaml
```

### How do I enable debug logging?

Currently, PortGuard logs to stdout. Run it in foreground to see logs:

```bash
portguard --config config.yaml
```

With systemd:
```bash
sudo journalctl -u portguard -f
```

### Checks are slow

Possible solutions:

1. Reduce timeout (if appropriate)
2. Check network latency
3. Verify target services are responsive
4. Use localhost for local services

## Performance

### How many ports can I monitor?

PortGuard can monitor hundreds of ports efficiently. Each check is independent and runs concurrently.

### What's the resource usage?

Very minimal:
- **Memory**: ~10-20 MB
- **CPU**: <1% when idle
- **Network**: Minimal (only health checks)

### Can I monitor thousands of ports?

Yes, but consider:
- Network bandwidth
- Check frequency (via external monitoring)
- Target service load

## Security

### Is PortGuard secure?

PortGuard follows security best practices:
- No authentication by default (use firewall/proxy)
- Minimal attack surface
- No data storage
- Runs with minimal privileges

### Should I expose PortGuard to the internet?

**No!** Use it behind:
- Firewall rules
- VPN
- Reverse proxy with authentication
- Load balancer

### How do I report security issues?

See [SECURITY.md](SECURITY.md) for our security policy.

## Development

### How do I contribute?

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### How do I build from source?

```bash
git clone https://github.com/mrwogu/portguard.git
cd portguard
go build -o portguard main.go
```

Or use Make:
```bash
make build
```

### How do I run tests?

```bash
make test
```

### Where can I get help?

- 📝 [File an Issue](https://github.com/mrwogu/portguard/issues/new)
- 💬 [Discussions](https://github.com/mrwogu/portguard/discussions)
- 📖 [Documentation](https://github.com/mrwogu/portguard)

## Still have questions?

If your question isn't answered here:

1. Check the [README](README.md)
2. Search [existing issues](https://github.com/mrwogu/portguard/issues)
3. Ask in [Discussions](https://github.com/mrwogu/portguard/discussions)
4. [Open a new issue](https://github.com/mrwogu/portguard/issues/new)
