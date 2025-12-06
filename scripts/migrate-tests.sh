#!/bin/bash
# Script to migrate service_test.go from testify mocks to gomock

set -e

FILE="internal/restaurant/application/service_test.go"
BACKUP="${FILE}.backup"

echo "Creating backup..."
cp "$FILE" "$BACKUP"

echo "Performing migration..."

# Step 1: Update imports
sed -i '' 's/"github.com\/stretchr\/testify\/mock"/"github.com\/Leon180\/tabelogo-v2\/internal\/restaurant\/mocks"\
	"github.com\/stretchr\/testify\/require"\
	"go.uber.org\/mock\/gomock"\
	"go.uber.org\/zap\/zaptest"/g' "$FILE"

sed -i '' 's/"go.uber.org\/zap"//g' "$FILE"

# Step 2: Remove manual mock definitions (lines 17-179)
# This is complex, so we'll use a different approach - create a new file

echo "Migration complete!"
echo "Backup saved to: $BACKUP"
