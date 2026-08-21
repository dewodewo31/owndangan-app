# owndangan-feature-wave - Work Plan

## TL;DR (For humans)
<!-- Fill this LAST, after the detailed plan below is written, so it summarizes the REAL plan. -->
<!-- Plain English for a non-engineer: NO file paths, NO todo numbers, NO wave/agent/tool names. -->

**What you'll get:** Live preview di editor yang memakai renderer undangan asli (bukan versi sederhana), filter template berdasarkan acara + template corporate baru, fitur sampah & pulihkan tamu, dashboard analitik (klik WA/maps/telepon + jumlah tampilan unik, khusus pengguna berbayar), mode gelap otomatis malam hari, email notifikasi (sukses bayar + pengingat kedaluwarsa), halaman FAQ/cara-order/testimoni, plus konsistensi harga & bahasa dan pembersihan repo.

**Why this approach:** Semua fitur dibangun di atas infrastruktur yang sudah ada (tabel analitik, entitlement, soft-delete tamu, audit log) — tanpa tabel baru, tanpa sistem antrean baru, tanpa ubah data harga. Harga backend dijadikan satu-satunya sumber kebenaran, tampilan tinggal disinkronkan.

**What it will NOT do:** Tidak mengubah data harga/paket di database, tidak membuat testimoni palsu, tidak merombak tampilan landing page, tidak menambah sistem billing baru, tidak mengaktifkan langganan dari sisi frontend.

**Effort:** XL
**Risk:** Medium - banyak lintas-stack (backend + frontend + docs) dalam 8 gelombang, tapi semuanya memakai primitif yang sudah ada.
**Decisions to sanity-check:** (1) sinkronisasi harga hanya tampilan (data backend tetap), (2) analitik khusus pengguna berbayar, (3) testimoni berupa placeholder, (4) email tanpa antrean (goroutine in-process), (5) kategori acara tanpa perubahan skema database.

Your next move: approve plan ini untuk dieksekusi worker, atau jalankan high-accuracy review dulu. Detail eksekusi lengkap di bawah.

---

> TL;DR (machine): 8 waves / 34 todos + 4 final audits (F1-F4); effort XL, risk Medium; reuses analytics_events, entitlement, guest soft-delete, audit_log — no new tables, no queue, no pricing-data change; tests-after per user decision (Q3).

## Scope
### Must have
- W1 Live Preview: editor (dashboard/editor) renders the REAL template renderer (`TemplateShell` + `selectTemplate` + `buildInvitationModel` via new adapter) inside a virtual-phone frame, live from editor state (no API call per keystroke, no save-to-preview roundtrip). Desktop: form left / sticky preview right. Mobile: preview as tab/toggle.
- W2 Template occasions: new `occasions: string[]` field on `TemplateDefinition` (types.ts) + values on all 10 existing templates; occasion filter UI in editor template tab AND marketing template showcase; new corporate template (`frontend/src/templates/corporate/`) with distinct non-wedding design language + seed row "Modern Corporate" in `backend/internal/database/seed.go` (count-check pattern, no migration).
- W3 Guest restore + Analytics: backend restore endpoint + list-deleted support; guests-page trash UI with restore; analytics click tracking (whatsapp_click, map_click, phone_click via new public POST endpoint writing existing `analytics_events` table), per-event dashboard analytics endpoint (unique views = COUNT DISTINCT ip_address of page_view per event), dashboard widget gated paid-only via `entitlement` resolver (free sees Event.ViewCount only); ICS timezone polish (TZID) in `sections.tsx`.
- W4 Dark/night mode: per-template `night` ThemeTokens variant + auto window 18:00–06:00 (event timezone, fallback browser), client-side/hydration-safe.
- W5 Email: gomail SMTP + config keys (.env.example + config.go) + EmailService + payment-success email on settlement webhook (async, never fails transaction) + daily expiry-reminder ticker (7 days) with dedupe via audit-log ledger (no schema change).
- W6 Trust pages: `/faq`, `/cara-order`, `/testimoni` (testimonials EXPLICITLY placeholder, no fabricated identities); fix dead `/packages` link in billing; nav/footer links.
- W7 Consistency: brand/copy normalized to Indonesian (layout.tsx metadata, marketing.ts, navbar/footer/dashboard/auth shells; brand stays "Owndangan"); pricing display synced to backend seed (Starter 99rb/30hr, Premium 299rb/60hr, All Access 999rb lifetime) — DISPLAY ONLY, no data change; onboarding auto-login (register page consumes returned tokens → /dashboard).
- W8 Docs + hygiene: public/internal doc split (docs/README.md restructured, dev-only docs moved under docs/internal/); delete root `-l` + `backend/-l`; .gitignore += `-l`, `*.tsbuildinfo`; `git rm --cached frontend/tsconfig.tsbuildinfo`; docs pricing ranges fixed to canonical values.
- Tests: backend Go tests for every new endpoint/behavior; frontend jest installed + unit tests for pure logic; final gates `make test`, `npm run build`, `npm run lint` with NEW vs PRE-EXISTING failures separated.

### Must NOT have (guardrails, anti-slop, scope boundaries)
- NO database schema change for template occasions (frontend mapping only).
- NO pricing/data change in `packages` table or seed VALUES (display sync only; user-approved Q1).
- NO new billing system; NO subscription activation from frontend callback (AGENTS.md forbidden).
- NO fabricated testimonials / customer identities (placeholder marked clearly).
- NO redesign of landing page or public renderer; NO removal of existing features.
- NO new queue system beyond in-process goroutine; NO third-party analytics SDK.
- NO rewrite of the simplified `InvitationPreview` — keep as fallback until LivePreview verified, then it may be deleted.
- NO tracking of password/token/private guest PII; analytics stores only event_id/type/ip/user_agent (existing columns).
- NO touching unrelated dirty-worktree WIP files (30+ modified, ~15 untracked; single initial commit 7d31e8c).

## Verification strategy
> Zero human intervention - all verification is agent-executed.
- Test decision: tests-after (user-confirmed Q3). Framework: Go `go test ./internal/...` (unit + integration + API via Makefile targets); frontend jest (to be installed in todo 1) for pure logic + component smoke; `npm run build` + `npm run lint` as final gates.
- Evidence: `.omo/evidence/owndangan-feature-wave/task-<N>-owndangan-feature-wave.<ext>` for every todo; final wave evidence under `.omo/evidence/owndangan-feature-wave/final-F<n>.`
- Every backend todo runs its tests with the exact Makefile target; every frontend todo runs `npx jest <file>` (or `npm run build` where noted); evidence files record the exact command, exit code, and assertion output (reject passing-logs-only claims).

## Execution strategy
### Parallel execution waves
- Wave 1 (Conversion - Live Preview): todos 1-5 (frontend).
- Wave 2 (Conversion - Templates): todos 6-11 (frontend + 1 backend seed row).
- Wave 3 (Retention - Guests+Analytics+ICS): todos 12-17 (backend + frontend).
- Wave 4 (Retention - Dark mode): todos 18-20 (frontend).
- Wave 5 (Professional - Email): todos 21-23 (backend).
- Wave 6 (Professional - Trust pages): todos 24-27 (frontend).
- Wave 7 (Consistency - Brand/Pricing/Onboarding): todos 28-31 (frontend).
- Wave 8 (Consistency - Docs/Hygiene/Gates): todos 32-34 (repo-wide).
Waves are independent except backend-frontend intra-wave dependencies; run waves in order, todos within a wave in order (each todo lists its blockers).

