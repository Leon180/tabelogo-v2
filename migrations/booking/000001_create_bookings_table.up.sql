-- Create bookings table
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,        -- Reference to auth_db.users (no FK due to microservices)
    restaurant_id UUID NOT NULL,  -- Reference to restaurant_db.restaurants (no FK due to microservices)
    booking_date TIMESTAMP NOT NULL,
    party_size INT NOT NULL CHECK (party_size > 0 AND party_size <= 50),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed', 'no_show')),
    external_booking_id VARCHAR(255),  -- OpenTable or other external service booking ID
    external_service VARCHAR(50),  -- 'opentable', 'tabelog', etc.
    special_requests TEXT,
    customer_name VARCHAR(100) NOT NULL,
    customer_phone VARCHAR(20) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    notes TEXT,  -- Internal notes
    confirmation_code VARCHAR(50) UNIQUE,  -- Unique confirmation code for the booking
    last_synced_at TIMESTAMP,  -- Last time synced with external service
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    cancelled_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL
);

-- Create indexes
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_restaurant_id ON bookings(restaurant_id);
CREATE INDEX idx_bookings_booking_date ON bookings(booking_date);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_confirmation_code ON bookings(confirmation_code);
CREATE INDEX idx_bookings_external_booking_id ON bookings(external_booking_id);
CREATE INDEX idx_bookings_created_at ON bookings(created_at DESC);

-- Composite indexes for common queries
CREATE INDEX idx_bookings_user_status ON bookings(user_id, status);
CREATE INDEX idx_bookings_restaurant_date ON bookings(restaurant_id, booking_date);
-- Partial index for upcoming bookings (removed NOW() check due to immutability requirement)
CREATE INDEX idx_bookings_upcoming ON bookings(booking_date, status)
    WHERE status IN ('pending', 'confirmed');

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create updated_at trigger
CREATE TRIGGER update_bookings_updated_at BEFORE UPDATE ON bookings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE bookings IS 'Restaurant booking records (synced with external services)';
COMMENT ON COLUMN bookings.user_id IS 'Reference to users.id from auth service (no FK)';
COMMENT ON COLUMN bookings.restaurant_id IS 'Reference to restaurants.id from restaurant service (no FK)';
COMMENT ON COLUMN bookings.status IS 'Booking status: pending, confirmed, cancelled, completed, no_show';
COMMENT ON COLUMN bookings.external_booking_id IS 'Booking ID from external service (e.g., OpenTable)';
COMMENT ON COLUMN bookings.external_service IS 'External service name: opentable, tabelog, etc.';
COMMENT ON COLUMN bookings.confirmation_code IS 'Unique confirmation code sent to customer';
COMMENT ON COLUMN bookings.party_size IS 'Number of people (1-50)';
COMMENT ON COLUMN bookings.last_synced_at IS 'Last synchronization time with external service';
