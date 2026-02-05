# ServQR Production Deployment Architecture

## 🏗️ System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Internet / Users                            │
└────────────────────────────────┬────────────────────────────────────┘
                                 │
                    ┌────────────▼─────────────┐
                    │   Firewall (UFW/firewalld) │
                    │   Ports: 80, 443, 22      │
                    └────────────┬─────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
         │            ┌──────────▼──────────┐            │
         │            │  Nginx Reverse Proxy │            │
         │            │  - SSL Termination   │            │
         │            │  - Load Balancing    │            │
         │            │  - Static Assets     │            │
         │            └──────────┬──────────┘            │
         │                       │                       │
    ┌────▼─────┐          ┌─────▼──────┐          ┌────▼─────┐
    │ Frontend │          │  Backend   │          │  Static  │
    │ Next.js  │◄────────►│  Go API    │          │  Assets  │
    │ Port 3000│          │ Port 8081  │          │   CDN    │
    └────┬─────┘          └─────┬──────┘          └──────────┘
         │                      │
         │              ┌───────▼────────┐
         │              │  Docker Engine  │
         │              └───────┬────────┘
         │                      │
         │              ┌───────▼────────────┐
         └──────────────► PostgreSQL 15     │
                        │  Container        │
                        │  Port 5432        │
                        │  Data: /opt/servqr/data/postgres │
                        └───────┬───────────┘
                                │
                        ┌───────▼───────┐
                        │  Persistent   │
                        │    Storage    │
                        │  - Database   │
                        │  - Backups    │
                        │  - Logs       │
                        │  - QR Codes   │
                        └───────────────┘
```

## 📦 Component Breakdown

### 1. Frontend (Next.js)
**Location:** `/opt/servqr/admin-ui`  
**Service:** `servqr-frontend.service`  
**Port:** 3000 (internal)  
**Technology:** Next.js 14, React 18, TypeScript

**Responsibilities:**
- User interface rendering
- Client-side routing
- API communication
- Real-time updates (React Query)
- Authentication UI

**Resources:**
- Memory: ~512 MB
- CPU: 1 core
- Disk: ~200 MB (built assets)

### 2. Backend (Go API)
**Location:** `/opt/servqr/platform` (binary)  
**Service:** `servqr-backend.service`  
**Port:** 8081 (internal)  
**Technology:** Go 1.23, Chi router, PostgreSQL

**Responsibilities:**
- RESTful API endpoints
- Business logic
- Database operations
- Authentication/Authorization
- AI integration
- Email notifications
- WhatsApp integration

**Resources:**
- Memory: ~256 MB
- CPU: 1 core
- Disk: ~50 MB (binary)

### 3. Database (PostgreSQL)
**Location:** Docker container `servqr-postgres`  
**Service:** `servqr-postgres.service`  
**Port:** 5432 (internal only)  
**Technology:** PostgreSQL 15 Alpine

**Responsibilities:**
- Data persistence
- Multi-tenant isolation
- ACID transactions
- Full-text search
- Audit logging

**Resources:**
- Memory: ~512 MB - 1 GB
- CPU: 1 core
- Disk: 10-20 GB (data)

### 4. Nginx (Reverse Proxy)
**Location:** `/etc/nginx/sites-available/servqr`  
**Service:** `nginx.service`  
**Ports:** 80, 443  
**Technology:** Nginx 1.18+

**Responsibilities:**
- SSL termination
- Reverse proxy (frontend/backend)
- Static file serving
- Load balancing (future)
- Rate limiting
- Compression

## 🔄 Request Flow

### User Request Flow
```
1. User → https://yourdomain.com
2. Nginx → SSL termination
3. Nginx → Forward to Frontend (localhost:3000)
4. Frontend → Render page
5. Frontend → API call to /api/v1/endpoint
6. Nginx → Forward to Backend (localhost:8081)
7. Backend → Process request
8. Backend → Query PostgreSQL
9. PostgreSQL → Return data
10. Backend → Return JSON response
11. Frontend → Update UI
12. User → See result
```

### Background Job Flow
```
1. Backend → Event trigger (ticket created)
2. Backend → Queue notification job
3. Background Worker → Process job
4. SendGrid API → Send email
5. Database → Log notification
```

## 📁 Directory Structure

```
/opt/servqr/
├── deployment/                # Deployment scripts and configs
│   ├── README.md             # Main deployment guide
│   ├── QUICKSTART.md         # Quick start guide
│   ├── DEPLOYMENT-CHECKLIST.md
│   ├── deploy-all.sh         # Master deployment script
│   ├── install-prerequisites.sh
│   ├── setup-docker.sh
│   ├── deploy-app.sh
│   ├── backup-database.sh    # Generated by setup-docker.sh
│   ├── restore-database.sh   # Generated by setup-docker.sh
│   ├── connect-database.sh   # Generated by setup-docker.sh
│   ├── docker-compose.yml    # Generated by setup-docker.sh
│   ├── .env.production.template
│   └── systemd/
│       ├── servqr-backend.service
│       └── servqr-frontend.service
│
├── cmd/platform/main.go      # Backend entry point
├── platform                  # Compiled backend binary
├── .env                      # Backend configuration (generated)
├── .db_password              # Database password (generated)
├── .jwt_secret               # JWT secret (generated)
│
├── admin-ui/                 # Frontend source
│   ├── .next/               # Built frontend (generated)
│   ├── .env.local           # Frontend config (generated)
│   └── ...
│
├── data/                     # Runtime data
│   ├── postgres/            # PostgreSQL data (persistent)
│   ├── qrcodes/             # Generated QR codes
│   └── whatsapp/            # WhatsApp media
│
├── logs/                     # Application logs
│   ├── deployment-*.log     # Deployment logs
│   ├── backend.log          # Backend logs
│   ├── backend-error.log    # Backend errors
│   ├── frontend.log         # Frontend logs
│   └── frontend-error.log   # Frontend errors
│
├── backups/                  # Database backups
│   └── servqr-backup-*.sql.gz
│
├── storage/                  # File uploads
│   ├── attachments/         # Ticket attachments
│   └── qr_codes/            # QR code images
│
├── migrations/               # Database migrations
│   └── *.sql
│
└── docs/                     # Documentation
    └── ...
