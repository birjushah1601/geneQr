# Production Deployment Guide - Attachment Service

This guide provides comprehensive instructions for deploying the ABY-MED medical platform attachment service to production environments.

## 🎯 Pre-Deployment Checklist

### System Requirements

- **Operating System**: Ubuntu 20.04+ / RHEL 8+ / Amazon Linux 2
- **CPU**: Minimum 2 cores, Recommended 4+ cores
- **RAM**: Minimum 4GB, Recommended 8GB+ 
- **Storage**: Minimum 50GB SSD, Recommended 200GB+ SSD
- **Network**: HTTPS (port 443) and SSH (port 22) access

### Prerequisites

- [x] Docker Engine 20.10+
- [x] Docker Compose v2.0+
- [x] PostgreSQL 14+ (managed service recommended)
- [x] OpenAI API key with vision model access
- [x] SSL/TLS certificate for domain
- [x] Domain name configured with DNS A record

## 🚀 Production Configuration

### 1. Environment Variables

Create a `.env.production` file with the following configuration:

```bash
# Application Configuration
NODE_ENV=production
GO_ENV=production
PORT=8080

# Database Configuration
DATABASE_URL=postgresql://username:password@host:5432/abymed_production?sslmode=require
DB_MAX_CONNECTIONS=50
DB_CONNECTION_TIMEOUT=30
DB_IDLE_TIMEOUT=300

# OpenAI Configuration
OPENAI_API_KEY=sk-your-production-openai-key-here
OPENAI_MODEL=gpt-4-vision-preview
OPENAI_MAX_TOKENS=2000
OPENAI_TEMPERATURE=0.3

# File Storage Configuration
STORAGE_TYPE=s3  # or 'local' for filesystem storage
STORAGE_PATH=/data/attachments
AWS_REGION=us-west-2
AWS_S3_BUCKET=abymed-attachments-prod
MAX_FILE_SIZE=100MB
ALLOWED_FILE_TYPES=image/jpeg,image/png,image/gif,image/bmp

# API Security
API_KEYS=prod-key-1:admin,prod-key-2:read,prod-key-3:upload
JWT_SECRET=your-256-bit-secret-key-here
CORS_ORIGINS=https://app.abymed.com,https://admin.abymed.com

# Rate Limiting
RATE_LIMIT_REQUESTS_PER_HOUR=1000
RATE_LIMIT_WINDOW_MINUTES=60
RATE_LIMIT_MAX_BURST=100

# Monitoring & Logging
LOG_LEVEL=info
LOG_FORMAT=json
ENABLE_METRICS=true
SENTRY_DSN=https://your-sentry-dsn@sentry.io/project-id

# SSL/TLS Configuration
SSL_CERT_PATH=/etc/ssl/certs/abymed.crt
SSL_KEY_PATH=/etc/ssl/private/abymed.key
FORCE_HTTPS=true

# Performance
GOMEMLIMIT=6GiB
GOMAXPROCS=4
```

