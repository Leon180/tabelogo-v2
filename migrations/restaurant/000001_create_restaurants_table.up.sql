-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create restaurants table
CREATE TABLE IF NOT EXISTS restaurants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    source VARCHAR(50) NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    address TEXT,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    rating DECIMAL(3, 2),
    price_range VARCHAR(10),
    cuisine_type VARCHAR(50),
    phone VARCHAR(20),
    website VARCHAR(500),
    opening_hours JSONB,
    metadata JSONB,
    view_count BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Create indexes
CREATE UNIQUE INDEX idx_restaurants_source_external_id
    ON restaurants(source, external_id) WHERE deleted_at IS NULL;

CREATE INDEX idx_restaurants_name ON restaurants(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_restaurants_cuisine ON restaurants(cuisine_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_restaurants_rating ON restaurants(rating DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_restaurants_created_at ON restaurants(created_at DESC);

-- Create GiST index for location-based queries (requires postgis extension)
-- Uncomment if using PostGIS:
-- CREATE INDEX idx_restaurants_location ON restaurants USING GIST(ll_to_earth(latitude, longitude));

-- Create GIN index for JSONB columns
CREATE INDEX idx_restaurants_opening_hours ON restaurants USING GIN(opening_hours);
CREATE INDEX idx_restaurants_metadata ON restaurants USING GIN(metadata);

-- Create updated_at trigger
CREATE TRIGGER update_restaurants_updated_at BEFORE UPDATE ON restaurants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE restaurants IS 'Restaurants aggregated from multiple sources';
COMMENT ON COLUMN restaurants.source IS 'Data source: tabelog, google, opentable, etc.';
COMMENT ON COLUMN restaurants.external_id IS 'ID from external source';
COMMENT ON COLUMN restaurants.opening_hours IS 'Opening hours in JSON format';
COMMENT ON COLUMN restaurants.metadata IS 'Additional metadata from source';
