#!/bin/sh
set -e

echo "==> Starting Restaurant Service..."

# Wait for PostgreSQL to be ready
echo "==> Waiting for PostgreSQL (${DB_NAME:-restaurant_db} at ${DB_HOST:-localhost}:${DB_PORT:-5432})..."
until nc -z ${DB_HOST:-localhost} ${DB_PORT:-5432}; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done
echo "==> PostgreSQL is up"

# Wait for Redis to be ready (optional)
if [ -n "${REDIS_HOST}" ]; then
  echo "==> Waiting for Redis..."
  until nc -z ${REDIS_HOST} ${REDIS_PORT:-6379}; do
    echo "Redis is unavailable - sleeping"
    sleep 1
  done
  echo "==> Redis is up"
fi

# Run migrations
echo "==> Running database migrations..."
if [ -f "/app/migrate" ]; then
  /app/migrate -path /app/migrations/restaurant -database "postgresql://${DB_USER:-postgres}:${DB_PASSWORD:-postgres}@${DB_HOST:-localhost}:${DB_PORT:-5432}/${DB_NAME:-restaurant_db}?sslmode=${DB_SSLMODE:-disable}" up
  echo "==> Migrations completed successfully"
else
  echo "WARNING: migrate binary not found, skipping migrations"
fi

# Start the application
echo "==> Starting Restaurant Service on port ${SERVER_PORT:-18082}..."
exec /app/restaurant-service