### Dependency matrix
| Todo | Depends on | Blocks | Can parallelize with |
| --- | --- | --- | --- |
| 1 (jest infra) | - | 2,3,4,5 | - |
| 2 (adapter) | 1 | 3,4,5 | - |
| 3 (LivePreview cmp) | 2 | 4,5 | - |
| 4 (editor wiring) | 3 | 5 | - |
| 5 (preview tests) | 1,4 | - | - |
| 6 (occasions field) | 1 | 7,8,9 | - |
| 7 (editor filter) | 6 | - | 8,9 |
| 8 (showcase filter) | 6 | - | 7,9 |
| 9 (corporate tpl) | 6 | 10,11 | 7,8 |
| 10 (seed row) | 9 | 11 | 7,8 |
| 11 (occasions tests) | 6,9 | - | - |
| 12 (guest restore BE) | - | 13 | 14 |
| 13 (trash UI) | 12 | - | - |
| 14 (analytics BE) | - | 15,16 | 12 |
| 15 (click tracking FE) | 14 | 16 | - |
| 16 (analytics widget) | 14,15 | - | - |
| 17 (ICS TZID) | 1 | - | 12,14 |
| 18 (night tokens) | 1 | 19,20 | - |
| 19 (night logic) | 18 | 20 | - |
| 20 (night tests) | 18,19 | - | - |
| 21 (email infra) | - | 22,23 | 12,14,18 |
| 22 (payment email) | 21 | - | 23 |
| 23 (expiry ticker) | 21 | - | 22 |
| 24 (/faq) | 1 | 26 | - |
| 25 (/cara-order + /testimoni) | 1 | 26 | 24 |
| 26 (nav/footer links + /packages fix) | 24,25 | 27 | - |
| 27 (trust pages tests) | 24,25,26 | - | - |
| 28 (brand normalization) | 1 | - | 29 |
| 29 (pricing display sync) | 1 | - | 28 |
| 30 (onboarding auto-login) | 1 | 31 | 28,29 |
| 31 (pricing+onboarding tests) | 29,30 | - | - |
| 32 (docs split) | 29 | - | 33 |
| 33 (-l + gitignore) | - | 34 | 32 |
| 34 (final gates) | 32,33 + all | - | - |
| F1-F4 | 34 | - | parallel with each other |

## Todos
> Implementation + Test = ONE todo. Never separate.
<!-- APPEND TASK BATCHES BELOW THIS LINE WITH edit/apply_patch - never rewrite the headers above. -->
- [x] 1. Frontend test infra: install jest + ts-jest + jest-environment-jsdom in frontend devDependencies, add jest.config.js (ts-jest preset, testMatch `src/**/*.test.{ts,tsx}`), add a smoke test `src/lib/__tests__/smoke.test.ts` asserting `1+1===2`; fix package.json "test" script to `jest --ci`.
  What to do / Must NOT do: npm install -D jest @types/jest ts-jest jest-environment-jsdom; configure; run `npm test` green. Must NOT remove existing deps; must NOT change app code.
  Parallelization: Wave 1 | Blocked by: - | Blocks: 2,3,4,5
  References (executor has NO interview context - be exhaustive): `frontend/package.json:11,38-44` (test script exists, jest missing), `frontend/tsconfig.json`, `frontend/jest.config.js` (does not exist — create).
  Acceptance criteria (agent-executable): `npm test` exits 0 and prints "1 passed"; `npx jest --listTests` shows the smoke test.
  QA scenarios (name the exact tool + invocation): happy — `npm test` → exit 0, 1 passed (Evidence `.omo/evidence/owndangan-feature-wave/task-1-owndangan-feature-wave.txt`); failure — temporarily break smoke assertion, `npm test` must exit non-zero with failing test name, then revert (Evidence same file, second run).
  Commit: Y | chore(frontend): add jest test infrastructure
- [x] 2. Editor→model adapter: create `frontend/src/lib/editor-preview.ts` exporting `editorStateToPublicEvent(event, sections, gallery, music, gift, loveStories, template) → PublicEventResponse` (reuse the shape from `frontend/src/lib/invitation.ts:5-59`) and `buildPreviewModel(...)` = `buildInvitationModel(adapterOutput, new URLSearchParams(), event.slug)`; map wedding_date/wedding_time, ceremony/reception fields, template css_config (from TemplateSummary), sections booleans, verse, music, gallery, love_stories, digital_gift, video.
  What to do / Must NOT do: pure function, no fetch; must NOT call save API; must NOT mutate inputs; must keep field names identical to PublicEventResponse so buildInvitationModel works unchanged.
  Parallelization: Wave 1 | Blocked by: 1 | Blocks: 3,4,5
  References: `frontend/src/lib/invitation.ts:5-59` (PublicEventResponse), `:83-171` (buildInvitationModel), `frontend/src/lib/types.ts:58-160` (WeddingEvent, EventSections, GalleryPhoto, Music, DigitalGift, TemplateSummary, LoveStory), editor state shape at `frontend/src/app/dashboard/editor/page.tsx:70-84`.
  Acceptance criteria (agent-executable): `npx jest src/lib/editor-preview` passes adapter mapping tests (todo 5 covers); `npm run build` type-checks clean.
  QA scenarios: happy — jest test maps a full WeddingEvent and asserts every PublicEventResponse field (names, dates, venues, sections flags, gift) (Evidence task-2 file); failure — test with null event must not throw (returns null) (same file).
  Commit: Y | feat(frontend): add editor-state to invitation model adapter
- [x] 3. LivePreview component: create `frontend/src/components/dashboard/live-preview.tsx` ("use client") rendering the REAL public renderer inside a phone frame: `selectTemplate(template.group_name, template.name)` + `TemplateShell` + `themeVars` from the template definition theme, fed by `buildPreviewModel`; phone frame = rounded-3xl border + notch bar + width ~375px centered on a subtle backdrop; render `CoverGate` too (it is part of the public page: `[slug]/page.tsx:51-55`).
  What to do / Must NOT do: import ONLY from `@/components/invitation/shell` (TemplateShell), `@/components/invitation/cover-gate` (CoverGate), `@/templates` (selectTemplate), `@/templates/theme` (themeVars), `@/lib/editor-preview`; must NOT fetch from API; must NOT use the simplified InvitationPreview internals; keep the public page ([slug]/page.tsx) untouched.
  Parallelization: Wave 1 | Blocked by: 2 | Blocks: 4,5
  References: `frontend/src/app/[slug]/page.tsx:39-56` (public renderer composition), `frontend/src/components/invitation/shell.tsx:30` (TemplateShell), `frontend/src/components/invitation/cover-gate.tsx`, `frontend/src/templates/index.ts:35-51` (selectTemplate), `frontend/src/templates/theme.ts:5-29` (themeVars), template defs e.g. `frontend/src/templates/modern-minimalist/index.tsx`.
  Acceptance criteria (agent-executable): jest smoke render (react-test-renderer or @testing-library/react — add if needed) renders LivePreview without throwing for a filled event; `npm run build` clean.
  QA scenarios: happy — jest render with sample data asserts TemplateShell content (couple names text present) (Evidence task-3 file); failure — render with `event=null` shows placeholder text "Pratinjau akan muncul setelah undangan dibuat" without crashing (same file).
  Commit: Y | feat(frontend): add real-renderer live preview component with phone frame
