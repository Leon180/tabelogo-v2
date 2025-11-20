-- Create email_logs table for audit and tracking
CREATE TABLE IF NOT EXISTS email_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email_queue_id UUID,  -- Reference to email_queue (can be NULL if queue entry is deleted)
    recipient_email VARCHAR(255) NOT NULL,
    subject VARCHAR(500) NOT NULL,
    status VARCHAR(20) NOT NULL,  -- sent, failed, bounced, opened, clicked
    event_type VARCHAR(50),  -- delivered, opened, clicked, bounced, spam_report, unsubscribed, etc.
    external_id VARCHAR(255),  -- ID from email service provider
    error_message TEXT,
    metadata JSONB,  -- Additional tracking data (IP, user agent, etc.)
    webhook_data JSONB,  -- Raw webhook data from email service provider
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_email_logs_email_queue_id ON email_logs(email_queue_id);
CREATE INDEX idx_email_logs_recipient ON email_logs(recipient_email);
CREATE INDEX idx_email_logs_status ON email_logs(status);
CREATE INDEX idx_email_logs_event_type ON email_logs(event_type);
CREATE INDEX idx_email_logs_created_at ON email_logs(created_at DESC);
CREATE INDEX idx_email_logs_external_id ON email_logs(external_id);

-- Composite index for tracking email timeline
CREATE INDEX idx_email_logs_queue_timeline ON email_logs(email_queue_id, created_at DESC);

-- Create GIN indexes for JSONB columns
CREATE INDEX idx_email_logs_metadata ON email_logs USING GIN(metadata);
CREATE INDEX idx_email_logs_webhook_data ON email_logs USING GIN(webhook_data);

-- Add comments
COMMENT ON TABLE email_logs IS 'Email delivery logs and tracking events';
COMMENT ON COLUMN email_logs.email_queue_id IS 'Reference to email_queue.id (can be NULL)';
COMMENT ON COLUMN email_logs.status IS 'Email status: sent, failed, bounced, opened, clicked';
COMMENT ON COLUMN email_logs.event_type IS 'Event type from email service: delivered, opened, clicked, bounced, spam_report, unsubscribed, etc.';
COMMENT ON COLUMN email_logs.external_id IS 'ID from email service provider';
COMMENT ON COLUMN email_logs.webhook_data IS 'Raw webhook data from email service provider';
COMMENT ON COLUMN email_logs.metadata IS 'Additional tracking data (IP, user agent, etc.)';
