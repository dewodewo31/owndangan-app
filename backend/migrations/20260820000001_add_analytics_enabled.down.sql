-- +goose Down
UPDATE packages SET features = features - 'analytics.enabled' WHERE code IN ('starter','premium','all');