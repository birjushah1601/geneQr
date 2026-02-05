# Log Management & Viewing Setup Guide

**Status:** 📋 BRAINSTORMING / PLANNING  
**Target:** Web-based log viewing and management  
**Last Updated:** February 6, 2026

---

## 🎯 Overview

Setup web-based log management tools to view, search, and analyze logs from anywhere through a web browser. Access secured through Nginx reverse proxy.

---

## 🔧 Options for Log Management

### Option 1: Dozzle (Simple, Docker-focused) ⭐ **RECOMMENDED for Docker**

**Pros:**
- ✅ Lightweight (single container)
- ✅ Real-time Docker logs
- ✅ Multi-container support
- ✅ Simple web UI
- ✅ No configuration needed
- ✅ Search & filter
- ✅ Container stats

**Cons:**
- ❌ Docker logs only (not file logs)
- ❌ No log retention/storage
- ❌ No alerting

### Option 2: Grafana Loki + Promtail (Powerful) 🔥 **RECOMMENDED for Production**

**Pros:**
- ✅ Powerful query language (LogQL)
- ✅ Long-term storage
- ✅ Beautiful Grafana dashboards
- ✅ Alert integration
- ✅ Multi-source logs (files + Docker)
- ✅ Label-based organization

**Cons:**
- ❌ More complex setup
- ❌ Higher resource usage

### Option 3: Graylog (Enterprise-grade)

**Pros:**
- ✅ Full-featured SIEM
- ✅ Advanced search
- ✅ Alerting & dashboards
- ✅ User management

**Cons:**
- ❌ Resource heavy
- ❌ Complex setup
- ❌ Overkill for small deployments

### Option 4: ELK Stack (Elasticsearch + Logstash + Kibana)

**Pros:**
- ✅ Industry standard
- ✅ Powerful analytics
- ✅ Rich visualizations

**Cons:**
- ❌ Very resource heavy
- ❌ Complex configuration
- ❌ Expensive for production

---

## 🚀 Recommended Setup: Grafana Loki + Promtail

This provides the best balance of features, performance, and ease of use.

---

## 📦 Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   Nginx (Port 443)                          │
│  /logs → Grafana (Port 3001)                               │
└─────────────┬───────────────────────────────────────────────┘
              │
    ┌─────────▼─────────┐
    │    Grafana        │
    │   (Dashboards)    │
    │    Port: 3001     │
    └─────────┬─────────┘
              │
    ┌─────────▼─────────┐
    │      Loki         │
    │  (Log Storage)    │
    │    Port: 3100     │
    └─────────┬─────────┘
              │
    ┌─────────▼─────────┐
    │    Promtail       │
    │  (Log Collector)  │
    └───────────────────┘
              │
    ┌─────────▼─────────────────────┐
    │  Log Files & Docker Logs      │
    │  - /var/servqr/logs/          │
    │  - Docker containers          │
    └───────────────────────────────┘
```

---

## 🔧 Setup Instructions

### Option A: Docker Compose (Recommended)

Add to `docker-compose.yml`:

```yaml
services:
  # ... existing services ...

  # Loki - Log aggregation
  loki:
    image: grafana/loki:2.9.3
    container_name: servqr-loki
    restart: unless-stopped
    ports:
      - "127.0.0.1:3100:3100"
    command: -config.file=/etc/loki/local-config.yaml
    volumes:
      - ./loki/loki-config.yaml:/etc/loki/local-config.yaml
      - /var/servqr/loki-data:/loki
    networks:
      - servqr-network
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  # Promtail - Log collector
  promtail:
    image: grafana/promtail:2.9.3
    container_name: servqr-promtail
    restart: unless-stopped
    volumes:
      - ./promtail/promtail-config.yaml:/etc/promtail/config.yaml
      - /var/servqr/logs:/var/log/servqr:ro
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
      - /var/run/docker.sock:/var/run/docker.sock
    command: -config.file=/etc/promtail/config.yaml
    depends_on:
      - loki
    networks:
      - servqr-network
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  # Grafana - Visualization
  grafana:
    image: grafana/grafana:10.2.3
    container_name: servqr-grafana
    restart: unless-stopped
    ports:
      - "127.0.0.1:3001:3000"
    environment:
      - GF_SERVER_ROOT_URL=https://servqr.com/logs/
      - GF_SERVER_SERVE_FROM_SUB_PATH=true
      - GF_AUTH_ANONYMOUS_ENABLED=false
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_ADMIN_PASSWORD}
      - GF_USERS_ALLOW_SIGN_UP=false
      - GF_INSTALL_PLUGINS=grafana-piechart-panel
    volumes:
      - /var/servqr/grafana-data:/var/lib/grafana
      - ./grafana/provisioning:/etc/grafana/provisioning
    depends_on:
      - loki
    networks:
      - servqr-network
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

