-- Remove name_ja column from restaurants table
DROP INDEX IF EXISTS idx_restaurants_name_ja;
ALTER TABLE restaurants DROP COLUMN IF EXISTS name_ja;
