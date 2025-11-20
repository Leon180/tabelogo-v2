-- Drop bookings table
DROP TRIGGER IF EXISTS update_bookings_updated_at ON bookings;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_bookings_upcoming;
DROP INDEX IF EXISTS idx_bookings_restaurant_date;
DROP INDEX IF EXISTS idx_bookings_user_status;
DROP INDEX IF EXISTS idx_bookings_created_at;
DROP INDEX IF EXISTS idx_bookings_external_booking_id;
DROP INDEX IF EXISTS idx_bookings_confirmation_code;
DROP INDEX IF EXISTS idx_bookings_status;
DROP INDEX IF EXISTS idx_bookings_booking_date;
DROP INDEX IF EXISTS idx_bookings_restaurant_id;
DROP INDEX IF EXISTS idx_bookings_user_id;
DROP TABLE IF EXISTS bookings;
