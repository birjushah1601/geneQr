# ABY-MED Deployment Guide

Quick deployment reference. For detailed guide, see [DEPLOYMENT-GUIDE.md](./DEPLOYMENT-GUIDE.md) and [PRODUCTION-DEPLOYMENT-CHECKLIST.md](./PRODUCTION-DEPLOYMENT-CHECKLIST.md).

---

## 🚀 Quick Deploy

### Prerequisites
- Go 1.21+, Node.js 18+, PostgreSQL 15+
- Domain with SSL certificate
- SMTP credentials (SendGrid)
- AI API keys (OpenAI, Anthropic)

### Backend Deployment
```bash
# Build
CGO_ENABLED=0 go build -o platform ./cmd/platform

# Run
./platform
```

### Frontend Deployment
```bash
cd admin-ui
npm run build
npm start
```

### Database Setup
```bash
# Apply migrations
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f database/migrations/*.sql
```

---

## 🔧 Environment Configuration

### Required Variables
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=<secure_password>
DB_NAME=med_platform

# Application
PORT=8081
APP_ENV=production

# AI Keys
AI_OPENAI_API_KEY=sk-...
OPENAI_API_KEY=sk-...

# Email
SENDGRID_API_KEY=SG...
SENDGRID_FROM_EMAIL=noreply@aby-med.com

# Security
CORS_ALLOWED_ORIGINS=https://app.aby-med.com
```

---

## 📊 Health Checks

```bash
# Backend health
curl http://localhost:8081/health

# Database connection
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1;"
```

---

## 🐳 Docker Deployment

```bash
# Build images
docker build -t aby-med-backend .
docker build -t aby-med-frontend ./admin-ui

# Run with compose
docker-compose -f docker-compose.prod.yml up -d
```

---

## 📚 Full Documentation

- **Detailed Guide:** [DEPLOYMENT-GUIDE.md](./DEPLOYMENT-GUIDE.md)
- **Production Checklist:** [PRODUCTION-DEPLOYMENT-CHECKLIST.md](./PRODUCTION-DEPLOYMENT-CHECKLIST.md)
- **External Services:** [EXTERNAL-SERVICES-SETUP.md](./EXTERNAL-SERVICES-SETUP.md)

**Last Updated:** December 23, 2025
