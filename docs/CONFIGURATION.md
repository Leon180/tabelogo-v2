# Configuration Guide

## Environment Variables

### Quick Start

1. **Copy the example file**:
   ```bash
   cp .env.example .env
   ```

2. **Update required variables**:
   - `JWT_SECRET`: Change to a strong secret (minimum 32 characters)
   - `GOOGLE_MAPS_API_KEY`: Your Google Maps API key

3. **Start services**:
   ```bash
   make up
   ```

---

## Configuration Variables

### JWT Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `JWT_SECRET` | Secret key for JWT signing | - | ✅ Yes |
| `JWT_ACCESS_TOKEN_EXPIRE` | Access token lifetime | `15m` | No |
| `JWT_REFRESH_TOKEN_EXPIRE` | Refresh token lifetime | `168h` | No |

> **⚠️ IMPORTANT**: `JWT_SECRET` must be at least 32 characters long in production.

### Redis Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `REDIS_HOST` | Redis server host | `redis` | Yes |
| `REDIS_PORT` | Redis server port | `6379` | Yes |
| `REDIS_PASSWORD` | Redis password | `` | No |
| `AUTH_REDIS_DB` | Redis DB for sessions | `0` | Yes |

> **📝 NOTE**: All services share `AUTH_REDIS_DB=0` for session storage.

### External APIs

| Variable | Description | Required |
|----------|-------------|----------|
| `GOOGLE_MAPS_API_KEY` | Google Maps API key | ✅ Yes |

---

## Security Best Practices

### 1. Never Commit `.env` Files

The `.env` file is gitignored by default. **Never** commit it to version control.

```bash
# ✅ Good - .env is gitignored
git status
# .env should NOT appear

# ❌ Bad - forcing .env into git
git add -f .env  # DON'T DO THIS!
```

### 2. Use Strong JWT Secrets

Generate a strong random secret:

```bash
# Option 1: Using openssl
openssl rand -base64 32

# Option 2: Using /dev/urandom
head -c 32 /dev/urandom | base64

# Option 3: Using Python
python3 -c "import secrets; print(secrets.token_urlsafe(32))"
```

### 3. Rotate Secrets Regularly

In production, rotate `JWT_SECRET` periodically:
- Recommended: Every 90 days
- Update `.env` file
- Restart all services
- Users will need to re-login

### 4. Use Secret Management Tools

For production deployments, consider:
- **AWS Secrets Manager**
- **HashiCorp Vault**
- **Kubernetes Secrets**
- **Azure Key Vault**

---

## Environment-Specific Configuration

### Development

```env
# .env (development)
JWT_SECRET=dev-secret-change-in-production-min-32-chars
ENVIRONMENT=development
LOG_LEVEL=debug
```

### Staging

```env
# .env.staging
JWT_SECRET=<strong-random-secret-32-chars>
ENVIRONMENT=staging
LOG_LEVEL=info
```

### Production

```env
# .env.production
JWT_SECRET=<strong-random-secret-from-secret-manager>
ENVIRONMENT=production
LOG_LEVEL=warn
```

---

## Docker Compose Integration

### How Environment Variables Are Loaded

1. **System environment variables** (highest priority)
2. **`.env` file** in project root
3. **Fallback values** in `docker-compose.yml` (lowest priority)

### Example

```yaml
# docker-compose.yml
services:
  auth-service:
    environment:
      # Uses .env value, falls back to default
      JWT_SECRET: ${JWT_SECRET:-change-me-in-production}
```

### Override for Specific Environments

```bash
# Use staging configuration
docker-compose --env-file .env.staging up -d

# Use production configuration
docker-compose --env-file .env.production up -d
```

---

## Troubleshooting

### Environment Variables Not Loading

**Problem**: Services use default values instead of `.env` values.

**Solution**:
```bash
# 1. Check .env file exists
ls -la .env

# 2. Verify .env content
cat .env | grep JWT_SECRET

# 3. Restart services with force recreate
docker-compose down
docker-compose up -d --force-recreate

# 4. Verify loaded values
docker exec tabelogo-auth-service env | grep JWT_SECRET
```

### JWT Signature Verification Failed

**Problem**: "Invalid or expired token" errors.

**Possible Causes**:
1. **Different JWT_SECRET across services**
   ```bash
   # Check all services have same secret
   docker exec tabelogo-auth-service env | grep JWT_SECRET
   docker exec tabelogo-spider-service env | grep JWT_SECRET
   ```

2. **JWT_SECRET changed after token issued**
   - Users need to re-login
   - Clear Redis sessions: `docker exec tabelogo-redis redis-cli FLUSHDB`

### Redis Connection Issues

**Problem**: "Session expired or revoked" errors.

**Solution**:
```bash
# 1. Check Redis is running
docker ps | grep redis

# 2. Test Redis connection
docker exec tabelogo-redis redis-cli ping
# Should return: PONG

# 3. Check session exists
docker exec tabelogo-redis redis-cli KEYS "session:*"

# 4. Verify AUTH_REDIS_DB is 0
docker exec tabelogo-auth-service env | grep AUTH_REDIS_DB
```

---

## Configuration Checklist

Before deploying:

- [ ] `.env` file created from `.env.example`
- [ ] `JWT_SECRET` is strong (min 32 chars)
- [ ] `GOOGLE_MAPS_API_KEY` is set
- [ ] `.env` is gitignored
- [ ] All services use same `JWT_SECRET`
- [ ] `AUTH_REDIS_DB=0` for all services
- [ ] Secrets are stored securely (production)
- [ ] Environment-specific configs prepared

---

## Additional Resources

- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [Redis Security](https://redis.io/topics/security)
- [Docker Secrets](https://docs.docker.com/engine/swarm/secrets/)
