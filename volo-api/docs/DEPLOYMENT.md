# Deployment

## Overview

Everything runs on a single VM. One `docker compose up -d` brings up the full stack. Caddy handles TLS automatically via Let's Encrypt.

```
┌─────────────────────────────────────────────────────────────┐
│  Your VM                                                     │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐  │
│  │  Docker Compose                                         │  │
│  │                                                         │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐             │  │
│  │  │ volo-api │  │ postgres │  │  redis   │             │  │
│  │  │  :8080   │  │  :5432   │  │  :6379   │             │  │
│  │  └──────────┘  └──────────┘  └──────────┘             │  │
│  │                                                         │  │
│  │  ┌──────────────────────────────────────────────────┐   │  │
│  │  │  Caddy (reverse proxy + auto TLS)                 │   │  │
│  │  │  :80, :443                                        │   │  │
│  │  │                                                    │   │  │
│  │  │  api.volo.domain → volo-api:8080                   │   │  │
│  │  │  volo.domain → /var/www/landing (static)           │   │  │
│  │  └──────────────────────────────────────────────────┘   │  │
│  │                                                         │  │
│  │  ┌──────────────────────────────────────────────────┐   │  │
│  │  │  Prometheus + Grafana (observability)              │   │  │
│  │  │  prometheus:9090 | grafana:3000                    │   │  │
│  │  └──────────────────────────────────────────────────┘   │  │
│  └────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## VM Requirements

| Resource | Minimum | Recommended |
|----------|---------|-------------|
| CPU | 1 vCPU | 2 vCPU |
| RAM | 1 GB | 2 GB |
| Storage | 10 GB | 20 GB |
| OS | Any Linux | Ubuntu 22.04 LTS |
| Ports | 80, 443 | 80, 443, 3000 (Grafana) |

---

## Domain Setup

Two subdomains pointing to the VM's IP:

```
api.volo.yourdomain.com  → VM IP (A record)
volo.yourdomain.com      → VM IP (A record)
```

Caddy handles TLS certificates automatically — no manual cert management.

---

## Directory Structure on VM

```
/opt/volo/
├── docker-compose.prod.yml
├── Caddyfile
├── .env                        ← secrets (not in git)
├── api/
│   ├── Dockerfile
│   ├── cmd/
│   ├── internal/
│   ├── migrations/
│   └── models/                 ← ONNX model files
├── landing/
│   └── dist/                   ← built static files
├── prometheus/
│   └── prometheus.yml
├── grafana/
│   └── provisioning/
└── data/
    ├── postgres/               ← persistent volume
    └── redis/                  ← persistent volume
```

---

## Production docker-compose

```yaml
version: "3.8"

services:
  volo-api:
    build: ./api
    restart: unless-stopped
    environment:
      - VOLO_PORT=8080
      - VOLO_DATABASE_URL=postgres://${DB_USER}:${DB_PASS}@postgres:5432/${DB_NAME}?sslmode=disable
      - VOLO_REDIS_ADDR=redis:6379
      - VOLO_REDIS_PASSWORD=${REDIS_PASS}
      - VOLO_JWT_SECRET=${JWT_SECRET}
      - VOLO_CORS_ORIGINS=https://volo.yourdomain.com,chrome-extension://${EXT_ID}
      - VOLO_LOG_LEVEL=info
      - VOLO_GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID}
      - VOLO_GOOGLE_CLIENT_SECRET=${GOOGLE_CLIENT_SECRET}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - internal

  postgres:
    image: postgres:16-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASS}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - ./data/postgres:/var/lib/postgresql/data
      - ./api/migrations/001_initial.sql:/docker-entrypoint-initdb.d/001_initial.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - internal

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    command: redis-server --requirepass ${REDIS_PASS}
    volumes:
      - ./data/redis:/data
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "${REDIS_PASS}", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - internal

  caddy:
    image: caddy:2-alpine
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - ./landing/dist:/var/www/landing
      - caddy_data:/data
      - caddy_config:/config
    networks:
      - internal

  prometheus:
    image: prom/prometheus:latest
    restart: unless-stopped
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - internal

  grafana:
    image: grafana/grafana:latest
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_PASS}
    volumes:
      - grafana_data:/var/lib/grafana
    networks:
      - internal

volumes:
  caddy_data:
  caddy_config:
  grafana_data:

networks:
  internal:
```

---

## Caddyfile

```
api.volo.yourdomain.com {
    reverse_proxy volo-api:8080
    header {
        X-Content-Type-Options nosniff
        X-Frame-Options DENY
        Referrer-Policy strict-origin-when-cross-origin
    }
}

