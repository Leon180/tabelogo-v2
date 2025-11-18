#!/bin/bash
# Migration 健康檢查腳本

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 使用說明
usage() {
    echo "Usage: $0 <service_name> <migrations_path>"
    echo ""
    echo "Example:"
    echo "  $0 auth migrations/auth"
    echo ""
    exit 1
}

# 檢查參數
if [ $# -ne 2 ]; then
    usage
fi

SERVICE=$1
MIGRATIONS_PATH=$2

echo "======================================"
echo "Migration Health Check: $SERVICE"
echo "======================================"
echo ""

# 檢查目錄是否存在
if [ ! -d "$MIGRATIONS_PATH" ]; then
    echo -e "${RED}❌ Directory not found: $MIGRATIONS_PATH${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Directory exists: $MIGRATIONS_PATH${NC}"
echo ""

# 列出所有 migration 檔案
echo "📋 Available migration files:"
echo "------------------------------"
ls -1 "$MIGRATIONS_PATH" | sort
echo ""

# 檢查配對
echo "🔍 Checking up/down pairs:"
echo "------------------------------"

up_files=$(ls "$MIGRATIONS_PATH"/*.up.sql 2>/dev/null || true)
down_files=$(ls "$MIGRATIONS_PATH"/*.down.sql 2>/dev/null || true)

if [ -z "$up_files" ]; then
    echo -e "${YELLOW}⚠️  No up migration files found${NC}"
    exit 0
fi

missing_down=0
for up in $up_files; do
    version=$(basename "$up" | sed 's/\([0-9]*\)_.*/\1/')
    down="${up%.up.sql}.down.sql"

    if [ -f "$down" ]; then
        echo -e "${GREEN}✅ $version - has both up and down${NC}"
    else
        echo -e "${RED}❌ $version - missing down migration${NC}"
        missing_down=$((missing_down + 1))
    fi
done
echo ""

# 檢查版本號是否重複
echo "🔢 Checking for duplicate versions:"
echo "------------------------------"
versions=$(ls "$MIGRATIONS_PATH"/*.up.sql 2>/dev/null | sed 's/.*\/\([0-9]*\)_.*/\1/' | sort)
duplicates=$(echo "$versions" | uniq -d)

if [ -z "$duplicates" ]; then
    echo -e "${GREEN}✅ No duplicate versions found${NC}"
else
    echo -e "${RED}❌ Duplicate versions found:${NC}"
    echo "$duplicates"
fi
echo ""

# 檢查檔案命名格式
echo "📝 Checking file naming format:"
echo "------------------------------"
invalid_names=0
for file in "$MIGRATIONS_PATH"/*.sql; do
    basename=$(basename "$file")
    # 檢查格式: {version}_{description}.{up|down}.sql
    if [[ ! "$basename" =~ ^[0-9]+_[a-z0-9_]+\.(up|down)\.sql$ ]]; then
        echo -e "${RED}❌ Invalid format: $basename${NC}"
        echo "   Expected: {version}_{description}.{up|down}.sql"
        invalid_names=$((invalid_names + 1))
    fi
done

if [ $invalid_names -eq 0 ]; then
    echo -e "${GREEN}✅ All filenames are correctly formatted${NC}"
fi
echo ""

# 檢查 SQL 語法 (基本檢查)
echo "🔍 Checking SQL syntax (basic):"
echo "------------------------------"
syntax_errors=0
for file in "$MIGRATIONS_PATH"/*.sql; do
    basename=$(basename "$file")

    # 檢查是否使用交易
    if ! grep -qi "BEGIN\|START TRANSACTION" "$file"; then
        echo -e "${YELLOW}⚠️  $basename - no BEGIN transaction found${NC}"
    fi

    # 檢查是否有 COMMIT
    if grep -qi "BEGIN\|START TRANSACTION" "$file" && ! grep -qi "COMMIT" "$file"; then
        echo -e "${RED}❌ $basename - has BEGIN but no COMMIT${NC}"
        syntax_errors=$((syntax_errors + 1))
    fi

    # 檢查 up migration 是否使用 IF NOT EXISTS
    if [[ "$basename" == *.up.sql ]] && grep -qi "CREATE TABLE" "$file" && ! grep -qi "IF NOT EXISTS" "$file"; then
        echo -e "${YELLOW}⚠️  $basename - CREATE TABLE without IF NOT EXISTS${NC}"
    fi

    # 檢查 down migration 是否使用 IF EXISTS
    if [[ "$basename" == *.down.sql ]] && grep -qi "DROP" "$file" && ! grep -qi "IF EXISTS" "$file"; then
        echo -e "${YELLOW}⚠️  $basename - DROP without IF EXISTS${NC}"
    fi
done

if [ $syntax_errors -eq 0 ]; then
    echo -e "${GREEN}✅ No critical syntax issues found${NC}"
fi
echo ""

# 統計資訊
echo "📊 Statistics:"
echo "------------------------------"
up_count=$(ls "$MIGRATIONS_PATH"/*.up.sql 2>/dev/null | wc -l)
down_count=$(ls "$MIGRATIONS_PATH"/*.down.sql 2>/dev/null | wc -l)
echo "Up migrations:   $up_count"
echo "Down migrations: $down_count"
echo ""

# 最終結果
echo "======================================"
echo "Summary"
echo "======================================"

total_issues=$((missing_down + syntax_errors))
if [ $total_issues -eq 0 ] && [ -z "$duplicates" ] && [ $invalid_names -eq 0 ]; then
    echo -e "${GREEN}✅ All checks passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Found $total_issues issue(s)${NC}"
    echo ""
    echo "Issues:"
    [ $missing_down -gt 0 ] && echo "  - $missing_down missing down migration(s)"
    [ -n "$duplicates" ] && echo "  - Duplicate version numbers found"
    [ $invalid_names -gt 0 ] && echo "  - $invalid_names invalid filename(s)"
    [ $syntax_errors -gt 0 ] && echo "  - $syntax_errors syntax error(s)"
    exit 1
fi