### 2. Docker Production Configuration

Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  abymed-platform:
    image: abymed/medical-platform:latest
    restart: unless-stopped
    ports:
      - \"80:8080\"
      - \"443:8443\"
    environment:
      - NODE_ENV=production
    env_file:
      - .env.production
    volumes:
      - ./storage:/data/attachments
      - ./ssl:/etc/ssl
      - ./logs:/var/log/abymed
    networks:
      - abymed-network
    depends_on:
      - postgres
      - redis
    deploy:
      resources:
        limits:
          memory: 6G
          cpus: '4'
        reservations:
          memory: 2G
          cpus: '2'
    healthcheck:
      test: [\"CMD\", \"curl\", \"-f\", \"http://localhost:8080/health/attachments\"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  postgres:
    image: postgres:14
    restart: unless-stopped
    environment:
      POSTGRES_DB: abymed_production
      POSTGRES_USER: abymed_user
      POSTGRES_PASSWORD: secure_production_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups
    networks:
      - abymed-network
    deploy:
      resources:
        limits:
          memory: 2G
          cpus: '2'

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    command: redis-server --appendonly yes --maxmemory 512mb --maxmemory-policy allkeys-lru
    volumes:
      - redis_data:/data
    networks:
      - abymed-network
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '1'

  nginx:
    image: nginx:alpine
    restart: unless-stopped
    ports:
      - \"80:80\"
      - \"443:443\"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/ssl:ro
      - ./logs/nginx:/var/log/nginx
    depends_on:
      - abymed-platform
    networks:
      - abymed-network

volumes:
  postgres_data:
  redis_data:

networks:
  abymed-network:
    driver: bridge
```

### 3. Nginx Configuration

Create `nginx/nginx.conf`:

```nginx
events {
    worker_connections 1024;
}

http {
    upstream abymed_backend {
        server abymed-platform:8080 max_fails=3 fail_timeout=30s;
    }

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=upload:10m rate=2r/s;

    # Security headers
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection \"1; mode=block\" always;
    add_header Referrer-Policy strict-origin-when-cross-origin always;
    add_header Content-Security-Policy \"default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:\" always;

    # Redirect HTTP to HTTPS
    server {
        listen 80;
        server_name _;
        return 301 https://$host$request_uri;
    }

    # HTTPS Configuration
    server {
        listen 443 ssl http2;
        server_name api.abymed.com;

        ssl_certificate /etc/ssl/abymed.crt;
        ssl_certificate_key /etc/ssl/abymed.key;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;

        # File upload size limit
        client_max_body_size 100M;

        # API endpoints
        location /api/ {
            limit_req zone=api burst=20 nodelay;
            proxy_pass http://abymed_backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_read_timeout 300s;
            proxy_connect_timeout 75s;
        }

        # File upload endpoint (higher limits)
        location /api/v1/attachments {
            limit_req zone=upload burst=5 nodelay;
            proxy_pass http://abymed_backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_read_timeout 600s;
            proxy_connect_timeout 75s;
            proxy_request_buffering off;
        }

        # Health checks
        location /health {
            proxy_pass http://abymed_backend;
            access_log off;
        }

        # Metrics (restrict access)
        location /metrics {
            allow 10.0.0.0/8;
            allow 172.16.0.0/12;
            allow 192.168.0.0/16;
            deny all;
            proxy_pass http://abymed_backend;
        }
    }
}
```

## 📦 Deployment Steps

### Step 1: Server Setup

1. **Update system packages:**
```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl wget gnupg2 software-properties-common
```

2. **Install Docker:**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
newgrp docker
```

3. **Install Docker Compose:**
```bash
sudo curl -L \"https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)\" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### Step 2: Application Deployment

1. **Clone repository:**
```bash
git clone https://github.com/your-org/aby-med.git
cd aby-med
git checkout main
```

2. **Configure environment:**
```bash
cp .env.example .env.production
# Edit .env.production with production values
nano .env.production
```

3. **Set up SSL certificates:**
```bash
mkdir -p ssl
# Copy your SSL certificate files to ssl/ directory
sudo cp your-domain.crt ssl/abymed.crt
sudo cp your-domain.key ssl/abymed.key
sudo chmod 600 ssl/abymed.key
```

4. **Create required directories:**
```bash
mkdir -p storage logs/nginx backups
sudo chown -R 1000:1000 storage logs backups
```

5. **Deploy application:**
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### Step 3: Database Setup

1. **Run database migrations:**
```bash
docker-compose -f docker-compose.prod.yml exec abymed-platform ./migrate up
```

2. **Create initial API keys:**
```bash
docker-compose -f docker-compose.prod.yml exec abymed-platform ./abymed create-api-key --name=\"Admin Key\" --permissions=admin
```

### Step 4: Verification

1. **Check service health:**
```bash
curl -k https://your-domain.com/health/attachments
```

2. **Test file upload:**
```bash
curl -X POST https://your-domain.com/api/v1/attachments \
  -H \"X-API-Key: your-api-key\" \
  -F \"file=@test-image.jpg\"
```

3. **Monitor logs:**
```bash
docker-compose -f docker-compose.prod.yml logs -f abymed-platform
```

## 🔒 Security Configuration

### Firewall Setup

```bash
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### SSL/TLS Best Practices

- Use certificates from a trusted CA (Let's Encrypt recommended)
- Enable HTTP Strict Transport Security (HSTS)
- Configure perfect forward secrecy
- Disable weak SSL/TLS versions (< TLS 1.2)

### API Security

- Rotate API keys regularly (every 90 days)
- Use different keys for different environments
- Implement IP whitelisting for admin keys
- Monitor API usage patterns

## 📊 Monitoring & Alerts

### Health Check Monitoring

Set up external monitoring for these endpoints:
- `GET /health/attachments` - Overall system health
- `GET /metrics/attachments` - Performance metrics
- `GET /status/ai-analysis` - AI service status

### Log Analysis

Configure log aggregation for:
- Application logs (`/var/log/abymed/`)
- Nginx access logs (`/var/log/nginx/`)
- System logs (`/var/log/syslog`)

### Key Metrics to Monitor

- **Performance**: Response time, throughput, error rate
- **Resource Usage**: CPU, memory, disk space, network I/O
- **Business Metrics**: Upload success rate, AI analysis completion rate
- **Security**: Failed authentication attempts, rate limit violations

## 🔄 Backup & Recovery

### Database Backup Script

Create `/home/deploy/backup-db.sh`:

```bash
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR=\"/data/backups\"
DB_NAME=\"abymed_production\"

# Create backup
docker exec postgres pg_dump -U abymed_user $DB_NAME | gzip > $BACKUP_DIR/db_backup_$DATE.sql.gz

# Keep only last 30 days of backups
find $BACKUP_DIR -name \"db_backup_*.sql.gz\" -mtime +30 -delete

# Upload to S3 (optional)
# aws s3 cp $BACKUP_DIR/db_backup_$DATE.sql.gz s3://abymed-backups/database/
```

### File Storage Backup

- For S3: Enable versioning and cross-region replication
- For local storage: Use rsync to remote backup server

### Recovery Procedures

Document step-by-step recovery procedures for:
- Database corruption
- Application server failure
- Complete system failure
- Data loss scenarios

## 🚦 Scaling Considerations

### Horizontal Scaling

- Use load balancer (ALB/NLB) for multiple app instances
- Implement database read replicas
- Use Redis cluster for session management
- Consider CDN for file serving

### Vertical Scaling

- Monitor resource usage patterns
- Scale CPU/memory based on traffic
- Optimize database queries and indexes
- Implement file compression and caching

## 📋 Maintenance Tasks

### Daily Tasks
- Check system health dashboards
- Review error logs
- Monitor disk usage
- Verify backup completion

### Weekly Tasks
- Review performance metrics
- Update security patches
- Clean up old log files
- Test disaster recovery procedures

### Monthly Tasks
- Rotate API keys
- Review access logs for anomalies
- Update documentation
- Performance optimization review

## 🆘 Troubleshooting Guide

### Common Issues

1. **High Memory Usage**
   - Check for memory leaks in logs
   - Review file upload patterns
   - Optimize database queries

2. **Slow Response Times**
   - Check database connection pool
   - Monitor OpenAI API response times
   - Review nginx access logs

3. **File Upload Failures**
   - Verify disk space availability
   - Check file size limits
   - Review rate limiting settings

4. **AI Analysis Errors**
   - Verify OpenAI API key validity
   - Check API rate limits
   - Monitor API costs

### Emergency Contacts

- **DevOps Team**: devops@abymed.com
- **Database Admin**: dba@abymed.com
- **Security Team**: security@abymed.com

---

## ✅ Production Readiness Checklist

- [ ] Environment variables configured
- [ ] SSL certificates installed and valid
- [ ] Database migrations applied
- [ ] API keys generated and tested
- [ ] Health checks responding correctly
- [ ] Monitoring and alerting configured
- [ ] Backup procedures tested
- [ ] Security scans completed
- [ ] Load testing performed
- [ ] Documentation updated
- [ ] Team trained on deployment process

**🎉 Congratulations! Your ABY-MED attachment service is production-ready!**

For support, contact the development team or refer to the [API Documentation](../api/ATTACHMENT-API.md).