- [x] 4. Editor wiring: in `frontend/src/app/dashboard/editor/page.tsx` replace the `<InvitationPreview .../>` usage (lines ~478-488) with `<LivePreview event={event} sections={sections} gallery={gallery} music={music} gift={gift} loveStories={loveStories} template={templates.find(t=>t.id===event?.template_id)||null} />`; keep desktop grid `lg:col-span-5 lg:sticky lg:top-6`; add mobile preview toggle: on small screens show a floating "Pratinjau" button opening the LivePreview in an overlay/tab (state `showMobilePreview`); keep all save/publish handlers and tab rendering intact; the existing `previewURL()` "Pratinjau" link (page.tsx:384-390) stays for published events.
  What to do / Must NOT do: only swap the preview column; do NOT touch renderDetail/renderSections/renderGallery/renderLoveStories/renderMusic/renderGifts/renderTemplate logic; do NOT alter fetchConfig/save debounce; do NOT add any API call.
  Parallelization: Wave 1 | Blocked by: 3 | Blocks: 5
  References: `frontend/src/app/dashboard/editor/page.tsx:418-490` (layout grid + preview column), `:350-390` (header/buttons), `:496-991` (render* tab functions — untouched).
  Acceptance criteria (agent-executable): `npm run build` clean; manual/agent browser check (playwright): editing the name input updates the preview text without page reload and without network request to /events (observe via devtools — assert with playwright route interception in QA).
  QA scenarios: happy — playwright: open /dashboard/editor, type in "Nama Mempelai" field, assert preview hero text updates within 500ms with zero XHR to /api/v1/events (Evidence task-4 screenshot+network log); failure — open editor on mobile viewport (375px), preview hidden by default, toggle button opens overlay (Evidence task-4 screenshot).
  Commit: Y | feat(frontend): wire live preview into editor layout
- [x] 5. Preview adapter/component tests: `frontend/src/lib/__tests__/editor-preview.test.ts` covering adapter mapping (full event, partial event, null event) and `frontend/src/components/dashboard/__tests__/live-preview.test.tsx` (renders real template, placeholder on null, template switch changes themeVars).
  What to do / Must NOT do: tests only; must NOT mock buildInvitationModel (test the real chain); keep tests hermetic (no network).
  Parallelization: Wave 1 | Blocked by: 1,4 | Blocks: -
  References: adapter (todo 2 output), LivePreview (todo 3 output), `frontend/src/lib/invitation.ts:83-171`.
  Acceptance criteria (agent-executable): `npx jest src/lib/__tests__/editor-preview.test.ts src/components/dashboard/__tests__/live-preview.test.tsx` — all pass, exit 0.
  QA scenarios: happy — full-event mapping asserts every field; failure — null-event returns null without throw (Evidence task-5 file).
  Commit: Y | test(frontend): cover editor-preview adapter and live preview
- [x] 6. Occasions field: add `occasions: string[]` to `TemplateDefinition` (`frontend/src/templates/types.ts:129-143`); set values in ALL 10 template definitions (modern-botanical, sundanese, rustic-bohemian, romantic-elegant, modern-minimalist, islamic, contemporary-editorial, luxury-black-gold, japanese-zen, javanese under `frontend/src/templates/*/index.ts:6`): wedding-themed → `["pernikahan"]` (javanese/sundanese/islamic also `["aqiqah"]`? NO — keep factual: only pernikahan unless a template genuinely fits another occasion; romantic-elegant/luxury-black-gold/islamic/rustic-bohemian/contemporary-editorial/japanese-zen/modern-minimalist/modern-botanical/sundanese/javanese all `["pernikahan"]` for now; corporate template (todo 9) gets `["corporate","event"]`); export a helper `templateOccasions()` and constants `OCCASION_LABELS: Record<string,string>` (pernikahan=Pernikahan, prewedding=Prewedding, ulang-tahun=Ulang Tahun, aqiqah=Aqiqah, khitanan=Khitanan, anniversary=Anniversary, graduation=Graduation, corporate=Corporate/Meeting, event=Event) in `frontend/src/templates/index.ts`.
  What to do / Must NOT do: must NOT change `category` semantics (style) or theme tokens; must NOT touch backend; all 10 files updated with the same single-item array (except corporate later).
  Parallelization: Wave 2 | Blocked by: 1 | Blocks: 7,8,9
  References: `frontend/src/templates/types.ts:129-143`, `frontend/src/templates/index.ts:13-24`, 10 def files under `frontend/src/templates/*/index.ts:1-10`.
  Acceptance criteria (agent-executable): `npx jest src/templates` (new test in todo 11) green; `npm run build` clean; grep confirms all 10 defs have `occasions`.
  QA scenarios: happy — helper returns correct labels; failure — a template missing `occasions` fails a type-level test (todo 11) (Evidence task-6 file).
  Commit: Y | feat(frontend): add occasion tags to template definitions
- [x] 7. Editor occasion filter: in `frontend/src/app/dashboard/editor/page.tsx` renderTemplate() add a filter bar (All + each occasion present in TEMPLATES via `templateOccasions()`), filtering the DB template list by mapping `t.name` → frontend definition (`selectTemplate(undefined, t.name)` or `TEMPLATES.find(x=>x.name.toLowerCase()===t.name.toLowerCase())`) and matching `occasions`; show occasion badge under each template card.
  What to do / Must NOT do: client-side filter only; must NOT change the assignTemplate call; must NOT add backend params; must NOT reorder templates otherwise.
  Parallelization: Wave 2 | Blocked by: 6 | Blocks: -
  References: `frontend/src/app/dashboard/editor/page.tsx` renderTemplate (grep `function renderTemplate`), `frontend/src/templates/index.ts` (TEMPLATES + OCCASION_LABELS), TemplateSummary at `frontend/src/lib/types.ts:147`.
  Acceptance criteria (agent-executable): `npm run build` clean; jest/playwright: selecting "Corporate" filter shows only the corporate template (after todo 9) (Evidence task-7 file).
  QA scenarios: happy — filter click narrows list; failure — template whose name has no frontend definition still shows under "Semua" (fallback, not hidden) (Evidence task-7).
  Commit: Y | feat(frontend): add occasion filter to editor template picker
- [x] 8. Showcase occasion filter: `frontend/src/components/marketing/template-showcase.tsx` — add the same occasion filter UI using TEMPLATES + OCCASION_LABELS (replace/augment the current hardcoded templateCategories from `frontend/src/data/marketing.ts:106` if used); show occasion badge on cards.
  What to do / Must NOT do: reuse the same constants; must NOT change landing layout beyond the filter; must NOT fabricate template thumbnails (keep existing image logic).
  Parallelization: Wave 2 | Blocked by: 6 | Blocks: -
  References: `frontend/src/components/marketing/template-showcase.tsx`, `frontend/src/data/marketing.ts:67-106`.
  Acceptance criteria (agent-executable): `npm run build` clean; playwright screenshot shows filter chips and filtering works (Evidence task-8).
  QA scenarios: happy — click "Corporate" shows corporate card; failure — zero-match filter shows empty-state text instead of broken grid (Evidence task-8).
  Commit: Y | feat(frontend): add occasion filter to template showcase
- [x] 9. Corporate template: create `frontend/src/templates/corporate/index.tsx` (definition + sections) with kind `"corporate"`, category `"Corporate"`, occasions `["corporate","event"]`, theme: neutral dark-slate/indigo palette, geometric decoration (`decoration: "geometric"`), no floral/batik, `nav: "bottom-floating"`; sections ordered for business events (cover with event title/subtitle, details/agenda via event blocks, location, rsvp, gallery, closing); register in `frontend/src/templates/index.ts:13-24`; add to templates list in marketing showcase data if thumbnails are static (fallback: definition thumbnail).
  What to do / Must NOT do: must be visually distinct from wedding templates (no couple names assumption — render names.full fallback "Undangan"; use model.events not wedding-specific labels; graceful when `names.full` is generic); must NOT break `selectTemplate` fallback (tier map stays); must NOT use wedding-only sections (no love-story/verse by default).
  Parallelization: Wave 2 | Blocked by: 6 | Blocks: 10,11
  References: template structure `frontend/src/templates/modern-minimalist/index.tsx` (pattern to copy), types `frontend/src/templates/types.ts:129-143`, shell/sections `frontend/src/components/invitation/shell.tsx:30`, `sections.tsx`, decoration styles `frontend/src/components/invitation/decorations.tsx`.
  Acceptance criteria (agent-executable): `npm run build` clean; jest render of corporate template with a corporate-style model (generic names, event blocks) renders without throwing; playwright screenshot shows non-wedding visual (Evidence task-9).
  QA scenarios: happy — template renders meeting-style event; failure — rendering with a wedding model (no events) still renders (falls back to date/time) without crash (Evidence task-9).
  Commit: Y | feat(frontend): add corporate invitation template
