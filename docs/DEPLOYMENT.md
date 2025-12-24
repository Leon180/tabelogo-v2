# Deployment Guide

## Prerequisites

- Docker & Docker Compose installed
- PostgreSQL (or use Docker)
- Redis (or use Docker)
- Go 1.24+ (for local development)

---

## Quick Start (Development)

### 1. Clone and Configure

```bash
# Clone repository
git clone <repository-url>
cd tabelogov2

# Copy environment template
cp .env.example .env

# Update .env with your values
# - JWT_SECRET (required)
# - GOOGLE_MAPS_API_KEY (required)
```

### 2. Start Services

```bash
# Start all services
make up

# Or manually
docker-compose -f deployments/docker-compose/docker-compose.yml up -d
```

### 3. Verify Deployment

```bash
# Check service health
curl http://localhost:8080/health  # Auth Service
curl http://localhost:18084/health # Spider Service
curl http://localhost:18082/health # Restaurant Service
curl http://localhost:8081/health  # Map Service

# Run integration tests
bash scripts/test_auth_integration.sh
```

---

## Production Deployment

### 1. Environment Setup

```bash
# Create production environment file
cp .env.example .env.production

# Generate strong JWT secret
openssl rand -base64 32

# Edit .env.production
nano .env.production
```

**Required Production Values**:
```env
# Strong random secret (min 32 chars)
JWT_SECRET=<generated-secret-from-above>

# Production environment
ENVIRONMENT=production
LOG_LEVEL=warn

# Your API keys
GOOGLE_MAPS_API_KEY=<your-production-api-key>

# Redis (if using external Redis)
REDIS_HOST=<your-redis-host>
REDIS_PORT=6379
REDIS_PASSWORD=<your-redis-password>
```

### 2. Security Checklist

Before deploying to production:

- [ ] **JWT Secret**: Strong random value (min 32 chars)
- [ ] **HTTPS**: Enable TLS/SSL
- [ ] **Database**: Use managed PostgreSQL (AWS RDS, etc.)
- [ ] **Redis**: Use managed Redis (AWS ElastiCache, etc.)
- [ ] **Secrets**: Use secret management (AWS Secrets Manager, Vault)
- [ ] **CORS**: Configure allowed origins
- [ ] **Rate Limiting**: Enabled and configured
- [ ] **Monitoring**: Prometheus + Grafana configured
- [ ] **Logging**: Centralized logging (ELK, CloudWatch)
- [ ] **Backups**: Database backup strategy

### 3. Build Production Images

```bash
# Build all services
docker-compose -f deployments/docker-compose/docker-compose.yml build

# Tag for registry
docker tag docker-compose-auth-service:latest your-registry/auth-service:v1.0.0
docker tag docker-compose-spider-service:latest your-registry/spider-service:v1.0.0
docker tag docker-compose-restaurant-service:latest your-registry/restaurant-service:v1.0.0
docker tag docker-compose-map-service:latest your-registry/map-service:v1.0.0

# Push to registry
docker push your-registry/auth-service:v1.0.0
docker push your-registry/spider-service:v1.0.0
docker push your-registry/restaurant-service:v1.0.0
docker push your-registry/map-service:v1.0.0
```

### 4. Deploy Services

```bash
# Using production environment file
docker-compose --env-file .env.production up -d

# Or with specific compose file
docker-compose -f docker-compose.prod.yml up -d
```

### 5. Health Checks

```bash
# Check all services
for port in 8080 18084 18082 8081; do
  echo "Checking port $port..."
  curl -f http://localhost:$port/health || echo "FAILED"
done

# Check logs
docker logs tabelogo-auth-service --tail 50
docker logs tabelogo-spider-service --tail 50
```

---

## Kubernetes Deployment

### 1. Create Secrets

```bash
# Create JWT secret
kubectl create secret generic jwt-secret \
  --from-literal=JWT_SECRET=$(openssl rand -base64 32)

# Create API keys
kubectl create secret generic api-keys \
  --from-literal=GOOGLE_MAPS_API_KEY=<your-key>
```

### 2. Deploy Services

