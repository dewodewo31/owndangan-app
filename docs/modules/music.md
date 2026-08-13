# Module: Music

## Purpose

Provide background music for the public invitation page. Uses preset tracks for Basic plan and allows MP3 upload for Premium/Pro plans. Music auto-plays when the invitation page is opened (standard wedding invitation UX in Indonesia).

## Responsibilities

- Provide a library of preset music tracks (Basic plan).
- Allow Premium/Pro users to upload custom MP3 files.
- Associate music with an event (one track at a time, or multiple with selection).
- Control music playback settings (autoplay, loop, volume via templates).
- Display music attribution/preset name.

## Non-Responsibilities

- Gallery videos (handled by Gallery module).
- Video background (not in scope).
- Music streaming or playlist management (single track or simple selection).
- Copyright clearance or licensing (user responsibility for uploads).

## Actors

- **User (couple):** Selects a preset track or uploads custom music.
- **Guest (public):** Hears background music on the invitation page (can toggle).
- **Admin:** Manages preset tracks (upload, activate, deactivate).

## Business Rules

- Basic plan: preset tracks only (from the system library).
- Premium/Pro plan: can upload MP3 files (up to 10MB, max 1 per event).
- Preset tracks are available to all plans (but Basic users can only use presets).
- `music.upload` capability gates custom upload (Premium/Pro only).
- Music file format: MP3 only (no other formats).
- Max file size: 10MB (reasonable for background music).
- Music autoplay is standard UX for Indonesian wedding invitations; provide a visible toggle for guests.
- One music configuration per event (1:1 via event_sections.music_id).
- Preset tracks are stored centrally; uploads are stored in object storage.
- Music file is deleted when the event is deleted or the music is changed.
- Volume control is handled by the frontend/template (default 0.3–0.5).

## Entities

- **Music:** `{ id, event_id, title, file_url, preset, is_preset, created_at }`
- **PresetTrack** (application config or DB table): `{ id, title, artist, file_url, duration, is_active }`

## Database

- Table: `music`
- `event_id`: nullable (for presets, event_id is null; for custom uploads, event_id is set).
- `is_preset`: boolean distinguishing system presets from user uploads.
- `preset`: identifier string for preset tracks (e.g., `canon_in_d`, `perfect_ed_sheeran`).
- `file_url`: storage URL for the audio file.
- `event_sections.music_id` references `music.id` for the event's selected track.

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/public/music/presets` | Public | List preset tracks |
| GET | `/api/v1/events/:id/music` | JWT | Get current music selection (returns `200` with `data: null` when none configured) |
| POST | `/api/v1/events/:id/music/upload` | JWT | Upload custom MP3 |
| GET | `/api/v1/events/:id/music/presets` | JWT | List preset tracks available to the caller's plan |
| POST | `/api/v1/events/:id/music/presets` | JWT | Select a preset track |
| DELETE | `/api/v1/events/:id/music` | JWT | Remove / disable music |

## Request Flow

```
POST /events/:id/music/upload
  → Handler: parse multipart form (file, title)
  → Service: verify event ownership
  → Service: check music.upload capability (Premium/Pro)
  → Service: validate file type (MP3) and size (max 10MB)
  → Service: upload to object storage, get URL
  → Service: create music record (is_preset=false), set event_sections.music_id
  → Handler: return created music (201)
```

```
POST /events/:id/music/presets
  → Handler: parse body { preset_id }
  → Service: verify event ownership
  → Service: load preset (must exist and be a preset)
  → Service: update event_sections.music_id to selected preset
  → Handler: return updated selection
```

```
DELETE /events/:id/music
  → Handler: parse event id
  → Service: verify event ownership
  → Service: clear event_sections.music_id and delete the selected music record
  → Handler: return 200 { message: "Music removed" }
```

## Validation

- Upload file: must be MP3 (audio/mpeg), max 10MB.
- Title: optional, max 255 chars.
- Preset ID: must match an existing active preset.
- File size checked at both application and reverse proxy level.

## Authorization

- Preset listing: public (no auth).
- Music selection/upload: JWT + event ownership + capability check.
- Admin: manages preset tracks.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing/invalid JWT |
| 403 | FORBIDDEN | Not owner, or music.upload not allowed by plan |
| 404 | NOT_FOUND | Event or preset not found |
| 413 | PAYLOAD_TOO_LARGE | File exceeds 10MB |
| 415 | UNSUPPORTED_MEDIA | Not MP3 |

## Security Considerations

- Uploaded MP3 files must be validated by content inspection (magic bytes), not just extension.
- Limit file size at multiple layers (reverse proxy, application).
- Store uploaded music in object storage outside web root.
- Serve music via streaming, not direct download (prevent hotlinking).
- Music URLs should use signed/pre-signed URLs if storage supports it.
- Preset tracks are system-controlled — no user modification.
- Ensure uploaded music does not overwrite existing files (use UUID filenames).

## Testing Requirements

- Unit tests for MP3 file type validation (magic bytes).
- Integration tests for preset selection.
- Integration tests for upload flow with entitlement.
- Test file size and type rejection.
- Test music removal (track deletion from storage).
- Test preset track listing (public).
- Test Basic plan cannot upload music.

## Dependencies

- Events module — event ownership.
- Subscriptions module — `music.upload` capability.
- Storage module — file upload and storage.
- Invitation Editor module — music_id in event_sections.

## Related Modules

- **Events** — Parent entity.
- **Invitation Editor** — Music section/settings.
- **Gallery** — Related media module.
- **Templates** — Music playback controls.

## Known Limitations

- No audio format conversion (accepts only MP3).
- No audio preview on selection (frontend concern).
- No playlist or multiple tracks per event.
- No crossfade or timed playback (e.g., start at a specific section).
- No copyright detection on uploaded tracks.
- No volume normalization across tracks.
- Preset tracks are loaded from backend config; no admin UI for managing them in current scope.

## TODO

- [ ] Add admin UI for managing preset tracks.
- [ ] Add audio format conversion (accept more formats, convert to MP3/OGG).
- [ ] Add audio preview in editor.
- [ ] Add playlist support (multiple tracks, sequential play).
- [ ] Add timed playback (synchronize with event schedule).
- [ ] Add copyright detection / takedown workflow.