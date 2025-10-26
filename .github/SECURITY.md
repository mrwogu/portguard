# Security Policy

## About This Project

PortGuard is an open-source project created with passion. We do our best to maintain it, but as with any open-source software, there are no guarantees.

## Reporting a Vulnerability

If you discover a security vulnerability within PortGuard, **please open a GitHub issue** describing the problem.

**We encourage you to submit a pull request with a fix** - this is the fastest way to get the issue resolved. Community contributions are always welcome!

### What to Include

When reporting a vulnerability, please include:

* Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
* Full paths of source file(s) related to the manifestation of the issue
* The location of the affected source code (tag/branch/commit or direct URL)
* Any special configuration required to reproduce the issue
* Step-by-step instructions to reproduce the issue
* Proof-of-concept or exploit code (if possible)
* Impact of the issue, including how an attacker might exploit it

### No Guarantees

Please note that:

* We cannot guarantee that we will fix every reported issue
* Response times may vary depending on maintainer availability
* **The best way to ensure a fix is to submit a pull request yourself**

## Security Best Practices

When deploying PortGuard:

1. **Run with minimal privileges**: Don't run as root unless necessary
2. **Use HTTPS**: If exposing to the internet, use a reverse proxy with TLS
3. **Network isolation**: Run in isolated networks when possible
4. **Keep updated**: Always use the latest version
5. **Secure configuration**: Protect your config files with appropriate permissions
6. **Monitor logs**: Regularly review logs for suspicious activity

## Known Security Considerations

* PortGuard performs TCP connections to configured hosts - ensure firewall rules are appropriate
* Configuration files may contain sensitive information (hostnames, ports) - protect with file permissions
* The HTTP API is unauthenticated by design - use firewall rules or reverse proxy for access control
