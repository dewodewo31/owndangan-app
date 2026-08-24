-- +goose Up
ALTER TABLE events ADD COLUMN IF NOT EXISTS video_url TEXT;

CREATE TABLE IF NOT EXISTS love_stories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    story TEXT NOT NULL,
    image_url TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_love_stories_event_id ON love_stories(event_id);

-- +goose Down
DROP TABLE IF EXISTS love_stories;
ALTER TABLE events DROP COLUMN IF EXISTS video_url;
