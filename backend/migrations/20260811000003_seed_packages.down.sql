-- +goose Down
DELETE FROM packages WHERE code IN ('free', 'basic', 'premium', 'pro');
