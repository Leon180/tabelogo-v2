-- Create user_favorites table
-- Users can favorite/bookmark restaurants for easy access
CREATE TABLE IF NOT EXISTS user_favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,  -- Reference to auth_db.users (no FK due to microservices)
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    notes TEXT,  -- User's private notes about this restaurant
    tags VARCHAR(255)[],  -- User-defined tags for organization
    visit_count INT DEFAULT 0,  -- How many times user marked as visited
    last_visited_at TIMESTAMP,  -- Last time user marked as visited
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL  -- Soft delete for unfavorite
);

-- Create indexes
CREATE INDEX idx_user_favorites_user_id ON user_favorites(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_favorites_restaurant_id ON user_favorites(restaurant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_favorites_created_at ON user_favorites(created_at DESC);
CREATE INDEX idx_user_favorites_last_visited ON user_favorites(last_visited_at DESC);

-- Unique constraint to prevent duplicate favorites
CREATE UNIQUE INDEX idx_user_favorites_unique ON user_favorites(user_id, restaurant_id) WHERE deleted_at IS NULL;

-- Create GIN index for tags array
CREATE INDEX idx_user_favorites_tags ON user_favorites USING GIN(tags);

-- Create updated_at trigger
CREATE TRIGGER update_user_favorites_updated_at BEFORE UPDATE ON user_favorites
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE user_favorites IS 'User favorite/bookmarked restaurants';
COMMENT ON COLUMN user_favorites.user_id IS 'Reference to users.id from auth service (no FK)';
COMMENT ON COLUMN user_favorites.restaurant_id IS 'Reference to restaurants.id';
COMMENT ON COLUMN user_favorites.notes IS 'User private notes about this restaurant';
COMMENT ON COLUMN user_favorites.tags IS 'User-defined tags for organization (e.g., want_to_visit, date_night)';
COMMENT ON COLUMN user_favorites.visit_count IS 'Number of times user marked as visited';
COMMENT ON COLUMN user_favorites.deleted_at IS 'Soft delete timestamp (unfavorite)';