---

### Configuration Files

#### 1. Loki Configuration

Create `loki/loki-config.yaml`:

```yaml
auth_enabled: false

server:
  http_listen_port: 3100
  grpc_listen_port: 9096

common:
  path_prefix: /loki
  storage:
    filesystem:
      chunks_directory: /loki/chunks
      rules_directory: /loki/rules
  replication_factor: 1
  ring:
    instance_addr: 127.0.0.1
    kvstore:
      store: inmemory

query_range:
  results_cache:
    cache:
      embedded_cache:
        enabled: true
        max_size_mb: 100

schema_config:
  configs:
    - from: 2020-10-24
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h

ruler:
  alertmanager_url: http://localhost:9093

limits_config:
  retention_period: 744h  # 31 days
  ingestion_rate_mb: 16
  ingestion_burst_size_mb: 32
  max_query_series: 100000
  max_query_parallelism: 32

chunk_store_config:
  max_look_back_period: 744h

table_manager:
  retention_deletes_enabled: true
  retention_period: 744h
```

#### 2. Promtail Configuration

Create `promtail/promtail-config.yaml`:

```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  # Backend application logs
  - job_name: backend
    static_configs:
      - targets:
          - localhost
        labels:
          job: backend
          app: servqr
          __path__: /var/log/servqr/backend/*.log

  # Frontend application logs
  - job_name: frontend
    static_configs:
      - targets:
          - localhost
        labels:
          job: frontend
          app: servqr
          __path__: /var/log/servqr/frontend/*.log

  # Nginx logs
  - job_name: nginx
    static_configs:
      - targets:
          - localhost
        labels:
          job: nginx
          app: servqr
          __path__: /var/log/servqr/nginx/*.log

  # PostgreSQL logs
  - job_name: postgres
    static_configs:
      - targets:
          - localhost
        labels:
          job: postgres
          app: servqr
          __path__: /var/log/servqr/postgres/*.log

  # Docker container logs
  - job_name: docker
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        refresh_interval: 5s
    relabel_configs:
      - source_labels: ['__meta_docker_container_name']
        regex: '/(.*)'
        target_label: 'container'
      - source_labels: ['__meta_docker_container_log_stream']
        target_label: 'stream'
```

#### 3. Grafana Provisioning

Create `grafana/provisioning/datasources/loki.yaml`:

```yaml
apiVersion: 1

datasources:
  - name: Loki
    type: loki
    access: proxy
    url: http://loki:3100
    jsonData:
      maxLines: 1000
      derivedFields:
        - datasourceUid: loki
          matcherRegex: "traceID=(\\w+)"
          name: TraceID
          url: "$${__value.raw}"
```

Create `grafana/provisioning/dashboards/dashboard.yaml`:

```yaml
apiVersion: 1

providers:
  - name: 'ServQR Dashboards'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    updateIntervalSeconds: 10
    allowUiUpdates: true
    options:
      path: /etc/grafana/provisioning/dashboards
```

---

### 4. Nginx Configuration

Add to `nginx/conf.d/servqr.conf`:

```nginx
# Inside the main HTTPS server block

# Grafana (Logs Dashboard)
location /logs/ {
    # Authentication required
    auth_basic "Log Access";
    auth_basic_user_file /etc/nginx/.htpasswd;

    proxy_pass http://grafana:3000/;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection 'upgrade';
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_cache_bypass $http_upgrade;

    # WebSocket support for live logs
    proxy_read_timeout 86400;
}

# Loki API (for Grafana)
location /loki/ {
    # Authentication required
    auth_basic "Log Access";
    auth_basic_user_file /etc/nginx/.htpasswd;

    proxy_pass http://loki:3100/;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

---

### 5. Setup Authentication

```bash
# Create htpasswd file for basic auth
sudo apt install apache2-utils