- [x] 10. Corporate seed row: in `backend/internal/database/seed.go` add a count-check block (pattern of seed.go:154-167) inserting `{Name: "Modern Corporate", GroupName: "standard", ThumbnailURL: <neutral office/event unsplash image>, CSSConfig: dark-slate/indigo palette matching the corporate theme, LayoutConfig: {"hero":"centered","gallery":"grid"}, IsActive: true}`; run `go build ./...` + `make test-unit`.
  What to do / Must NOT do: follow the existing count-check pattern EXACTLY (idempotent on existing DBs); must NOT modify existing seed rows or prices; must NOT change `model.Template`.
  Parallelization: Wave 2 | Blocked by: 9 | Blocks: 11
  References: `backend/internal/database/seed.go:70-167` (pattern), `backend/internal/model/model.go:151-161` (Template model fields).
  Acceptance criteria (agent-executable): `go build ./...` exit 0; `make test-unit` passes; unit/integration test asserts the row is seeded once (count-check) — extend `backend/internal/database` tests if present, else assert via a temporary integration check in `make test-integration` style test.
  QA scenarios: happy — seed inserts exactly one row; failure — running seed twice inserts zero additional rows (Evidence task-10 file).
  Commit: Y | feat(backend): seed corporate template row
- [x] 11. Occasions/corporate tests: `frontend/src/templates/__tests__/occasions.test.ts` — every TEMPLATES entry has non-empty `occasions` and valid keys from OCCASION_LABELS; corporate definition has `["corporate","event"]`; `selectTemplate` still returns a definition for unknown names (fallback).
  What to do / Must NOT do: tests only; must NOT change selectTemplate behavior.
  Parallelization: Wave 2 | Blocked by: 6,9 | Blocks: -
  References: `frontend/src/templates/index.ts`, `frontend/src/templates/types.ts`.
  Acceptance criteria (agent-executable): `npx jest src/templates/__tests__/occasions.test.ts` exit 0.
  QA scenarios: happy — all 11 templates validated; failure — temporarily set an empty occasions array on one template → test fails, revert (Evidence task-11 file).
  Commit: Y | test(frontend): validate occasion metadata and corporate template
