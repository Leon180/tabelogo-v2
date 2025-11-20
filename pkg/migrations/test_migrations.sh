#!/bin/bash
# Migration 功能測試腳本

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 測試配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

DB_HOST="localhost"
DB_PORT="5433"
DB_USER="testuser"
DB_PASS="testpass"
DB_NAME="testdb"
DSN="postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
MIGRATIONS_PATH="$PROJECT_ROOT/migrations/auth"
SERVICE_NAME="auth"
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.test.yml"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Migration 功能測試${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 檢查 Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker 未安裝${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose 未安裝${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Docker 環境檢查通過${NC}"
echo ""

# 啟動測試資料庫
echo -e "${BLUE}📦 啟動 PostgreSQL 測試容器...${NC}"
docker-compose -f "$COMPOSE_FILE" up -d

# 等待資料庫就緒
echo -e "${YELLOW}⏳ 等待資料庫就緒...${NC}"
MAX_TRIES=30
TRIES=0
while [ $TRIES -lt $MAX_TRIES ]; do
    if docker-compose -f "$COMPOSE_FILE" exec -T postgres pg_isready -U $DB_USER -d $DB_NAME &> /dev/null; then
        echo -e "${GREEN}✅ 資料庫已就緒${NC}"
        break
    fi
    TRIES=$((TRIES + 1))
    echo -n "."
    sleep 1
done

if [ $TRIES -eq $MAX_TRIES ]; then
    echo -e "${RED}❌ 資料庫啟動超時${NC}"
    docker-compose -f "$COMPOSE_FILE" logs
    docker-compose -f "$COMPOSE_FILE" down
    exit 1
fi
echo ""

# 測試函數
test_passed=0
test_failed=0

run_test() {
    local test_name=$1
    local command=$2

    echo -e "${BLUE}🧪 測試: ${test_name}${NC}"
    if eval "$command"; then
        echo -e "${GREEN}✅ 通過${NC}"
        test_passed=$((test_passed + 1))
    else
        echo -e "${RED}❌ 失敗${NC}"
        test_failed=$((test_failed + 1))
    fi
    echo ""
}

# 建置 CLI 工具
echo -e "${BLUE}🔨 建置 Migration CLI 工具...${NC}"
mkdir -p "$SCRIPT_DIR/bin"
cd "$PROJECT_ROOT/pkg"
go build -o "$SCRIPT_DIR/bin/migrate" "$SCRIPT_DIR/cmd/migrate/main.go"
cd "$SCRIPT_DIR"
echo -e "${GREEN}✅ CLI 工具建置完成${NC}"
echo ""

# 測試 1: 檢查版本（應該沒有 migration）
run_test "檢查初始版本" "
    ./bin/migrate -dsn '$DSN' -path '$MIGRATIONS_PATH' -service '$SERVICE_NAME' -command version 2>&1 | grep -q 'no migration'
"

# 測試 2: 執行 Up（執行所有 migrations）
run_test "執行 Up migration" "
    ./bin/migrate -dsn '$DSN' -path '$MIGRATIONS_PATH' -service '$SERVICE_NAME' -command up
"

# 測試 3: 檢查版本（應該是 2）
run_test "檢查版本是否為 2" "
    ./bin/migrate -dsn '$DSN' -path '$MIGRATIONS_PATH' -service '$SERVICE_NAME' -command version 2>&1 | grep -q 'Current version: 2'
"

# 測試 4: 驗證表格是否建立
run_test "驗證 users 表格存在" "
    docker-compose -f "$COMPOSE_FILE" exec -T postgres psql -U $DB_USER -d $DB_NAME -c '\dt users' | grep -q 'users'
"

run_test "驗證 refresh_tokens 表格存在" "
    docker-compose -f "$COMPOSE_FILE" exec -T postgres psql -U $DB_USER -d $DB_NAME -c '\dt refresh_tokens' | grep -q 'refresh_tokens'
"

