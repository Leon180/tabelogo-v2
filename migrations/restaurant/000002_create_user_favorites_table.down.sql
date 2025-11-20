-- Drop user_favorites table
DROP TRIGGER IF EXISTS update_user_favorites_updated_at ON user_favorites;
DROP INDEX IF EXISTS idx_user_favorites_tags;
DROP INDEX IF EXISTS idx_user_favorites_unique;
DROP INDEX IF EXISTS idx_user_favorites_last_visited;
DROP INDEX IF EXISTS idx_user_favorites_created_at;
DROP INDEX IF EXISTS idx_user_favorites_restaurant_id;
DROP INDEX IF EXISTS idx_user_favorites_user_id;
DROP TABLE IF EXISTS user_favorites;
