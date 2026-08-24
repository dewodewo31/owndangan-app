-- +goose Down
DROP TABLE IF EXISTS love_stories;
ALTER TABLE events DROP COLUMN IF EXISTS video_url;
