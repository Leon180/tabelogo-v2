#!/bin/sh
set -e

echo "Starting Spider Service..."

# Wait for dependencies
if [ -n "$REDIS_HOST" ]; then
    echo "Waiting for Redis at $REDIS_HOST:${REDIS_PORT:-6379}..."
    timeout 30 sh -c 'until nc -z $REDIS_HOST ${REDIS_PORT:-6379}; do sleep 1; done' || echo "Redis not available, continuing anyway..."
fi

# Start the service
exec ./spider-service
