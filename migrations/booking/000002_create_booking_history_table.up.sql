-- Create booking_history table for event sourcing
CREATE TABLE IF NOT EXISTS booking_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,
    changed_by UUID,  -- user_id who made the change
    change_type VARCHAR(50) NOT NULL,  -- created, updated, confirmed, cancelled, completed, synced, etc.
    previous_value JSONB,  -- Previous state before change
    new_value JSONB,       -- New state after change
    notes TEXT,
    metadata JSONB,  -- Additional metadata about the change
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_booking_history_booking_id ON booking_history(booking_id);
CREATE INDEX idx_booking_history_created_at ON booking_history(created_at DESC);
CREATE INDEX idx_booking_history_change_type ON booking_history(change_type);
CREATE INDEX idx_booking_history_changed_by ON booking_history(changed_by);

-- Composite index for querying booking timeline
CREATE INDEX idx_booking_history_booking_timeline ON booking_history(booking_id, created_at DESC);

-- Add comments
COMMENT ON TABLE booking_history IS 'Event sourcing table for booking state changes';
COMMENT ON COLUMN booking_history.booking_id IS 'Reference to bookings.id';
COMMENT ON COLUMN booking_history.changed_by IS 'User ID who made the change (from auth service)';
COMMENT ON COLUMN booking_history.change_type IS 'Type of change: created, updated, confirmed, cancelled, completed, synced, etc.';
COMMENT ON COLUMN booking_history.previous_value IS 'Previous state in JSON format';
COMMENT ON COLUMN booking_history.new_value IS 'New state in JSON format';
COMMENT ON COLUMN booking_history.metadata IS 'Additional metadata (IP, user agent, sync source, etc.)';
