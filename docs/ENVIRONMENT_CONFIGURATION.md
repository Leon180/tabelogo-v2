# Environment Configuration Guide

## Overview

Sensitive environment variables (API keys, secrets, passwords) are stored in `.env` files that are **gitignored** and never committed to version control.

## Setup Instructions

### For Docker Compose

1. **Copy the template**:
   ```bash
   cd deployments/docker-compose
   cp .env.example .env
   ```

2. **Edit `.env` file** with your actual values:
   ```bash
   # Required
   GOOGLE_MAPS_API_KEY=your_actual_google_maps_api_key
   
   # Optional (defaults are provided in docker-compose.yml)
   POSTGRES_PASSWORD=your_secure_password
   JWT_SECRET=your_jwt_secret_at_least_32_chars
   ```

3. **Start services**:
   ```bash
   docker-compose up -d
   ```

Docker Compose automatically loads variables from `.env` file in the same directory.

### For Local Development

Each service can have its own `.env` file in `cmd/<service-name>/.env`:

```bash
# Example: cmd/map-service/.env
GOOGLE_MAPS_API_KEY=your_api_key_here
REDIS_HOST=localhost
REDIS_PORT=6379
```

## Environment Variable Priority

1. **System environment variables** (highest priority)
2. **`.env` file** in the same directory as docker-compose.yml
3. **Default values** in docker-compose.yml (lowest priority)

## Security Best Practices

### ✅ DO:
- Keep `.env` files gitignored
- Use `.env.example` as a template (committed to git)
- Document required variables in `.env.example`
- Use different values for dev/staging/production
- Rotate API keys and secrets regularly

### ❌ DON'T:
- Commit `.env` files to git
- Hardcode secrets in docker-compose.yml
- Share `.env` files via email/chat
- Use production secrets in development

## Required Environment Variables

### Map Service
- `GOOGLE_MAPS_API_KEY` - **Required** - Google Maps Platform API key

### Auth Service
- `JWT_SECRET` - Optional (default provided, change in production)
- `POSTGRES_PASSWORD` - Optional (default: postgres)

### Restaurant Service
- `POSTGRES_PASSWORD` - Optional (default: postgres)
- `JWT_SECRET` - Optional (default provided, change in production)

## Verification

Check if environment variables are loaded correctly:

```bash
# Check Map Service
docker exec tabelogo-map-service env | grep GOOGLE_MAPS_API_KEY

# Should show: GOOGLE_MAPS_API_KEY=AIzaSy...
# Should NOT show: GOOGLE_MAPS_API_KEY=your_api_key_here
```

## Troubleshooting

### Variables not loading

1. **Check .env file location**: Must be in `deployments/docker-compose/.env`
2. **Restart services**: `docker-compose down && docker-compose up -d`
3. **Check syntax**: No spaces around `=`, no quotes needed

### API key invalid

1. **Verify key in Google Cloud Console**
2. **Check API is enabled**: Places API (New)
3. **Check billing**: API requires billing enabled
4. **Check restrictions**: Remove HTTP referrer restrictions for testing

## Example .env File

```bash
# Tabelogo v2 - Environment Configuration
GOOGLE_MAPS_API_KEY=your_actual_google_maps_api_key_here
```

## Production Deployment

For production, use:
- **AWS Secrets Manager**
- **HashiCorp Vault**
- **Kubernetes Secrets**
- **Environment variables** in CI/CD platform

Never use `.env` files in production deployments.
