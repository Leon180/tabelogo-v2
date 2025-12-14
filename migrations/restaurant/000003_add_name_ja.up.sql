-- Add name_ja column to restaurants table for Japanese name support
ALTER TABLE restaurants 
ADD COLUMN name_ja VARCHAR(255);

-- Create index for Japanese name searches
CREATE INDEX idx_restaurants_name_ja ON restaurants(name_ja);

-- Add comment
COMMENT ON COLUMN restaurants.name_ja IS 'Japanese name of the restaurant for better Tabelog search results';