```

## 🔐 Security Layers

### 1. Network Security
- Firewall (UFW/firewalld) - Only ports 80, 443, 22 open
- Private internal network for services
- Database not exposed externally

### 2. Application Security
- HTTPS only (SSL/TLS)
- JWT authentication
- CORS policy enforcement
- Rate limiting (100 req/min)
- Input sanitization
- SQL injection prevention (prepared statements)

### 3. Data Security
- Database password encrypted at rest
- JWT secret secure random generation
- Audit logging all operations
- Multi-tenant data isolation
- Role-based access control (RBAC)

### 4. System Security
- systemd service isolation
- Read-only file systems where possible
- No new privileges flag
- Private /tmp directories
- Resource limits (CPU, memory, files)

## 🔄 Service Management

### systemd Services
```
servqr-postgres.service    → PostgreSQL (Docker)
  ├── servqr-backend.service    → Go API
  │     └── servqr-frontend.service   → Next.js
  │
  └── nginx.service          → Reverse Proxy
```

**Dependency Chain:**
1. Docker starts
2. PostgreSQL container starts
3. Backend waits for database
4. Frontend starts after backend
5. Nginx routes traffic

### Auto-Start on Boot
All services configured to start automatically:
```bash
sudo systemctl is-enabled servqr-postgres  # enabled
sudo systemctl is-enabled servqr-backend   # enabled
sudo systemctl is-enabled servqr-frontend  # enabled
sudo systemctl is-enabled nginx            # enabled
```

## 💾 Data Persistence

### PostgreSQL Data
**Location:** `/opt/servqr/data/postgres`  
**Mount:** Docker volume mounted to host  
**Backup:** Daily automated backups at 2 AM  
**Retention:** 7 days

### Application Storage
**QR Codes:** `/opt/servqr/data/qrcodes`  
**WhatsApp Media:** `/opt/servqr/data/whatsapp`  
**File Uploads:** `/opt/servqr/storage`

### Logs
**System Logs:** `journalctl` (systemd)  
**Application Logs:** `/opt/servqr/logs/`  
**Rotation:** Daily, keep 14 days

## 🔍 Monitoring Points

### Health Checks
1. **Backend:** `GET /health` → `{"status":"ok"}`
2. **Frontend:** `GET /` → HTTP 200
3. **Database:** `pg_isready -U servqr`
4. **Nginx:** `systemctl status nginx`

### Metrics to Monitor
- CPU usage (target: <70%)
- Memory usage (target: <80%)
- Disk space (target: >20% free)
- Database connections (target: <50)
- API response time (target: <500ms)
- Error rate (target: <1%)

### Log Monitoring
- Backend errors: `sudo journalctl -u servqr-backend | grep ERROR`
- Frontend errors: `sudo journalctl -u servqr-frontend | grep ERROR`
- Database errors: `docker logs servqr-postgres | grep ERROR`
- Nginx errors: `tail -f /var/log/nginx/error.log`

## 🚀 Scaling Strategy

### Vertical Scaling (Current)
- Increase VM resources (CPU, RAM)
- Tune PostgreSQL parameters
- Optimize application code

### Horizontal Scaling (Future)
```
Load Balancer
  ├── Backend Server 1
  ├── Backend Server 2
  └── Backend Server 3
        ↓
  Database (Primary + Replicas)
```

**Steps for horizontal scaling:**
1. Add load balancer (Nginx/HAProxy)
2. Deploy multiple backend instances
3. Add PostgreSQL read replicas
4. Implement Redis for sessions/cache
5. Use CDN for static assets

## 📊 Resource Planning

### Small Deployment (100-500 users)
- VM: 4 GB RAM, 2 CPU cores
- Database: 10 GB storage
- Bandwidth: 100 GB/month

### Medium Deployment (500-5000 users)
- VM: 8 GB RAM, 4 CPU cores
- Database: 50 GB storage
- Bandwidth: 500 GB/month
- Consider: Database replica, Redis cache

### Large Deployment (5000+ users)
- VMs: Multiple instances (4x 8 GB RAM, 4 CPU)
- Database: Dedicated server (16 GB RAM, 100+ GB)
- Load balancer
- CDN for static assets
- Redis cluster
- Monitoring stack (Prometheus + Grafana)

## 🔧 Maintenance Windows

### Daily
- Automated database backups (2 AM)
- Log rotation (midnight)

### Weekly
- Review logs for errors
- Check disk space
- Monitor performance metrics

### Monthly
- System package updates
- Security patches
- Database vacuum/analyze
- Backup restoration test

### Quarterly
- Full system audit
- Performance tuning
- Capacity planning review
- Disaster recovery drill

## 🆘 Disaster Recovery

### Backup Strategy
- **Frequency:** Daily
- **Retention:** 7 days
- **Location:** Local + offsite (recommended)
- **RTO:** 1 hour
- **RPO:** 24 hours

### Recovery Procedure
1. Provision new server
2. Run deployment scripts
3. Restore latest database backup
4. Update DNS (if needed)
5. Verify all services
6. Monitor for 24 hours

---

**Document Version:** 1.0  
**Last Updated:** 2026-02-06  
**Maintained By:** DevOps Team