# Create password file
sudo htpasswd -c /etc/nginx/.htpasswd logviewer

# Add more users
sudo htpasswd /etc/nginx/.htpasswd admin
```

Or in Docker:

```bash
# Create password in nginx container
docker compose exec nginx sh -c "apk add apache2-utils && htpasswd -c /etc/nginx/.htpasswd logviewer"
```

---

### 6. Update Environment Variables

Add to `.env.production`:

```bash
# Grafana
GRAFANA_ADMIN_PASSWORD=CHANGE_THIS_STRONG_PASSWORD
```

---

### 7. Create Volumes

```bash
# Create data directories
sudo mkdir -p /var/servqr/{loki-data,grafana-data}
sudo chown -R $USER:$USER /var/servqr/{loki-data,grafana-data}
```

---

### 8. Deploy

```bash
# Start new services
docker compose up -d loki promtail grafana

# Check status
docker compose ps

# View logs
docker compose logs -f grafana
```

---

## 🌐 Accessing Logs

### Via Web Browser

1. **Navigate to:** `https://servqr.com/logs/`
2. **Login:**
   - Username: `admin`
   - Password: (from GRAFANA_ADMIN_PASSWORD)

### Creating Dashboards

1. **Go to:** Explore → Select "Loki" datasource
2. **Query logs:**
   ```logql
   {job="backend"} |= "error"
   ```
3. **Common queries:**

```logql
# All backend errors
{job="backend"} |= "error"

# Specific container logs
{container="servqr-backend"}

# Last 5 minutes of nginx access logs
{job="nginx"} | logfmt | __error__=""

# Backend logs with specific pattern
{job="backend"} |~ "user.*login"

# Count errors per minute
sum(rate({job="backend"} |= "error" [1m])) by (level)
```

---

## 📊 Pre-built Dashboards

Create `grafana/provisioning/dashboards/servqr-logs.json`:

```json
{
  "dashboard": {
    "title": "ServQR Application Logs",
    "panels": [
      {
        "title": "Backend Errors (Last Hour)",
        "targets": [
          {
            "expr": "{job=\"backend\"} |= \"error\"",
            "refId": "A"
          }
        ]
      },
      {
        "title": "Frontend Logs",
        "targets": [
          {
            "expr": "{job=\"frontend\"}",
            "refId": "A"
          }
        ]
      },
      {
        "title": "Nginx Access Logs",
        "targets": [
          {
            "expr": "{job=\"nginx\", filename=\"/var/log/servqr/nginx/access.log\"}",
            "refId": "A"
          }
        ]
      },
      {
        "title": "Database Logs",
        "targets": [
          {
            "expr": "{job=\"postgres\"}",
            "refId": "A"
          }
        ]
      }
    ]
  }
}
```

---

## 🎯 Alternative: Simple Option (Dozzle)

If Grafana is too complex, use Dozzle for simple Docker log viewing:

### Docker Compose

```yaml
services:
  # Dozzle - Simple log viewer
  dozzle:
    image: amir20/dozzle:latest
    container_name: servqr-dozzle
    restart: unless-stopped
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    ports:
      - "127.0.0.1:8080:8080"
    environment:
      DOZZLE_LEVEL: info
      DOZZLE_TAILSIZE: 300
      DOZZLE_FILTER: "name=servqr-*"
      DOZZLE_NO_ANALYTICS: "true"
    networks:
      - servqr-network
```

### Nginx Configuration for Dozzle

```nginx
# Dozzle (Simple Docker logs)
location /dozzle/ {
    auth_basic "Log Access";
    auth_basic_user_file /etc/nginx/.htpasswd;

    proxy_pass http://dozzle:8080/;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection 'upgrade';
    proxy_set_header Host $host;
    proxy_cache_bypass $http_upgrade;
}
```

**Access:** `https://servqr.com/dozzle/`

---

## 🔒 Security Best Practices

1. **Always use authentication** (htpasswd or Grafana auth)
2. **Use HTTPS only** (never expose on HTTP)
3. **Restrict by IP** if possible:
   ```nginx
   location /logs/ {
       allow 192.168.1.0/24;  # Office network
       allow 10.0.0.0/8;       # VPN
       deny all;
       
       # ... rest of proxy config
   }
   ```
