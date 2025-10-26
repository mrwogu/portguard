# PortGuard Example Configurations

This directory contains example configurations for various use cases.

## Basic Mail Server

Monitor essential mail server ports:

```yaml
server:
  port: "8888"
  timeout: 2s

checks:
  - host: "mail.example.com"
    port: 25
    name: "SMTP"
    description: "Mail Transfer Protocol"
  
  - host: "mail.example.com"
    port: 587
    name: "SMTP Submission"
    description: "Mail Submission"
  
  - host: "mail.example.com"
    port: 993
    name: "IMAPS"
    description: "IMAP over SSL"
```

## Web Application Stack

Monitor a complete web application:

```yaml
server:
  port: "8888"
  timeout: 3s

checks:
  # Web Server
  - host: "localhost"
    port: 80
    name: "HTTP"
    description: "Web Server HTTP"
  
  - host: "localhost"
    port: 443
    name: "HTTPS"
    description: "Web Server HTTPS"
  
  # Application Server
  - host: "localhost"
    port: 8080
    name: "App Server"
    description: "Application Backend"
  
  # Database
  - host: "localhost"
    port: 5432
    name: "PostgreSQL"
    description: "Database Server"
  
  # Cache
  - host: "localhost"
    port: 6379
    name: "Redis"
    description: "Cache Server"
  
  # Message Queue
  - host: "localhost"
    port: 5672
    name: "RabbitMQ"
    description: "Message Broker"
```

## Kubernetes Cluster Services

Monitor services in a Kubernetes cluster:

```yaml
server:
  port: "8888"
  timeout: 5s

checks:
  # API Server
  - host: "kubernetes.default.svc"
    port: 443
    name: "K8s API"
    description: "Kubernetes API Server"
  
  # Ingress Controller
  - host: "ingress-nginx-controller"
    port: 80
    name: "Ingress HTTP"
    description: "Ingress Controller HTTP"
  
  - host: "ingress-nginx-controller"
    port: 443
    name: "Ingress HTTPS"
    description: "Ingress Controller HTTPS"
  
  # Service Mesh
  - host: "istio-ingressgateway"
    port: 15021
    name: "Istio Health"
    description: "Istio Health Port"
```

## Database Cluster

Monitor a database cluster:

```yaml
server:
  port: "8888"
  timeout: 2s

checks:
  # PostgreSQL Primary
  - host: "pg-primary.example.com"
    port: 5432
    name: "PG Primary"
    description: "PostgreSQL Primary"
  
  # PostgreSQL Replicas
  - host: "pg-replica-1.example.com"
    port: 5432
    name: "PG Replica 1"
    description: "PostgreSQL Replica 1"
  
  - host: "pg-replica-2.example.com"
    port: 5432
    name: "PG Replica 2"
    description: "PostgreSQL Replica 2"
  
  # PgBouncer Connection Pooler
  - host: "pgbouncer.example.com"
    port: 6432
    name: "PgBouncer"
    description: "Connection Pooler"
```

## Microservices Architecture

Monitor multiple microservices:

```yaml
server:
  port: "8888"
  timeout: 3s

checks:
  # API Gateway
  - host: "api-gateway"
    port: 8080
    name: "API Gateway"
    description: "Main API Gateway"
  
  # User Service
  - host: "user-service"
    port: 8081
    name: "User Service"
    description: "User Management"
  
  # Auth Service
  - host: "auth-service"
    port: 8082
    name: "Auth Service"
    description: "Authentication Service"
  
  # Order Service
  - host: "order-service"
    port: 8083
    name: "Order Service"
    description: "Order Processing"
  
  # Payment Service
  - host: "payment-service"
    port: 8084
    name: "Payment Service"
    description: "Payment Processing"
  
  # Notification Service
  - host: "notification-service"
    port: 8085
    name: "Notification Service"
    description: "Notification Service"
```

## Network Infrastructure

Monitor network infrastructure components:

```yaml
server:
  port: "8888"
  timeout: 5s

checks:
  # DNS
  - host: "ns1.example.com"
    port: 53
    name: "DNS Primary"
    description: "Primary DNS Server"
  
  # LDAP
  - host: "ldap.example.com"
    port: 389
    name: "LDAP"
    description: "Directory Service"
  
  - host: "ldap.example.com"
    port: 636
    name: "LDAPS"
    description: "LDAP over SSL"
  
  # VPN
  - host: "vpn.example.com"
    port: 1194
    name: "OpenVPN"
    description: "VPN Server"
  
  # Proxy
  - host: "proxy.example.com"
    port: 3128
    name: "Squid Proxy"
    description: "HTTP Proxy"
```

## High Timeout for Remote Services

For monitoring remote or slow services:

```yaml
server:
  port: "8888"
  timeout: 10s  # Longer timeout for remote services

checks:
  - host: "remote-api.partner.com"
    port: 443
    name: "Partner API"
    description: "External Partner API"
  
  - host: "backup-server.remote.com"
    port: 22
    name: "Backup SSH"
    description: "Remote Backup Server"
```

## Localhost Development

For local development monitoring:

```yaml
server:
  port: "8888"
  timeout: 1s

checks:
  - host: "localhost"
    port: 3000
    name: "React Dev"
    description: "React Development Server"
  
  - host: "localhost"
    port: 5000
    name: "API Dev"
    description: "Backend API"
  
  - host: "localhost"
    port: 5432
    name: "Dev DB"
    description: "Development Database"
  
  - host: "localhost"
    port: 6379
    name: "Dev Redis"
    description: "Development Cache"
```
