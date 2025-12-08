-- Remove area field from restaurants table
ALTER TABLE restaurants 
DROP COLUMN IF EXISTS area;

-- Drop index
DROP INDEX IF EXISTS idx_restaurants_area;