# 測試 5: 檢查版本控制表
run_test "檢查版本控制表存在" "
    docker-compose -f "$COMPOSE_FILE" exec -T postgres psql -U $DB_USER -d $DB_NAME -c 'SELECT * FROM schema_migrations_auth' | grep -q '2'
"

# 測試 6: 執行 Down（回滾一個 migration）
run_test "執行 Down migration" "
    ./bin/migrate -dsn '$DSN' -path '$MIGRATIONS_PATH' -service '$SERVICE_NAME' -command down
"

# 測試 7: 檢查版本（應該是 1）
run_test "檢查版本是否回到 1" "
    ./bin/migrate -dsn '$DSN' -path '$MIGRATIONS_PATH' -service '$SERVICE_NAME' -command version 2>&1 | grep -q 'Current version: 1'
"

# 測試 8: 驗證 refresh_tokens 表格已刪除
run_test "驗證 refresh_tokens 表格已刪除" "
    ! docker-compose -f "$COMPOSE_FILE" exec -T postgres psql -U $DB_USER -d $DB_NAME -c '\dt refresh_tokens' | grep -q 'refresh_tokens'
"

# 測試 9: 執行 Steps（向上 1 步）
run_test "執行 Steps +1" "
    ./bin/migrate -dsn '$DSN' -path '$MIGRATIONS_PATH' -service '$SERVICE_NAME' -command steps -steps 1
"

# 測試 10: 檢查版本（應該又是 2）
run_test "檢查版本是否回到 2" "
    ./bin/migrate -dsn '$DSN' -path '$MIGRATIONS_PATH' -service '$SERVICE_NAME' -command version 2>&1 | grep -q 'Current version: 2'
"

# 測試 11: 驗證狀態（不應該是 dirty）
run_test "驗證狀態不是 dirty" "
    ./bin/migrate -dsn '$DSN' -path '$MIGRATIONS_PATH' -service '$SERVICE_NAME' -command validate
"

# 測試 12: 檢查索引是否建立
run_test "檢查 users 表的索引" "
    docker-compose -f "$COMPOSE_FILE" exec -T postgres psql -U $DB_USER -d $DB_NAME -c '\di' | grep -q 'idx_users_email'
"

# 測試 13: 檢查觸發器是否建立
run_test "檢查 updated_at 觸發器" "
    docker-compose -f "$COMPOSE_FILE" exec -T postgres psql -U $DB_USER -d $DB_NAME -c 'SELECT tgname FROM pg_trigger WHERE tgname LIKE '\''%users%'\''' | grep -q 'update_users_updated_at'
"

# 顯示資料庫狀態
echo -e "${BLUE}📊 資料庫狀態:${NC}"
echo ""
echo "版本控制表內容:"
docker-compose -f "$COMPOSE_FILE" exec -T postgres psql -U $DB_USER -d $DB_NAME -c "SELECT * FROM schema_migrations_auth;"
echo ""

echo "所有表格:"
docker-compose -f "$COMPOSE_FILE" exec -T postgres psql -U $DB_USER -d $DB_NAME -c "\dt"
echo ""

# 測試結果統計
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  測試結果統計${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✅ 通過: ${test_passed}${NC}"
echo -e "${RED}❌ 失敗: ${test_failed}${NC}"
echo -e "${BLUE}總計: $((test_passed + test_failed))${NC}"
echo ""

# 清理
echo -e "${YELLOW}🧹 清理測試環境...${NC}"
read -p "是否清理 Docker 容器? (y/n) " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker-compose -f "$COMPOSE_FILE" down -v
    echo -e "${GREEN}✅ 清理完成${NC}"
else
    echo -e "${YELLOW}⚠️  容器仍在執行，可手動清理: docker-compose -f "$COMPOSE_FILE" down -v${NC}"
fi

# 返回結果
if [ $test_failed -eq 0 ]; then
    echo -e "${GREEN}🎉 所有測試通過！${NC}"
    exit 0
else
    echo -e "${RED}❌ 有 ${test_failed} 個測試失敗${NC}"
    exit 1
fi
