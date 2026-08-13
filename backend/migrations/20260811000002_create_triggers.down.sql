-- +goose Down
DROP TRIGGER IF EXISTS trigger_users_updated_at ON users;
DROP TRIGGER IF EXISTS trigger_packages_updated_at ON packages;
DROP TRIGGER IF EXISTS trigger_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS trigger_subscriptions_updated_at ON subscriptions;
DROP TRIGGER IF EXISTS trigger_templates_updated_at ON templates;
DROP TRIGGER IF EXISTS trigger_events_updated_at ON events;
DROP TRIGGER IF EXISTS trigger_event_sections_updated_at ON event_sections;
DROP TRIGGER IF EXISTS trigger_guests_updated_at ON guests;
DROP TRIGGER IF EXISTS trigger_guestbook_messages_updated_at ON guestbook_messages;
DROP TRIGGER IF EXISTS trigger_digital_gifts_updated_at ON digital_gifts;
DROP FUNCTION IF EXISTS update_updated_at_column();
