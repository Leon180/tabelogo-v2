-- Drop email_queue table
DROP TRIGGER IF EXISTS update_email_queue_updated_at ON email_queue;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_email_queue_template_data;
DROP INDEX IF EXISTS idx_email_queue_metadata;
DROP INDEX IF EXISTS idx_email_queue_retry;
DROP INDEX IF EXISTS idx_email_queue_pending;
DROP INDEX IF EXISTS idx_email_queue_template;
DROP INDEX IF EXISTS idx_email_queue_created_at;
DROP INDEX IF EXISTS idx_email_queue_scheduled_at;
DROP INDEX IF EXISTS idx_email_queue_recipient;
DROP INDEX IF EXISTS idx_email_queue_priority;
DROP INDEX IF EXISTS idx_email_queue_status;
DROP TABLE IF EXISTS email_queue;
