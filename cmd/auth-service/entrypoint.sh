#!/bin/sh
set -e

echo "==> Starting Auth Service..."

# Wait for PostgreSQL to be ready
echo "==> Waiting for PostgreSQL..."
until nc -z ${DB_HOST:-localhost} ${DB_PORT:-5432}; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done
echo "==> PostgreSQL is up"

# Run migrations
echo "==> Running database migrations..."
if [ -f "/app/migrate" ]; then
  /app/migrate -path /app/migrations/auth -database "postgresql://${DB_USER:-postgres}:${DB_PASSWORD:-postgres}@${DB_HOST:-localhost}:${DB_PORT:-5432}/${DB_NAME:-auth_db}?sslmode=${DB_SSLMODE:-disable}" up
  echo "==> Migrations completed successfully"
else
  echo "WARNING: migrate binary not found, skipping migrations"
fi

# Start the application
echo "==> Starting application..."
exec /app/auth-service
