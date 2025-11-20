-- Drop trigger
DROP TRIGGER IF EXISTS update_restaurants_updated_at ON restaurants;

-- Drop indexes
DROP INDEX IF EXISTS idx_restaurants_metadata;
DROP INDEX IF EXISTS idx_restaurants_opening_hours;
-- DROP INDEX IF EXISTS idx_restaurants_location; -- Uncomment if using PostGIS
DROP INDEX IF EXISTS idx_restaurants_created_at;
DROP INDEX IF EXISTS idx_restaurants_rating;
DROP INDEX IF EXISTS idx_restaurants_cuisine;
DROP INDEX IF EXISTS idx_restaurants_name;
DROP INDEX IF EXISTS idx_restaurants_source_external_id;

-- Drop restaurants table
DROP TABLE IF EXISTS restaurants;
