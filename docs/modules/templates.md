# Module: Templates

## Purpose

Manage the invitation visual themes (templates) available to users. Templates define the look, layout, and styling of invitations, grouped by tier (standard, premium, all). This module handles template lifecycle, group assignment, and versioning.

## Responsibilities

- Define templates with name, group, thumbnail, CSS config, and layout config.
- Group templates by tier: `standard`, `premium`, `all`.
- Control template availability based on the user's plan.
- Assign a template to an event.
- Activate/deactivate templates.
- Manage template assets (CSS, JS, thumbnails).

## Non-Responsibilities

- Section content management (handled by Invitation Editor module).
- Rendering the invitation (handled by the frontend).
- Uploading user media (handled by Gallery/Music modules).
- Custom template creation by users (future Pro feature).

## Actors

- **User (couple):** Selects a template for their event from the available group.
- **Admin:** Uploads, updates, activates, or deactivates templates.
- **Guest (public):** Views the invitation with the assigned template (read-only).

## Business Rules

- Templates are grouped: `standard` (Basic), `premium` (Premium), `all` (Pro).
- Basic/Free users can only use `standard` group templates.
- Premium users can use `standard` + `premium` groups.
- Pro users can use all templates and can request custom templates (`template.custom_request`).
- Template assignment is per event, not per subscription.
- Only `is_active = true` templates are selectable.
- Changing an event's template preserves the event content (content is template-agnostic).
- Template name is unique.
- Template deletion is not allowed if referenced by events; deactivate instead.
- Templates have no hard version history in current scope; updates overwrite config (see Known Limitations).
- `css_config` holds theme CSS variables; `layout_config` holds layout options (JSONB).

## Entities

- **Template:** `{ id, name, group_name, thumbnail_url, css_config, layout_config, is_active, created_at, updated_at }`

## Database

- Table: `templates`
- No soft delete (use `is_active`).
- Unique constraint on `name`.
- `group_name` values: `standard`, `premium`, `all`.
- `events.template_id` references `templates.id` (nullable).

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/templates` | JWT | List templates available to user's plan |
| PUT | `/api/v1/events/:id/template` | JWT | Assign template to event |
| GET | `/api/v1/admin/templates` | Admin | List all templates |
| POST | `/api/v1/admin/templates` | Admin | Create/upload template |
| PUT | `/api/v1/admin/templates/:id` | Admin | Update template |
| DELETE | `/api/v1/admin/templates/:id` | Admin | Deactivate template |

## Request Flow

```
GET /templates
  → Handler: parse JWT
  → Service: determine user's plan/template group entitlement
  → Service: load is_active templates for allowed groups
  → Handler: return template list (name, thumbnail, group)
```

```
PUT /events/:id/template
  → Handler: parse template_id
  → Service: verify event ownership
  → Service: verify template exists, is_active, and group allowed by plan
  → Service: update event.template_id
  → Handler: return updated event
```

## Validation

- name: required, unique, max 100 chars.
- group_name: one of `standard`, `premium`, `all`.
- css_config: valid JSON object of CSS variables.
- layout_config: valid JSON object.
- thumbnail_url: valid URL.
- Template assignment requires an active subscription matching the group.

## Authorization

- Template listing: JWT required (plan-dependent filter).
- Template assignment: JWT + event ownership + plan entitlement.
- Template CRUD: JWT with admin role.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing/invalid JWT |
| 403 | FORBIDDEN | Not owner, or template group not allowed by plan |
| 404 | NOT_FOUND | Template/event not found |
| 409 | CONFLICT | Duplicate template name, delete referenced template |
| 422 | VALIDATION_ERROR | Invalid config JSON |

## Security Considerations

- Template CSS/config is executed/rendered by the frontend — treat as trusted code from admins only, but sanitize on render.
- User-supplied content is injected into templates — rendering must escape content to prevent stored XSS.
- Template group entitlement must be enforced server-side; a Basic user must not receive premium templates even if they guess IDs.
- Admin template uploads should validate file types (images) and sizes.

## Testing Requirements

- Unit tests for group entitlement filtering.
- Integration tests for template assignment (success and forbidden cases).
- Admin CRUD tests.
- Test deactivation blocks new assignments.
- Test duplicate name rejection.
- Test template change preserves event content.
- Test content sanitization within template rendering.

## Dependencies

- Subscriptions module — plan/group entitlement.
- Packages module — `template_group` and `template.custom_request` capabilities.
- Events module — template assignment target.
- Storage (S3-compatible) — thumbnail and asset hosting.

## Related Modules

- **Events** — Template assigned to events.
- **Invitation Editor** — Content rendered within templates.
- **Packages** — Group entitlement mapping.
- **Admin** — Template management UI.

## Known Limitations

- No template versioning (config overwrites, no rollback).
- No user-created custom templates (Pro custom request is manual).
- No template analytics (which templates convert best).
- No localized template variants.
- Template assets are not hot-reloadable in production (redeploy required).

## TODO

- [ ] Implement template versioning (revision history, rollback).
- [ ] Add custom template request workflow for Pro users.
- [ ] Add template usage analytics.
- [ ] Add template asset hot-reload (storage-backed).
- [ ] Add template preview screenshots per device.