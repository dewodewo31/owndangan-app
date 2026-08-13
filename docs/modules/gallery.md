# Module: Gallery

## Purpose

Manage the photo and video gallery displayed on the public invitation page. Allows couples to upload wedding photos and embed YouTube videos to showcase their journey.

## Responsibilities

- Upload and store wedding photos.
- Associate photos with captions and sort order.
- Embed YouTube videos (URL-based).
- Enforce photo count limits based on subscription plan.
- Enforce video availability based on subscription plan.
- Display gallery on the public invitation page.
- Manage photo deletion and reordering.

## Non-Responsibilities

- Background music (handled by Music module).
- Album organization (flat gallery, no albums/sub-albums).
- Image editing or filters (upload as-is).
- Video hosting (YouTube only — no video upload).

## Actors

- **User (couple):** Uploads photos, embeds videos, manages gallery.
- **Guest (public):** Views the gallery on the invitation.
- **Admin:** Read-only.

## Business Rules

- Photo limit per plan: Free/Basic=5, Premium=20, Pro=unlimited.
- Video limit: Free/Basic=0 (disabled), Premium=1, Pro=10 (or unlimited depending on `gallery.video.max`).
- Video is only available on plans with `video.enabled: true`.
- Gallery section must be enabled (`event_sections.gallery_enabled`).
- Photos are stored in S3-compatible object storage (or local in dev).
- Supported image formats: JPEG, PNG, WebP. Max file size: 5MB per photo.
- Video embed: YouTube URL only (not direct video file upload).
- Sort order is numeric (0-based); users can reorder via the editor.
- Deleting a photo frees up quota for additional uploads.
- Photos are hard-deleted (no soft delete).
- Captions are optional, max 255 characters.
- Gallery section toggle is independent of photo/video count enforcement.

## Entities

- **GalleryPhoto:** `{ id, event_id, image_url, caption, sort_order, created_at }`

## Database

- Table: `gallery_photos`
- Hard delete.
- Index on `event_id`, `sort_order`.
- Video URLs stored in `event_sections` or a dedicated `gallery_videos` table (TBD).

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/events/:id/gallery/photos` | JWT | List photos for event |
| POST | `/api/v1/events/:id/gallery/photos` | JWT | Upload photo (multipart) |
| PUT | `/api/v1/events/:id/gallery/photos/:pid` | JWT | Update photo caption/sort |
| DELETE | `/api/v1/events/:id/gallery/photos/:pid` | JWT | Delete photo |
| POST | `/api/v1/events/:id/gallery/photos/reorder` | JWT | Reorder photos (array of IDs) |
| PUT | `/api/v1/events/:id/gallery/video` | JWT | Set/update YouTube embed URL |
| GET | `/api/v1/public/events/:slug/gallery` | Public | Get gallery for public |

## Request Flow

```
POST /events/:id/gallery/photos
  → Handler: parse multipart form (file, caption)
  → Service: verify event ownership
  → Service: check subscription plan's photo limit
  → Service: validate file type (JPEG/PNG/WebP) and size (max 5MB)
  → Service: upload to object storage, get URL
  → Service: create gallery_photo record with URL, caption, sort_order
  → Handler: return created photo (201)
```

```
PUT /events/:id/gallery/video
  → Handler: parse body { video_url }
  → Service: verify event ownership + video.enabled capability
  → Service: validate YouTube URL format
  → Service: extract YouTube video ID, store embed URL
  → Handler: return updated video config
```

## Validation

- Photo file: JPEG, PNG, or WebP; max 5MB.
- Caption: optional, max 255 chars.
- Video URL: must be a valid YouTube URL (youtube.com or youtu.be).
- Sort order: non-negative integer.
- Photo count: checked against `gallery.photo.max` capability.
- Video count: checked against `gallery.video.max` capability.

## Authorization

- Photo/video management: JWT + event ownership.
- Public viewing: no auth (published event only).
- Admin: read-only.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing/invalid JWT |
| 403 | FORBIDDEN | Not owner, or feature not allowed |
| 404 | NOT_FOUND | Event or photo not found |
| 413 | PAYLOAD_TOO_LARGE | File exceeds 5MB |
| 415 | UNSUPPORTED_MEDIA | Invalid file type |
| 422 | LIMIT_EXCEEDED | Photo/video limit reached |

## Security Considerations

- Photo upload must validate file type by inspecting magic bytes, not just extension.
- Limit file size at the web server/reverse proxy level (Nginx/client_max_body_size).
- Uploaded files must be stored outside the web root (in object storage).
- Generated image URLs should be signed (pre-signed URLs for S3) to prevent unauthorized access.
- YouTube URL validation must prevent arbitrary URL injection (iframe XSS risk).
- Delete uploaded file from storage when photo record is deleted.

## Testing Requirements

- Unit tests for file type validation (magic bytes).
- Unit tests for YouTube URL parsing and validation.
- Integration tests for photo upload with limit enforcement.
- Test photo count enforcement per plan.
- Test video embedding per plan.
- Test reorder logic.
- Test file deletion removes from storage.
- Test malicious file upload rejection.

## Dependencies

- Events module — event ownership.
- Subscriptions module — gallery.photo.max, video.enabled, gallery.video.max capabilities.
- Storage module — object storage for photo uploads.
- Invitation Editor module — gallery section toggle.

## Related Modules

- **Events** — Parent entity.
- **Invitation Editor** — Gallery section toggle.
- **Music** — Related media module (gallery photos + background music).
- **Templates** — Gallery rendering styles.

## Known Limitations

- No image compression/resizing on upload.
- No album/sub-album organization.
- No slideshow or lightbox customization.
- No video upload (YouTube only).
- No preloaded stock photos.
- No drag-and-drop reorder in current scope (API supports numeric sort).

## TODO

- [ ] Add image compression/resizing pipeline.
- [ ] Add album organization (albums/sub-albums).
- [ ] Add slideshow configuration (autoplay speed, transition).
- [ ] Add video upload support (future, requires transcoding).
- [ ] Add image optimization (WebP conversion, responsive sizes).