volo.yourdomain.com {
    root * /var/www/landing
    file_server
    header {
        X-Content-Type-Options nosniff
        Cache-Control "public, max-age=3600"
    }
}
```

---

## Environment File (.env)

```bash
# Database
DB_USER=volo
DB_PASS=<strong-random-password>
DB_NAME=volo

# Redis
REDIS_PASS=<strong-random-password>

# Auth
JWT_SECRET=<64-char-random-string>
GOOGLE_CLIENT_ID=<from-google-console>
GOOGLE_CLIENT_SECRET=<from-google-console>

# Extension
EXT_ID=<chrome-extension-id-after-publishing>

# Grafana
GRAFANA_PASS=<admin-password>
```

Generate secrets:
```bash
openssl rand -hex 32   # for JWT_SECRET
openssl rand -hex 16   # for DB_PASS, REDIS_PASS
```

---

## Deployment Steps

### First Time

```bash
# 1. SSH into VM
ssh user@your-vm-ip

# 2. Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# 3. Clone repo
git clone <your-repo> /opt/volo
cd /opt/volo

# 4. Create .env file
cp .env.example .env
nano .env  # fill in secrets

# 5. Build landing page (on your machine, then copy)
# On local: cd volo-landing && npm run build
# Then: scp -r dist/ user@vm:/opt/volo/landing/dist/

# 6. Start everything
docker compose -f docker-compose.prod.yml up -d --build

# 7. Verify
curl https://api.volo.yourdomain.com/health
# → {"data":{"status":"ok"}}
```

### Updates

```bash
# Pull latest code
cd /opt/volo && git pull

# Rebuild and restart API (zero-downtime with rolling update)
docker compose -f docker-compose.prod.yml up -d --build volo-api

# Update landing page
scp -r dist/ user@vm:/opt/volo/landing/dist/
# Caddy serves static files — no restart needed

# Update ONNX model (hot-swap)
scp models/intent.onnx user@vm:/opt/volo/api/models/
docker compose -f docker-compose.prod.yml restart volo-api
```

### Rollback

```bash
# If something breaks, revert to previous image
docker compose -f docker-compose.prod.yml down volo-api
git checkout HEAD~1
docker compose -f docker-compose.prod.yml up -d --build volo-api
```

---

## Monitoring (Prometheus + Grafana)

### Prometheus Config

```yaml
# prometheus/prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'volo-api'
    static_configs:
      - targets: ['volo-api:8080']
    metrics_path: '/metrics'
```

### Metrics Exposed by API

| Metric | Type | Description |
|--------|------|-------------|
| `volo_http_requests_total` | Counter | Total requests by method, path, status |
| `volo_http_request_duration_seconds` | Histogram | Request latency |
| `volo_commands_processed_total` | Counter | Commands processed by action type |
| `volo_auth_registrations_total` | Counter | New device registrations |
| `volo_rate_limit_hits_total` | Counter | Rate limit rejections |

### Grafana Dashboards

Access at `http://your-vm-ip:3000` (or behind Caddy if you add a subdomain).

Suggested panels:
- Request rate (req/s) over time
- P50/P95/P99 latency
- Error rate (4xx, 5xx)
- Active users (unique user_ids in last hour)
- Command breakdown by action type
- Rate limit hits

---

## Backups

### Postgres

```bash
# Daily backup via cron
0 3 * * * docker exec postgres pg_dump -U volo volo | gzip > /opt/volo/backups/volo-$(date +\%Y\%m\%d).sql.gz

# Keep last 7 days
find /opt/volo/backups -name "*.sql.gz" -mtime +7 -delete
```

### Redis

Redis persistence is handled via RDB snapshots (default config). Volume is mounted at `./data/redis`.

---

## Security Checklist

- [ ] `.env` file has 600 permissions (`chmod 600 .env`)
- [ ] Postgres not exposed to public (only on internal Docker network)
- [ ] Redis not exposed to public (only on internal Docker network)
- [ ] JWT secret is 64+ characters
- [ ] CORS origins locked to specific domains (not `*`)
- [ ] Caddy security headers enabled (X-Frame-Options, etc.)
- [ ] SSH key-only auth on VM (no password login)
- [ ] Firewall: only 80, 443, 22 open
- [ ] Regular OS updates (`unattended-upgrades`)
- [ ] Grafana behind auth (not publicly accessible without password)

---

## Future Considerations

- **CI/CD:** GitHub Actions → build Docker image → push to registry → SSH deploy
- **Multiple instances:** If traffic grows, add a load balancer and run 2+ API containers
- **Managed DB:** Move Postgres to a managed service (e.g., Supabase, Neon, RDS) for automatic backups and scaling
- **CDN:** Put landing page behind Cloudflare for global edge caching
- **Log aggregation:** Add Loki + Grafana for centralized log search
