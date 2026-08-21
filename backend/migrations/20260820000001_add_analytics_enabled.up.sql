-- +goose Up
UPDATE packages SET features = features || '{"analytics.enabled":true}'::jsonb WHERE code IN ('starter','premium','all');