```yaml
# auth-service-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: your-registry/auth-service:v1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: jwt-secret
              key: JWT_SECRET
        - name: REDIS_HOST
          value: redis-service
        - name: ENVIRONMENT
          value: production
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

```bash
# Apply deployments
kubectl apply -f k8s/auth-service-deployment.yaml
kubectl apply -f k8s/spider-service-deployment.yaml
kubectl apply -f k8s/restaurant-service-deployment.yaml
kubectl apply -f k8s/map-service-deployment.yaml
```

---

## Monitoring

### Prometheus

Access Prometheus at: `http://localhost:9090`

**Key Metrics to Monitor**:
- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - Request latency
- `auth_login_attempts_total` - Login attempts
- `auth_token_validations_total` - Token validations

### Grafana

Access Grafana at: `http://localhost:3000` (admin/admin)

**Dashboards**:
1. Service Overview
2. Authentication Metrics
3. API Performance
4. Error Rates

---

## Scaling

### Horizontal Scaling

```bash
# Scale specific service
docker-compose up -d --scale auth-service=3

# Kubernetes
kubectl scale deployment auth-service --replicas=5
```

### Load Balancing

Use nginx or cloud load balancer:

```nginx
upstream auth_backend {
    server auth-service-1:8080;
    server auth-service-2:8080;
    server auth-service-3:8080;
}

server {
    listen 80;
    location /api/v1/auth {
        proxy_pass http://auth_backend;
    }
}
```

---

## Backup & Recovery

### Database Backup

```bash
# PostgreSQL backup
docker exec tabelogo-postgres-auth pg_dump -U postgres auth_db > backup_auth_$(date +%Y%m%d).sql

# Restore
docker exec -i tabelogo-postgres-auth psql -U postgres auth_db < backup_auth_20231224.sql
```

### Redis Backup

```bash
# Save Redis snapshot
docker exec tabelogo-redis redis-cli SAVE

# Copy snapshot
docker cp tabelogo-redis:/data/dump.rdb ./backup_redis_$(date +%Y%m%d).rdb
```

---

## Troubleshooting

### Service Won't Start

```bash
# Check logs
docker logs tabelogo-auth-service

# Check environment variables
docker exec tabelogo-auth-service env

# Restart service
docker-compose restart auth-service
```

### Authentication Failures

```bash
# Check JWT secret consistency
docker exec tabelogo-auth-service env | grep JWT_SECRET
docker exec tabelogo-spider-service env | grep JWT_SECRET

# Check Redis connection
docker exec tabelogo-redis redis-cli ping

# Check sessions
docker exec tabelogo-redis redis-cli KEYS "session:*"
```

### Performance Issues

```bash
# Check resource usage
docker stats

# Check database connections
docker exec tabelogo-postgres-auth psql -U postgres -c "SELECT count(*) FROM pg_stat_activity;"

# Check Redis memory
docker exec tabelogo-redis redis-cli INFO memory
```

---

## Rollback

### Quick Rollback

```bash
# Stop current version
docker-compose down

# Checkout previous version
git checkout v1.0.0

# Restart services
docker-compose up -d
```

### Database Rollback

```bash
# Restore from backup
docker exec -i tabelogo-postgres-auth psql -U postgres auth_db < backup_auth_previous.sql
```

---

## Maintenance

### Update Dependencies

```bash
# Update Go modules
go get -u ./...
go mod tidy

# Rebuild images
docker-compose build

# Test before deploying
bash scripts/test_auth_integration.sh
```

### Rotate Secrets

```bash
# 1. Generate new secret
NEW_SECRET=$(openssl rand -base64 32)

# 2. Update .env
echo "JWT_SECRET=$NEW_SECRET" >> .env.new

# 3. Rolling update (zero downtime)
# Update one service at a time
docker-compose up -d --no-deps auth-service

# 4. Users will need to re-login
```

---

## Support

For issues or questions:
- Check logs: `docker logs <container-name>`
- Run tests: `bash scripts/test_auth_integration.sh`
- Review docs: `docs/CONFIGURATION.md`
