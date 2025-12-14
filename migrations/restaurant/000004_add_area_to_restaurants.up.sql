-- Add area field to restaurants table
-- Area stores the administrative_area_level_1 from Google Maps (e.g., "Tokyo", "Osaka")
ALTER TABLE restaurants 
ADD COLUMN area VARCHAR(100);

-- Create index for area queries
CREATE INDEX idx_restaurants_area ON restaurants(area);

-- Add comment
COMMENT ON COLUMN restaurants.area IS 'Administrative area level 1 from Google Maps (e.g., Tokyo, Osaka)';
