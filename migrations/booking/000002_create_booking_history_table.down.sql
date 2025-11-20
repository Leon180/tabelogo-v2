-- Drop booking_history table
DROP INDEX IF EXISTS idx_booking_history_booking_timeline;
DROP INDEX IF EXISTS idx_booking_history_changed_by;
DROP INDEX IF EXISTS idx_booking_history_change_type;
DROP INDEX IF EXISTS idx_booking_history_created_at;
DROP INDEX IF EXISTS idx_booking_history_booking_id;
DROP TABLE IF EXISTS booking_history;
