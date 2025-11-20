-- Drop email_logs table
DROP INDEX IF EXISTS idx_email_logs_webhook_data;
DROP INDEX IF EXISTS idx_email_logs_metadata;
DROP INDEX IF EXISTS idx_email_logs_queue_timeline;
DROP INDEX IF EXISTS idx_email_logs_external_id;
DROP INDEX IF EXISTS idx_email_logs_created_at;
DROP INDEX IF EXISTS idx_email_logs_event_type;
DROP INDEX IF EXISTS idx_email_logs_status;
DROP INDEX IF EXISTS idx_email_logs_recipient;
DROP INDEX IF EXISTS idx_email_logs_email_queue_id;
DROP TABLE IF EXISTS email_logs;
