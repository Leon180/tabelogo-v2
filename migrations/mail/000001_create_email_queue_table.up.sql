-- Create email_queue table
CREATE TABLE IF NOT EXISTS email_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_email VARCHAR(255) NOT NULL,
    recipient_name VARCHAR(100),
    sender_email VARCHAR(255) DEFAULT 'noreply@tabelogo.com',
    sender_name VARCHAR(100) DEFAULT 'Tabelogo',
    subject VARCHAR(500) NOT NULL,
    body TEXT NOT NULL,
    html_body TEXT,  -- HTML version of the email
    template_name VARCHAR(100),  -- Email template name (e.g., welcome, booking_confirmation, password_reset)
    template_data JSONB,  -- Data to populate the template
    attachments JSONB DEFAULT '[]'::JSONB,  -- Array of attachment info
    priority INT DEFAULT 5 CHECK (priority >= 1 AND priority <= 10),  -- 1 is highest priority
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'sending', 'sent', 'failed', 'cancelled')),
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    scheduled_at TIMESTAMP DEFAULT NOW(),  -- When to send the email
    sent_at TIMESTAMP,
    failed_at TIMESTAMP,
    error_message TEXT,
    error_details JSONB,
    external_id VARCHAR(255),  -- ID from email service provider (e.g., SendGrid message ID)
    metadata JSONB,  -- Additional metadata (user_id, booking_id, etc.)
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_email_queue_status ON email_queue(status, scheduled_at);
CREATE INDEX idx_email_queue_priority ON email_queue(priority DESC, created_at);
CREATE INDEX idx_email_queue_recipient ON email_queue(recipient_email);
CREATE INDEX idx_email_queue_scheduled_at ON email_queue(scheduled_at);
CREATE INDEX idx_email_queue_created_at ON email_queue(created_at DESC);
CREATE INDEX idx_email_queue_template ON email_queue(template_name);

-- Index for pending emails to be sent (removed NOW() check due to immutability requirement)
CREATE INDEX idx_email_queue_pending ON email_queue(priority DESC, scheduled_at)
    WHERE status = 'pending';

-- Index for failed emails that can be retried
CREATE INDEX idx_email_queue_retry ON email_queue(retry_count, scheduled_at)
    WHERE status = 'failed' AND retry_count < max_retries;

-- Create GIN indexes for JSONB columns
CREATE INDEX idx_email_queue_metadata ON email_queue USING GIN(metadata);
CREATE INDEX idx_email_queue_template_data ON email_queue USING GIN(template_data);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create updated_at trigger
CREATE TRIGGER update_email_queue_updated_at BEFORE UPDATE ON email_queue
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE email_queue IS 'Email queue for asynchronous email sending';
COMMENT ON COLUMN email_queue.priority IS 'Email priority: 1 (highest) to 10 (lowest)';
COMMENT ON COLUMN email_queue.status IS 'Email status: pending, sending, sent, failed, cancelled';
COMMENT ON COLUMN email_queue.template_name IS 'Email template name (e.g., welcome, booking_confirmation, password_reset)';
COMMENT ON COLUMN email_queue.template_data IS 'Data to populate the template in JSON format';
COMMENT ON COLUMN email_queue.scheduled_at IS 'When to send the email (allows scheduling)';
COMMENT ON COLUMN email_queue.external_id IS 'ID from email service provider';
COMMENT ON COLUMN email_queue.metadata IS 'Additional metadata (user_id, booking_id, etc.) in JSON format';
