-- +goose Up
INSERT INTO packages (id, name, code, price, duration_days, guest_limit, template_group, features, is_active, created_at, updated_at)
VALUES
(gen_random_uuid(), 'Free', 'free', 0, 7, 100, 'standard', '{"guest.max": 100, "theme.max": 5, "gallery.photo.max": 10, "gallery.video.max": 0, "video.enabled": false, "music.upload": false, "music.preset": true, "custom_domain": false, "watermark.removed": false, "whatsapp.bulk": false, "guestbook.qr": false, "rsvp.export": false, "digital_gift.qris": false, "template.custom_request": false, "event.max": 1}', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(gen_random_uuid(), 'Basic', 'basic', 50000, 90, 200, 'standard', '{"guest.max": 200, "theme.max": 10, "gallery.photo.max": 20, "gallery.video.max": 0, "video.enabled": false, "music.upload": true, "music.preset": true, "custom_domain": false, "watermark.removed": false, "whatsapp.bulk": false, "guestbook.qr": false, "rsvp.export": true, "digital_gift.qris": false, "template.custom_request": false, "event.max": 3}', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(gen_random_uuid(), 'Premium', 'premium', 150000, 365, 500, 'premium', '{"guest.max": 500, "theme.max": 25, "gallery.photo.max": 50, "gallery.video.max": 1, "video.enabled": true, "music.upload": true, "music.preset": true, "custom_domain": false, "watermark.removed": true, "whatsapp.bulk": true, "guestbook.qr": true, "rsvp.export": true, "digital_gift.qris": false, "template.custom_request": true, "event.max": 10}', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(gen_random_uuid(), 'Pro', 'pro', 350000, NULL, NULL, 'all', '{"guest.max": null, "theme.max": null, "gallery.photo.max": null, "gallery.video.max": 10, "video.enabled": true, "music.upload": true, "music.preset": true, "custom_domain": true, "watermark.removed": true, "whatsapp.bulk": true, "guestbook.qr": true, "rsvp.export": true, "digital_gift.qris": true, "template.custom_request": true, "event.max": null}', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (name) DO NOTHING;

-- +goose Down
DELETE FROM packages WHERE code IN ('free', 'basic', 'premium', 'pro');