- [x] 12. Guest restore backend: add `Restore(ctx, eventID, guestID, userID)` to `backend/internal/service/guest/service.go` (load guest, verify ownership via event.UserID, call new repo method `Restore` = `db.Unscoped().Model(&Guest{}).Where("id=? AND event_id=?", id, eventID).Update("deleted_at", nil)`, guard double-restore), add `ListDeleted(ctx, eventID)` repo method (`Unscoped().Where("event_id=? AND deleted_at IS NOT NULL")`); add handler route in `guest_handler.go` RegisterRoutes: `r.Post("/{guestID}/restore", h.Restore)`; extend `guest_repo.go` (repo already has SoftDelete at guest_repo.go:86-87). Also extend `List` with optional `?include_deleted=true` query (or rely on ListDeleted; pick ListDeleted — do NOT touch existing List contract).
  What to do / Must NOT do: ownership check identical to Delete (service.go:168-182 pattern); must NOT hard-delete anything; must NOT change RSVP records (deleted guest's RSVP stays; guest reappears with same ID — no duplicates because PK unchanged); must NOT alter existing DELETE behavior.
  Parallelization: Wave 3 | Blocked by: - | Blocks: 13
  References: `backend/internal/service/guest/service.go:168-182` (Delete pattern), `backend/internal/repository/guest_repo.go:86-87` (SoftDelete), `backend/internal/model/model.go:163-176` (Guest.DeletedAt), `backend/internal/api/handler/guest_handler.go:107-161` (Delete handler + routes).
  Acceptance criteria (agent-executable): `make test-unit` + `make test-integration` pass; new tests: delete→list hides guest, restore→list shows guest again, restore wrong-user returns 404/403, double-restore idempotent.
  QA scenarios: happy — create guest, delete, restore, assert present with same ID and RSVP intact (Evidence task-12 file); failure — restoring another user's guest returns error (Evidence same file).
  Commit: Y | feat(backend): add guest restore endpoint with ownership check
- [x] 13. Guests trash UI: in `frontend/src/app/dashboard/guests/page.tsx` add a "Sampah" (Trash) toggle listing deleted guests (`api.get('/events/{id}/guests?deleted=1')` or the ListDeleted endpoint from todo 12), each row with "Pulihkan" (Restore) button calling `api.post('/events/{id}/guests/{gid}/restore')` then refreshing both lists; active list unchanged; deleted guests excluded from counts.
  What to do / Must NOT do: reuse existing list/table UI and styling; must NOT change delete flow; must NOT touch RSVP/import.
  Parallelization: Wave 3 | Blocked by: 12 | Blocks: -
  References: `frontend/src/app/dashboard/guests/page.tsx:88-98,206` (delete flow), `frontend/src/lib/api.ts` (client), backend routes todo 12.
  Acceptance criteria (agent-executable): `npm run build` clean; playwright: delete a guest → appears in trash; restore → back in active list; active count excludes deleted.
  QA scenarios: happy — restore flow end-to-end; failure — restore API error shows alert and list unchanged (Evidence task-13 file).
  Commit: Y | feat(frontend): add guest trash and restore UI
- [x] 14. Analytics backend: (a) public endpoint `POST /api/v1/analytics/events` (no auth; rate-limited by existing per-IP middleware) accepting `{event_id, type}` with type ∈ whatsapp_click|map_click|phone_click, resolving event must be published, inserting `model.AnalyticsEvent{EventID, EventType, IPAddress (server-side r.RemoteAddr), UserAgent}` — reuse `analyticsRepo.Create`; (b) owner endpoint `GET /api/v1/events/{id}/analytics` (authRequired + ownership via EventService) returning `{views, unique_views, whatsapp_clicks, map_clicks, phone_clicks, rsvp_count}` where views=CountByType('page_view'), unique_views=COUNT(DISTINCT ip_address) WHERE event_id AND type='page_view' (new repo method `CountUniqueByEvent`), clicks=CountByType per type; (c) extend `backend/internal/service/entitlement/entitlement.go` with `AnalyticsEnabled() bool` (false for code "free"; otherwise features JSON `analytics.enabled` defaulting true) + add additive migration `backend/migrations/<next>_add_analytics_enabled.up.sql` = `UPDATE packages SET features = features || '{"analytics.enabled":true}'::jsonb WHERE code IN ('starter','premium','all');` (+ matching .down.sql) and add `"analytics.enabled":true` to the three paid seeds in seed.go (new inserts only); wire GET endpoint behind entitlement check.
  What to do / Must NOT do: must NOT store PII beyond existing ip/user_agent; must NOT add a schema change to analytics_events (columns exist, model.go:253-261); migration must be additive/backward-compatible; must NOT touch ViewCount increment logic (event_service.go:376-380 stays).
  Parallelization: Wave 3 | Blocked by: - | Blocks: 15,16
  References: `backend/internal/model/model.go:253-261` (AnalyticsEvent), `backend/internal/repository/guest_repo.go` analytics section (Create/CountByType), `backend/internal/api/server.go:205` (public routes area), `backend/internal/service/entitlement/entitlement.go:99-139` (resolver methods), `backend/internal/database/seed.go:16-68` (package seeds), `backend/migrations/` (numbering: existing 20260815000001/2), `backend/internal/api/handler/event_handler.go` (GetByID ownership pattern), event publish check in `event_service.go` PublicView.
  Acceptance criteria (agent-executable): `make test-unit` + `make test-integration` pass; tests: unauth POST records event (rate-limited), auth GET returns aggregated numbers matching inserted rows, free package → GET returns 403/feature-disabled, paid → 200.
  QA scenarios: happy — insert 3 page_views from 2 IPs + 1 whatsapp click → GET returns views=3, unique_views=2, whatsapp_clicks=1 (Evidence task-14 file); failure — POST with invalid type returns 400 (same file).
  Commit: Y | feat(backend): add analytics tracking endpoint, owner dashboard API, and paid gating
- [x] 15. Click tracking frontend: in `frontend/src/components/invitation/sections.tsx` and `share.tsx` add a `track(type)` helper (fire-and-forget `navigator.sendBeacon` or `fetch(..., {keepalive:true})` to `/api/v1/analytics/events` with `{event_id: model.eventId, type}`) invoked on: WhatsApp share click (share.tsx:69), map link clicks (sections.tsx:520,570,744-753), directions/Google Maps click (sections.tsx:744-753); phone clicks — add the same hook to any phone link if one exists on the page (none today — note in code comment); silently ignore failures (no console noise, no UI impact).
  What to do / Must NOT do: must NOT block navigation (fire-and-forget); must NOT track guestbook/RSVP interactions; must NOT add analytics to the editor/dashboard (public page only); must NOT track when `model.eventId` is missing.
  Parallelization: Wave 3 | Blocked by: 14 | Blocks: 16
  References: `frontend/src/components/invitation/share.tsx:46-74` (WA/Telegram links), `frontend/src/components/invitation/sections.tsx:520,570,744-753` (map links), `frontend/src/lib/api.ts`, model.eventId at `frontend/src/templates/types.ts:44`.
  Acceptance criteria (agent-executable): `npm run build` clean; unit test for `track()` with fetch mocked asserts correct URL+body and no-throw on failure; playwright: click WA share on a public page → network log shows POST /api/v1/analytics/events.
  QA scenarios: happy — click emits one POST (Evidence task-15 file); failure — network down: navigation still happens, no error UI (same file).
  Commit: Y | feat(frontend): track WhatsApp and map clicks on public invitations
- [x] 16. Analytics dashboard widget: create `frontend/src/components/dashboard/analytics-widget.tsx` ("use client") fetching `GET /events/{id}/analytics` (only when current subscription entitlement `analytics.enabled` — expose via existing auth/subscription context or `api.get('/users/me/subscription')` features) and rendering views/unique/WA/map/phone + RSVP count cards; mount it on `frontend/src/app/dashboard/page.tsx` (per-event card or section); free users: show only the existing ViewCount (from event list) with an upsell note linking to billing.
  What to do / Must NOT do: must NOT fetch analytics when entitlement false (404/403 tolerated); must NOT create a new route; must NOT add charting library (plain numbers/cards); must NOT show data for other users' events.
  Parallelization: Wave 3 | Blocked by: 14,15 | Blocks: -
  References: `frontend/src/app/dashboard/page.tsx`, `frontend/src/app/dashboard/billing/page.tsx` (subscription shape: package {name, code, price} at billing/page.tsx:25), `frontend/src/lib/api.ts`, `frontend/src/providers/auth-context.tsx` (auth/subscription context — verify what it exposes).
  Acceptance criteria (agent-executable): `npm run build` clean; jest (mocked api) renders numbers from GET payload; playwright with paid seed user shows widget, free user shows upsell (Evidence task-16).
  QA scenarios: happy — paid user sees 5 metrics; failure — GET returns 403 → widget hides gracefully with upsell (Evidence task-16).
  Commit: Y | feat(frontend): add gated analytics widget to dashboard
- [x] 17. ICS timezone polish: in `frontend/src/components/invitation/sections.tsx` update `icsFile` (386-415) to emit timezone-explicit DTSTART/DTEND — add `TZID` param when a named zone can be derived (event timezone stored? none today → use `Intl.DateTimeFormat().resolvedOptions().timeZone` at click time) e.g. `DTSTART;TZID=Asia/Jakarta:...`, and add `X-WR-TIMEZONE` to VCALENDAR; keep floating-time fallback when zone unavailable; update `gcalStamp` (363-369) to pass `ctz` param when available. Add unit test for ICS output shape.
  What to do / Must NOT do: must NOT change Google Calendar URL behavior when no tz; must NOT add a dependency; keep escapeIcs unchanged.
  Parallelization: Wave 3 | Blocked by: 1 | Blocks: -
  References: `frontend/src/components/invitation/sections.tsx:363-425`.
  Acceptance criteria (agent-executable): `npx jest src/components/invitation` (new ics test) green; `npm run build` clean.
  QA scenarios: happy — generated ICS contains `TZID=` and `X-WR-TIMEZONE` (Evidence task-17 file); failure — event without date produces no DTSTART (skips, no crash) (same file).
  Commit: Y | fix(frontend): make ICS export timezone-explicit
- [x] 18. Night theme tokens: extend `frontend/src/templates/types.ts` with `night?: Partial<ThemeTokens>` on TemplateDefinition; add `night` palette (dark background/surface/text/muted, dimmed primary) to all 11 template definitions; update `frontend/src/templates/theme.ts` `themeVars(theme, night?)` to emit overridden `--t-*` vars when night active (merge `theme` then `theme.night`).
  What to do / Must NOT do: must NOT touch globals.css `.dark` (dashboard palette, stays); must NOT change default (day) tokens; must NOT force dark on any template by default.
  Parallelization: Wave 4 | Blocked by: 1 | Blocks: 19,20
  References: `frontend/src/templates/types.ts:72-91` (ThemeTokens), `frontend/src/templates/theme.ts:5-29` (themeVars), template defs `frontend/src/templates/*/index.ts` (theme object).
  Acceptance criteria (agent-executable): `npm run build` clean; jest asserts themeVars(theme, night) overrides the 8 color vars and keeps day values when night omitted.
  QA scenarios: happy — night override applied; failure — no night field → identical to day output (Evidence task-18 file).
  Commit: Y | feat(frontend): add night theme token variants per template
- [x] 19. Auto night mode: create `frontend/src/hooks/use-night-mode.ts`: `useNightMode(opts: {start:"18:00", end:"06:00", timezone?: string})` returning `{night: boolean}` — computes current time in the event timezone (Intl.DateTimeFormat with timeZone; fallback browser tz), night=true when local hour ∈ [start, end) treating the window crossing midnight (18:00–06:00 → night when hour>=18 || hour<6); client-only (useEffect + useState, initial false to avoid hydration mismatch); mount in `TemplateShell` (`frontend/src/components/invitation/shell.tsx`) or LivePreview + public page wrapper: apply `themeVars(theme, night)`. Keep public [slug] page server component unchanged (client child handles it).
  What to do / Must NOT do: must NOT use `window` during render (hydration-safe); must NOT force dark in SSR; must NOT add a settings UI (config via template `darkMode` default; expose constant in types.ts `darkMode?: {auto:boolean; start:string; end:string}` on TemplateDefinition, default auto true for all templates).
  Parallelization: Wave 4 | Blocked by: 18 | Blocks: 20
  References: `frontend/src/components/invitation/shell.tsx:30` (TemplateShell, themeVars usage), `frontend/src/templates/theme.ts`, types.ts (add darkMode field), `frontend/src/app/[slug]/page.tsx:51-55`.
  Acceptance criteria (agent-executable): jest unit tests for window math (todo 20) green; `npm run build` clean.
  QA scenarios: happy — mocked clock 22:00 Asia/Jakarta → night=true; failure — 12:00 → night=false; boundary 18:00 → true, 06:00 → false (Evidence task-19 file).
  Commit: Y | feat(frontend): add auto night mode with event timezone
- [x] 20. Night mode tests: `frontend/src/hooks/__tests__/use-night-mode.test.ts` (mock Date/Intl): window boundaries, midnight crossing, timezone conversion, SSR-safety (no window access during first render).
  What to do / Must NOT do: tests only.
  Parallelization: Wave 4 | Blocked by: 18,19 | Blocks: -
  References: hook (todo 19 output).
  Acceptance criteria (agent-executable): `npx jest src/hooks/__tests__/use-night-mode.test.ts` exit 0.
  QA scenarios: happy/failure — boundaries and fallback tz cases (Evidence task-20 file).
  Commit: Y | test(frontend): cover auto night mode logic
- [x] 21. Email infra: add gomail dep (`gopkg.in/gomail.v2` or `github.com/wneessen/go-mail` — pick gomail.v2 per docs) to backend/go.mod; add config keys SMTP_HOST/PORT/USERNAME/PASSWORD/FROM/FROM_NAME to `backend/internal/config` (pattern of existing config) + `.env.example`; create `backend/internal/service/email/service.go`: `EmailService{cfg, client}` with `Send(to, subject, htmlBody)` (DialAndSend, 10s timeout) + `SendWithRetry` (3 attempts, 1m/5m/30m backoff) + HTML templates via `html/template` for welcome, payment-success, expiry-reminder, expired; log failures via zerolog, NEVER return blocking errors to callers (async wrapper `SendAsync`).
  What to do / Must NOT do: must NOT fail the caller on send error; must NOT store credentials in code (config only); must NOT add a queue dependency; must NOT send marketing email.
  Parallelization: Wave 5 | Blocked by: - | Blocks: 22,23
  References: `docs/integrations/email.md:88-133` (gomail + queue spec), `backend/go.mod`, `backend/internal/config` (existing config structs), `backend/internal/pkg/` layout.
  Acceptance criteria (agent-executable): `go build ./...` exit 0; `make test-unit` passes; unit test: SendWithRetry fails 3x → returns error after 3 attempts (with backoff shortened via injectable clock); template render test escapes user content (no HTML injection).
  QA scenarios: happy — template renders with escaped names; failure — SMTP down → SendAsync returns immediately, error logged (Evidence task-21 file).
  Commit: Y | feat(backend): add SMTP email service with retry and templates
- [x] 22. Payment success email: in `backend/internal/service/payment_service.go` after `ActivateOnSettlement` succeeds (line ~167), call `emailSvc.SendAsync` with payment-success template (plan name, amount formatted IDR, activation/expiry dates from subscription) — inject EmailService into PaymentService (server.go wiring, payment service constructor at payment_service.go:32); must run AFTER transaction state persisted; must not alter webhook response or idempotency (WebhookIdempotency logic untouched); dedupe: send only when txn transitions to settlement (guard on status change, existing logic at payment_service.go:142-167).
  What to do / Must NOT do: must NOT block/alter webhook; must NOT double-send on duplicate webhook (idempotency + status guard); must NOT send on non-settlement statuses.
  Parallelization: Wave 5 | Blocked by: 21 | Blocks: -
  References: `backend/internal/service/payment_service.go:113-180` (HandleWebhook), `:142-167` (settlement branch), `backend/internal/api/server.go:87-91` (PaymentService construction), `docs/integrations/email.md:44-51` (content spec).
  Acceptance criteria (agent-executable): `make test-unit` + `make test-integration` pass; new test: settlement webhook → email service received payload (mock sender), duplicate webhook → exactly one send; deny/cancel → no email.
  QA scenarios: happy — settlement sends once; failure — email send fails → webhook still returns 200 and subscription activated (Evidence task-22 file).
  Commit: Y | feat(backend): send payment-success email on settlement
- [x] 23. Expiry reminder ticker: in `backend/cmd/server` (or a small `internal/service/email/worker.go` started from main) add a daily goroutine ticker (interval 24h, offset configurable; run once at startup) that: scans subscriptions with `status='active'` and `expires_at` in [now+6d, now+8d] → sends expiry-reminder email; scans `expires_at < now` and status still 'active' → sends expired email + marks status 'expired' (existing status transition — reuse SubscriptionService if it has expiry handling; else do status update via repo); dedupe per subscription per day by checking `audit_log` for `action='email.expiry_reminder'` + `entity_type='subscription'` + `entity_id=<id>` created today (audit log exists: model.go:242-251, AuditLogRepo wired).
  What to do / Must NOT do: must NOT add a new table/migration; must NOT double-send (audit-log dedupe); must NOT touch user-facing subscription activation flow; ticker must be stoppable on shutdown (context cancel).
  Parallelization: Wave 5 | Blocked by: 21 | Blocks: -
  References: `backend/internal/model/model.go:56-71` (Subscription.ExpiresAt), `:242-251` (AuditLog), `backend/cmd/server/main.go` (startup wiring), `backend/internal/repository` (SubscriptionRepository methods — verify expiry query exists or add `ListExpiring/ListExpired`), `docs/integrations/email.md:62-76` (content spec).
  Acceptance criteria (agent-executable): `go build ./...` exit 0; `make test-unit` passes; unit test with mocked clock: subscription expiring in 7d → reminder email sent once (audit-log dedupe blocks second run); expired active sub → expired email + status flipped.
  QA scenarios: happy — reminder fired for the 7-day window only; failure — same-day re-run sends nothing (Evidence task-23 file).
  Commit: Y | feat(backend): add daily subscription expiry reminder ticker
- [x] 24. FAQ page: create `frontend/src/app/faq/page.tsx` reusing `faqItems` from `frontend/src/data/marketing.ts:160-189` — extract a shared `FaqList` component (`frontend/src/components/marketing/faq-list.tsx`) used by both the landing FAQ section and the new page; page has simple hero header + FaqList + CTA back to landing; add route link.
  What to do / Must NOT do: reuse existing data + one shared component; must NOT duplicate FAQ markup; must NOT change landing FAQ section behavior.
  Parallelization: Wave 6 | Blocked by: 1 | Blocks: 26
  References: `frontend/src/data/marketing.ts:160-189`, `frontend/src/components/marketing/` (existing section components pattern), landing FAQ usage (grep `faqItems`).
  Acceptance criteria (agent-executable): `npm run build` clean; playwright: navigate /faq, assert accordion toggles work and content equals marketing data.
  QA scenarios: happy — /faq renders all questions from faqItems; failure — empty data shows empty-state (not crash) (Evidence task-24 file).
  Commit: Y | feat(frontend): add FAQ page reusing marketing data
- [x] 25. Cara Order + Testimoni pages: create `frontend/src/app/cara-order/page.tsx` (steps from `howItWorks` marketing.ts:214-230, plus contact/WhatsApp CTA) and `frontend/src/app/testimoni/page.tsx` — testimonials EXPLICITLY placeholder: heading "Testimoni" + body "Kumpulan testimoni pelanggan akan segera hadir. Hubungi kami untuk berbagi pengalaman Anda." with a "WA Kami" button; include a clearly-labeled sample section (e.g. "Contoh tampilan" with `isSample: true` badges) — NO fabricated customer names/faces/quotes.
  What to do / Must NOT do: must NOT invent testimonial identities or quotes; must NOT add a testimonial data model or CMS; placeholder must be visually distinct (badge "Contoh").
  Parallelization: Wave 6 | Blocked by: 1 | Blocks: 26 | Can run with: 24
  References: `frontend/src/data/marketing.ts:214-230` (howItWorks), WhatsApp CTA pattern `frontend/src/components/marketing/` (grep `whatsapp`), landing layout `frontend/src/app/page.tsx`.
  Acceptance criteria (agent-executable): `npm run build` clean; playwright: /testimoni shows placeholder text + badge, /cara-order shows 4 steps.
  QA scenarios: happy — pages render and links work; failure — page without data does not crash (Evidence task-25 file).
  Commit: Y | feat(frontend): add cara-order and testimoni pages (placeholder testimonials)
- [x] 26. Nav/footer links + /packages fix: add /faq, /cara-order, /testimoni links to the marketing navbar and footer (components in `frontend/src/components/marketing/` and `frontend/src/app/layout.tsx` shell if used); fix the dead link in `frontend/src/app/dashboard/billing/page.tsx:161` (currently href `/packages` which does not exist) → point to landing `/#paket` anchor; verify all internal links resolve (grep for href="/packages" everywhere).
  What to do / Must NOT do: only link fixes and additions; must NOT create a /packages page; must NOT change nav behavior on dashboard.
  Parallelization: Wave 6 | Blocked by: 24,25 | Blocks: 27
  References: `frontend/src/components/marketing/navbar.tsx`, `frontend/src/components/marketing/footer.tsx` (verify paths), `frontend/src/app/dashboard/billing/page.tsx:161`, landing sections anchors (`#paket` in `frontend/src/app/page.tsx`).
  Acceptance criteria (agent-executable): `npm run build` clean; grep: no `href="/packages"` remains; playwright: click each new nav link lands on the right page.
  QA scenarios: happy — billing page link goes to landing #paket; failure — a broken internal link is caught by the grep gate (Evidence task-26 file).
  Commit: Y | fix(frontend): wire trust-page links and fix /packages dead link
- [x] 27. Trust pages tests: `frontend/src/app/__tests__/trust-pages.test.tsx` (or per-page test files) — render /faq, /cara-order, /testimoni and assert key content present, placeholder testimonials have no fabricated identity (no person names in testimonials data), nav/footer include the new links.
  What to do / Must NOT do: tests only; must NOT mock page internals beyond router/nav.
  Parallelization: Wave 6 | Blocked by: 24,25,26 | Blocks: -
  References: pages from todos 24-26.
  Acceptance criteria (agent-executable): `npx jest src/app/__tests__/trust-pages` exit 0.
  QA scenarios: happy — pages render expected headings; failure — placeholder testimonial with a name fails the no-fabricated-identity assertion (Evidence task-27 file).
  Commit: Y | test(frontend): cover trust pages and link integrity
- [x] 28. Brand normalization: normalize metadata/copy to Indonesian: `frontend/src/app/layout.tsx:4-7` (title/description/metadata → Indonesian, brand stays "Owndangan"); audit `frontend/src/data/marketing.ts`, navbar/footer, dashboard/auth shells for English strings → Indonesian (keep brand, technical labels, and payment provider names); keep `lang="id"` in html tag.
  What to do / Must NOT do: display copy only; must NOT change routes/URLs, API contracts, or template content; must NOT rename brand.
  Parallelization: Wave 7 | Blocked by: 1 | Blocks: - | Can run with: 29
  References: `frontend/src/app/layout.tsx:4-7`, `frontend/src/data/marketing.ts`, `frontend/src/app/dashboard/layout.tsx`, auth pages under `frontend/src/app/(auth)/`.
  Acceptance criteria (agent-executable): `npm run build` clean; grep audit report: no remaining user-facing English in the audited shells (allow list: brand, "Premium", provider names).
  QA scenarios: happy — metadata is Indonesian in built HTML; failure — a hardcoded English string remains → grep gate flags it (Evidence task-28 file).
  Commit: Y | fix(frontend): normalize UI copy to Indonesian
- [x] 29. Pricing display sync: change frontend pricing display (`frontend/src/data/marketing.ts:108-158`) to canonical backend values — Starter Rp99.000/30 hari, Premium Rp299.000/60 hari, All Access Rp999.000/lifetime — add comment `// Source of truth: backend/internal/database/seed.go`; update docs pricing ranges to match: `docs/prd.md:55`, `docs/terminology.md:93`, `docs/project-context.md:29`; verify billing page (`frontend/src/app/dashboard/billing/page.tsx`) shows the same values.
  What to do / Must NOT do: DISPLAY ONLY — must NOT change backend seed VALUES or `packages` table data (user-approved Q1); must NOT add new pricing fields; keep 3 plans.
  Parallelization: Wave 7 | Blocked by: 1 | Blocks: 31 | Can run with: 28
  References: `backend/internal/database/seed.go:29-68` (canonical), `frontend/src/data/marketing.ts:108-158`, `frontend/src/app/dashboard/billing/page.tsx`, `docs/prd.md:55`, `docs/terminology.md:93`, `docs/project-context.md:29`.
  Acceptance criteria (agent-executable): grep: marketing.ts prices equal seed.go prices; `npm run build` clean; docs mention only canonical ranges.
  QA scenarios: happy — landing shows 99rb/299rb/999rb; failure — a mismatched price is caught by the grep comparison (Evidence task-29 file).
  Commit: Y | fix(frontend,docs): sync pricing display to backend canonical values
- [x] 30. Onboarding auto-login: in `frontend/src/app/(auth)/register/page.tsx` after successful register, consume the tokens already returned by the API (`backend/internal/api/handler/auth_handler.go:25-37` returns token pair) — store via the existing auth/session storage (same mechanism login uses; grep how login page stores tokens) and redirect to `/dashboard` instead of linking to /login (page.tsx:181); handle token-missing fallback (redirect to /login with a notice).
  What to do / Must NOT do: must NOT add a new auth endpoint or change backend; must NOT call a second login request; must reuse the existing token-storage helper.
  Parallelization: Wave 7 | Blocked by: 1 | Blocks: 31 | Can run with: 28,29
  References: `frontend/src/app/(auth)/register/page.tsx` (form + current success path, :181 link), `frontend/src/app/(auth)/login/page.tsx` (token storage pattern), `frontend/src/lib/` (api/auth client — grep `setToken`/`localStorage`), `backend/internal/api/handler/auth_handler.go:25-37`.
  Acceptance criteria (agent-executable): `npm run build` clean; playwright: register a new user (or mock API) → lands on /dashboard with authenticated session, no manual login.
  QA scenarios: happy — fresh register goes straight to dashboard; failure — API returns no tokens → redirects to /login with notice (Evidence task-30 file).
  Commit: Y | feat(frontend): auto-login after registration
- [x] 31. Pricing+onboarding tests: `frontend/src/data/__tests__/pricing.test.ts` (prices match `seed.go` canonical constants — hardcode the 3 canonical values in the test with the same Source-of-truth comment; fail if changed without updating both) and `frontend/src/app/(auth)/__tests__/register.test.tsx` (mock register API: success → token stored + redirected to /dashboard; token-less response → redirects to /login).
  What to do / Must NOT do: tests only; must NOT mock the pricing data (test the real marketing.ts export).
  Parallelization: Wave 7 | Blocked by: 29,30 | Blocks: -
  References: outputs of todos 29-30.
  Acceptance criteria (agent-executable): `npx jest src/data/__tests__/pricing.test.ts src/app/(auth)/__tests__/register.test.tsx` exit 0.
  QA scenarios: happy — canonical prices pass; failure — price drift (e.g. Starter 120000) fails the test (Evidence task-31 file).
  Commit: Y | test(frontend): lock pricing display and onboarding flow
- [x] 32. Docs split: restructure `docs/` — move dev-only/internal docs (e.g. `docs/integrations/`, architecture notes) under `docs/internal/` via `git mv`; keep public-facing docs (product/PRD, terminology, project context) at `docs/` root; update `docs/README.md` index and any cross-references; document the split rule ("internal" = implementation details, "root" = product/architecture overview).
  What to do / Must NOT do: must NOT delete doc content; must NOT move docs that public readers need at root (prd, terminology, project-context stay); update all internal cross-links found by grep.
  Parallelization: Wave 8 | Blocked by: 29 | Blocks: 34 | Can run with: 33
  References: `docs/` tree (list first), `docs/README.md`, `docs/integrations/email.md` (moves under internal), grep for `docs/integrations` references across repo.
  Acceptance criteria (agent-executable): `find docs -type f` shows the new layout; grep: no stale references to moved paths; `docs/README.md` index matches reality.
  QA scenarios: happy — all links resolve after move; failure — a stale path in README or code comment is caught by grep (Evidence task-32 file).
  Commit: Y | docs: split public and internal documentation
- [x] 33. Repo hygiene: delete untracked 0-byte files `-l` (root) and `backend/-l`; add `-l` and `*.tsbuildinfo` to `.gitignore`; `git rm --cached frontend/tsconfig.tsbuildinfo` (keep the file locally via .gitignore); verify `git status` clean of these artifacts.
  What to do / Must NOT do: must NOT delete any real source file; must NOT touch the 30+ WIP modified files; git rm --cached only for tsbuildinfo.
  Parallelization: Wave 8 | Blocked by: - | Blocks: 34 | Can run with: 32
  References: `git status` (verify paths), `.gitignore` (root), `frontend/tsconfig.tsbuildinfo` (tracked).
  Acceptance criteria (agent-executable): `git status --porcelain` shows no `-l` or `tsbuildinfo` entries; `git ls-files | grep -c tsbuildinfo` = 0.
  QA scenarios: happy — artifacts gone from git status; failure — re-running `git status` after a fresh build shows no tsbuildinfo (Evidence task-33 file).
  Commit: Y | chore: remove stray -l files and untrack tsbuildinfo
- [x] 34. Final gates: run `make test` (backend), `npm run build`, `npm run lint` (frontend) — classify EVERY failure as NEW (caused by this plan's changes) vs PRE-EXISTING (present on the initial commit/WIP baseline) by re-running on a stashed baseline if needed; all NEW failures must be fixed; PRE-EXISTING ones documented with evidence (baseline output) and left untouched; produce a gate report.
  What to do / Must NOT do: must NOT fix pre-existing failures unrelated to the plan; must NOT delete/modify tests to pass; must record baseline evidence before claiming pre-existing.
  Parallelization: Wave 8 | Blocked by: 32,33 + all todos | Blocks: F1-F4
  References: Makefile (backend targets), `frontend/package.json` scripts; git stash/`git stash -u` for baseline if needed (careful with WIP — prefer a worktree or `git diff` comparison instead of stash).
  Acceptance criteria (agent-executable): gate report lists each command + exit code + NEW/PRE-EXISTING classification + evidence path; 0 NEW failures.
  QA scenarios: happy — all gates green, report written to `.omo/evidence/owndangan-feature-wave/task-34-owndangan-feature-wave.md`; failure — a NEW failure is fixed and re-run green (Evidence task-34 file).
  Commit: Y | chore: final gate verification report

## Final verification wave
> Runs in parallel after ALL todos. ALL must APPROVE. Surface results and wait for the user's explicit okay before declaring complete.
- [x] F1. Plan compliance audit: for every todo 1-34, verify the implemented code matches the todo's What/Must NOT/Acceptance criteria and references; check no Must NOT have guardrail was violated (esp. pricing data untouched, no fabricated testimonials, no new queue/table, no schema change for occasions, no subscription activation from frontend). Evidence: `.omo/evidence/owndangan-feature-wave/final-F1/`.
  Fail if: any todo's acceptance criteria unmet, or any guardrail violated.
  Approve only when: all 34 todos verified against their written criteria.
- [x] F2. Code quality review: review the full diff for this plan's changes (backend Go + frontend TS) against AGENTS.md rules — thin handlers, services/repos layering, server-side authorization (analytics ownership, guest restore ownership), input validation, no secrets, no `as any`/`@ts-ignore`, no empty catch blocks, error handling that never silently eats failures (async email failures are logged by design). Evidence: `.omo/evidence/owndangan-feature-wave/final-F2/`.
  Fail if: layering violations, auth bypass, unvalidated input at a trust boundary, or suppressed type errors.
  Approve only when: no blocking findings.
- [x] F3. Real manual QA: browser-driven QA (playwright, no API mocking) against a running stack — walk the user journeys: editor live preview updates without network calls; template occasion filter; guest trash → restore; analytics widget shows clicks/views after simulated clicks on a published invitation; dark-mode auto-switch (clock mocking or manual); payment-success email (SMTP test sink) and expiry reminder (clock skew); /faq + /cara-order + /testimoni render; onboarding auto-login; pricing display matches seed. Screenshots per journey in `.omo/evidence/owndangan-feature-wave/final-F3/`.
  Fail if: any core journey broken or visually broken (screenshots show layout breakage).
  Approve only when: all journeys pass with screenshots.
- [x] F4. Scope fidelity: diff the delivered change-set against Scope IN/OUT — no out-of-scope feature shipped, no scope item missed, no unrelated WIP files touched (verify via `git status` + diff against the initial commit 7d31e8c); confirm the 16-section brief (master-prompt.md) is fully addressed section by section. Evidence: `.omo/evidence/owndangan-feature-wave/final-F4/`.
  Fail if: out-of-scope changes or missed brief sections.
  Approve only when: scope matches the plan exactly.

## Commit strategy
- One commit per todo (see `Commit:` line on each todo), conventional commits, scope-prefixed — follow repo style (`git log --oneline -10` to confirm; repo uses conventional style).
- Commits stay in the working tree (no push, no PR) unless the user explicitly asks — the repo is a single-commit WIP tree (7d31e8c + 30+ modified, ~15 untracked); do NOT commit unrelated WIP files. Stage only the files each todo touched.
- Order: commit in todo order per wave; wave-level commits optional if a wave's todos are too interleaved to separate cleanly (prefer per-todo).
- Final verification wave produces NO code commits (audit/report evidence only, under `.omo/evidence/`).

## Success criteria
The plan is DONE when ALL of these hold:
- All todos 1-34 checked off with evidence files; final wave F1-F4 all approved by the user.
- `make test`, `npm run build`, `npm run lint` exit 0 with zero NEW failures (pre-existing ones documented).
- Live preview renders the real public template renderer live in the editor (no save-to-preview roundtrip, no per-keystroke API calls) — verified by network-log assertion.
- All 11 templates carry occasion metadata; corporate template + seed row live and selectable in editor and showcase; filter works.
- Guest trash/restore round-trip works; analytics captures WA/map/phone clicks and shows a paid-only dashboard widget with unique-view counts.
- Night mode auto-applies 18:00-06:00 per event timezone on public page and preview, hydration-safe.
- Payment-success email fires once per settled transaction (webhook unaffected); expiry reminder fires at 7 days with audit-log dedupe; SMTP config documented in .env.example.
- /faq, /cara-order, /testimoni (placeholder, clearly marked) live and linked from nav/footer; /packages dead link fixed.
- All user-facing copy Indonesian; pricing display matches backend seed (99rb/30hr, 299rb/60hr, 999rb) with source-of-truth comment; docs pricing ranges corrected.
- Register auto-logs-in to /dashboard.
- docs/ split public/internal; `-l` files deleted; tsbuildinfo untracked; `git status` free of these artifacts.