4. **Limit log retention** (31 days in Loki config)
5. **Regular backups** of Grafana dashboards
6. **Monitor disk space** (logs can grow quickly)

---

## 📊 Comparison Table

| Feature | Dozzle | Loki + Grafana | Graylog | ELK |
|---------|--------|----------------|---------|-----|
| **Setup Complexity** | ⭐ Easy | ⭐⭐ Medium | ⭐⭐⭐ Hard | ⭐⭐⭐⭐ Very Hard |
| **Resource Usage** | 50MB RAM | 500MB RAM | 2GB RAM | 4GB+ RAM |
| **Search** | Basic | Powerful | Advanced | Most Powerful |
| **Real-time** | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes |
| **Log Retention** | ❌ None | ✅ Configurable | ✅ Yes | ✅ Yes |
| **Dashboards** | ❌ No | ✅ Grafana | ✅ Built-in | ✅ Kibana |
| **Alerts** | ❌ No | ✅ Yes | ✅ Yes | ✅ Yes |
| **Cost** | Free | Free | Free/Paid | Free/Paid |

---

## 🎯 Recommendations

### For Small/Medium Deployments (< 1GB logs/day)
**Use:** Loki + Promtail + Grafana
- ✅ Good balance of features
- ✅ Reasonable resource usage
- ✅ Professional dashboards
- ✅ Future-proof

### For Simple Needs (Just Docker logs)
**Use:** Dozzle
- ✅ Super simple setup
- ✅ Minimal resources
- ✅ Real-time viewing
- ✅ Perfect for troubleshooting

### For Large Deployments (> 10GB logs/day)
**Use:** Graylog or ELK
- ✅ Enterprise features
- ✅ Advanced analytics
- ✅ Compliance features
- ⚠️ Requires dedicated resources

---

## 🚀 Quick Start Commands

```bash
# Setup Grafana + Loki + Promtail
cd /opt/servqr

# Create directories
mkdir -p loki promtail grafana/provisioning/{datasources,dashboards}
sudo mkdir -p /var/servqr/{loki-data,grafana-data}

# Copy configuration files (from this guide)
vim loki/loki-config.yaml
vim promtail/promtail-config.yaml
vim grafana/provisioning/datasources/loki.yaml

# Update docker-compose.yml (add services from this guide)

# Setup authentication
docker compose exec nginx sh -c "apk add apache2-utils && htpasswd -c /etc/nginx/.htpasswd admin"

# Deploy
docker compose up -d loki promtail grafana

# Access logs at:
# https://servqr.com/logs/
```

---

## 📱 Mobile Access

Grafana has a mobile-responsive interface, so you can view logs from:
- 📱 Phone browser
- 💻 Tablet
- 🖥️ Desktop

All through the same URL: `https://servqr.com/logs/`

---

## 🔍 Troubleshooting

### Grafana not showing logs
```bash
# Check Loki is running
docker compose logs loki

# Check Promtail is collecting
docker compose logs promtail

# Test Loki API
curl http://localhost:3100/ready
```

### Promtail not collecting logs
```bash
# Check file permissions
ls -la /var/servqr/logs/

# Check Promtail config
docker compose exec promtail cat /etc/promtail/config.yaml
```

### Can't access via Nginx
```bash
# Check Nginx config
docker compose exec nginx nginx -t

# Check authentication
cat /etc/nginx/.htpasswd
```

---

## 📋 Checklist

### Setup
- [ ] Loki + Promtail + Grafana added to docker-compose.yml
- [ ] Configuration files created
- [ ] Nginx reverse proxy configured
- [ ] Authentication setup (htpasswd)
- [ ] Data volumes created
- [ ] Services deployed and running

### Post-Setup
- [ ] Can access https://servqr.com/logs/
- [ ] Grafana login works
- [ ] Logs visible in Explore
- [ ] Created sample dashboard
- [ ] Tested search functionality
- [ ] Verified log retention settings

---

**Status:** 📋 READY FOR IMPLEMENTATION  
**Recommendation:** Use Loki + Grafana for production  
**Alternative:** Use Dozzle for quick/simple setup

---

## 📚 Related Documentation

- [Production Docker Setup Guide](./PRODUCTION-DOCKER-SETUP-GUIDE.md)
- [Production VM Setup Guide](./PRODUCTION-VM-SETUP-GUIDE.md)
