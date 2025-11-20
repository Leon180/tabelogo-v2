#!/bin/bash

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Migration 設定
MIGRATIONS_DIR="../migrations"

# 資料庫連線資訊
AUTH_DB="postgresql://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
RESTAURANT_DB="postgresql://postgres:postgres@localhost:5433/restaurant_db?sslmode=disable"
BOOKING_DB="postgresql://postgres:postgres@localhost:5434/booking_db?sslmode=disable"
SPIDER_DB="postgresql://postgres:postgres@localhost:5435/spider_db?sslmode=disable"
MAIL_DB="postgresql://postgres:postgres@localhost:5436/mail_db?sslmode=disable"

# 函數：印出訊息
print_header() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# 函數：執行 migration
run_migration() {
    local service=$1
    local db_url=$2
    local migration_path="$MIGRATIONS_DIR/$service"

    print_header "執行 $service Migration"

    # 檢查 migration 目錄是否存在
    if [ ! -d "$migration_path" ]; then
        print_error "Migration 目錄不存在: $migration_path"
        return 1
    fi

    # 執行 migration
    if migrate -path "$migration_path" -database "$db_url" up; then
        print_success "$service migration 執行成功"

        # 顯示當前版本
        version=$(migrate -path "$migration_path" -database "$db_url" version 2>&1 | tail -1)
        print_info "當前版本: $version"
    else
        print_error "$service migration 執行失敗"
        return 1
    fi

    echo ""
}

# 函數：顯示 migration 版本
show_version() {
    local service=$1
    local db_url=$2
    local migration_path="$MIGRATIONS_DIR/$service"

    version=$(migrate -path "$migration_path" -database "$db_url" version 2>&1 | tail -1)
    echo -e "${service}: ${GREEN}${version}${NC}"
}

# 主程式
main() {
    print_header "🚀 Tabelogo Microservices Migration Tool"

    # 檢查 migrate 是否安裝
    if ! command -v migrate &> /dev/null; then
        print_error "golang-migrate 未安裝"
        print_info "安裝方式："
        echo "  macOS: brew install golang-migrate"
        echo "  Linux: https://github.com/golang-migrate/migrate/releases"
        exit 1
    fi

    print_success "golang-migrate 已安裝"
    echo ""

    # 檢查 Docker 容器是否運行
    print_header "檢查資料庫容器狀態"

    containers=("tabelogo-auth-db" "tabelogo-restaurant-db" "tabelogo-booking-db" "tabelogo-spider-db" "tabelogo-mail-db")
    all_running=true

    for container in "${containers[@]}"; do
        if docker ps --format '{{.Names}}' | grep -q "^${container}$"; then
            print_success "$container 運行中"
        else
            print_error "$container 未運行"
            all_running=false
        fi
    done

    if [ "$all_running" = false ]; then
        print_info "請先啟動資料庫容器："
        echo "  cd deployments/docker-compose"
        echo "  docker-compose up -d"
        exit 1
    fi

    echo ""
    sleep 2

    # 執行所有 migrations
    print_header "開始執行 Migrations"

    run_migration "auth" "$AUTH_DB"
    run_migration "restaurant" "$RESTAURANT_DB"
    run_migration "booking" "$BOOKING_DB"
    run_migration "spider" "$SPIDER_DB"
    run_migration "mail" "$MAIL_DB"

    # 顯示所有版本
    print_header "📊 所有服務的 Migration 版本"
    show_version "auth" "$AUTH_DB"
    show_version "restaurant" "$RESTAURANT_DB"
    show_version "booking" "$BOOKING_DB"
    show_version "spider" "$SPIDER_DB"
    show_version "mail" "$MAIL_DB"
    echo ""

    print_success "所有 Migrations 執行完成！"
}

# 執行主程式
main
