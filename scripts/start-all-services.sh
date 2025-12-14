#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=================================================="
echo "Starting Tabelogo v2 Microservices"
echo -e "==================================================${NC}"
echo ""

# Navigate to docker-compose directory
cd "$(dirname "$0")/../deployments/docker-compose"

# Function to print colored messages
print_error() {
    echo -e "${RED}ERROR: $1${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}→ $1${NC}"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker Desktop and try again."
    exit 1
fi

print_success "Docker is running"
echo ""

# Show current docker-compose services
print_info "Services available in docker-compose.yml:"
echo ""
echo "  Infrastructure:"
echo "    • postgres-auth (port 5432)"
echo "    • postgres-restaurant (port 5433)"
echo "    • postgres-booking (port 5434)"
echo "    • redis (port 6379)"
echo "    • kafka (port 9092)"
echo "    • zookeeper (port 2181)"
echo ""
echo "  Microservices:"
echo "    • auth-service (HTTP: 8080, gRPC: 9090) ✓ Ready"
echo "    • restaurant-service (HTTP: 18082, gRPC: 19082) ✓ Ready"
echo "    • booking-service (commented out - not implemented)"
echo "    • api-gateway (commented out - not implemented)"
echo ""
echo "  Monitoring:"
echo "    • prometheus (port 9090)"
echo "    • grafana (port 3000)"
echo ""

# Ask user what to start
echo -e "${YELLOW}What would you like to start?${NC}"
echo "  1) Infrastructure only (databases, redis, kafka)"
echo "  2) Infrastructure + Auth Service"
echo "  3) Infrastructure + Restaurant Service"
echo "  4) Infrastructure + Both Services (auth + restaurant)"
echo "  5) Everything (infrastructure + services + monitoring)"
echo "  6) Custom selection"
echo ""
read -p "Enter your choice (1-6): " choice

case $choice in
    1)
        print_info "Starting infrastructure services..."
        SERVICES="postgres-auth postgres-restaurant postgres-booking redis kafka zookeeper"
        ;;
    2)
        print_info "Starting infrastructure + auth service..."
        SERVICES="postgres-auth redis kafka zookeeper auth-service"
        ;;
    3)
        print_info "Starting infrastructure + restaurant service..."
        SERVICES="postgres-restaurant redis kafka zookeeper restaurant-service"
        ;;
    4)
        print_info "Starting infrastructure + both services..."
        SERVICES="postgres-auth postgres-restaurant redis kafka zookeeper auth-service restaurant-service"
        ;;
    5)
        print_info "Starting everything..."
        SERVICES="" # Empty means all services
        ;;
    6)
        print_info "Available services: postgres-auth postgres-restaurant postgres-booking redis kafka zookeeper auth-service restaurant-service prometheus grafana"
        read -p "Enter services to start (space-separated): " SERVICES
        ;;
    *)
        print_error "Invalid choice"
        exit 1
        ;;
esac

echo ""
print_info "Building services (if needed)..."
echo ""

# Build services using legacy builder (to avoid BuildKit corruption)
if [ -z "$SERVICES" ]; then
    # Build all
    DOCKER_BUILDKIT=0 docker compose build
else
    # Build only specified services that need building
    for service in $SERVICES; do
        if [ "$service" = "auth-service" ] || [ "$service" = "restaurant-service" ]; then
            print_info "Building $service..."
            DOCKER_BUILDKIT=0 docker compose build $service
        fi
    done
fi

echo ""
print_info "Starting services..."
echo ""

# Start services using legacy builder
if [ -z "$SERVICES" ]; then
    DOCKER_BUILDKIT=0 docker compose up -d
else
    DOCKER_BUILDKIT=0 docker compose up -d $SERVICES
fi

echo ""
print_info "Waiting for services to be healthy (30 seconds)..."
sleep 30

echo ""
print_success "Service Status:"
echo ""
docker compose ps

echo ""
print_success "=================================================="
print_success "Services Started Successfully!"
print_success "=================================================="
echo ""

# Show endpoints
echo -e "${BLUE}Available Endpoints:${NC}"
echo ""

# Check which services are running
RUNNING_SERVICES=$(docker compose ps --services --filter "status=running")

if echo "$RUNNING_SERVICES" | grep -q "auth-service"; then
    echo "  Auth Service:"
    echo "    • HTTP API:  http://localhost:8080"
    echo "    • gRPC:      localhost:9090"
    echo "    • Health:    http://localhost:8080/health"
    echo "    • Swagger:   http://localhost:8080/swagger/index.html"
    echo ""
fi

if echo "$RUNNING_SERVICES" | grep -q "restaurant-service"; then
    echo "  Restaurant Service:"
    echo "    • HTTP API:  http://localhost:18082"
    echo "    • gRPC:      localhost:19082"
    echo "    • Health:    http://localhost:18082/health"
    echo "    • Swagger:   http://localhost:18082/swagger/index.html"
    echo ""
fi

if echo "$RUNNING_SERVICES" | grep -q "grafana"; then
    echo "  Monitoring:"
    echo "    • Grafana:   http://localhost:3000 (admin/admin)"
    echo "    • Prometheus: http://localhost:9090"
    echo ""
fi

echo -e "${BLUE}Useful Commands:${NC}"
echo ""
echo "  # View logs"
echo "  docker compose logs -f [service-name]"
echo ""
echo "  # View all logs"
echo "  docker compose logs -f"
echo ""
echo "  # Stop all services"
echo "  docker compose down"
echo ""
echo "  # Stop and remove volumes (WARNING: deletes data)"
echo "  docker compose down -v"
echo ""
echo "  # Restart a service"
echo "  docker compose restart [service-name]"
echo ""
echo "  # Check service status"
echo "  docker compose ps"
echo ""

# Show recent logs from services
if [ "$choice" != "1" ]; then
    echo -e "${BLUE}Recent logs (last 20 lines per service):${NC}"
    echo ""
    if echo "$RUNNING_SERVICES" | grep -q "auth-service"; then
        echo "--- Auth Service ---"
        docker compose logs --tail=20 auth-service
        echo ""
    fi
    if echo "$RUNNING_SERVICES" | grep -q "restaurant-service"; then
        echo "--- Restaurant Service ---"
        docker compose logs --tail=20 restaurant-service
        echo ""
    fi
fi

print_success "Done! Services are running."
