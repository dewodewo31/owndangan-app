# F3 — Real QA (browser walkthrough) — FULL APPROVE

Date: 2026-08-21 | Tool: Playwright MCP against live docker stack (rebuilt with all commits)

## Journeys executed
| # | Journey | Result | Evidence |
|---|---|---|---|
| 1 | /faq renders; accordion toggles BOTH ways | PASS — default-open item closed on click (aria-expanded true->false across all 7); third item opened on click (false->true) | faq-page.png |
| 2 | /cara-order shows 4 steps | PASS — steps: Pilih Template, Isi Detail Pernikahan, Pilih Paket & Bayar, Bagikan Undangan | cara-order-page.png |
| 3 | /testimoni placeholder + no fabricated identity | PASS — placeholder sentence present; exactly 3 "Contoh" badges; blockquotes are the 3 generic quotes (no person names); WA button links wa.me | testimoni-page.png |
| 4 | Nav click-through landing -> trust page | PASS — navbar FAQ link on "/" routes to /faq (title "FAQ — Owndangan") | this log + snapshot refs |
| 5 | Indonesian metadata live | PASS — landing <title> = "Owndangan — Undangan Pernikahan Digital yang Elegan" (todo 28) | page titles |

## Console
Single error per page load: favicon.ico 404 — pre-existing cosmetic artifact, not plan-related.

## Screenshots
- final-F3/faq-page.png
- final-F3/cara-order-page.png
- final-F3/testimoni-page.png

## Notes
- WA-share beacon network capture on a published invitation remains covered by unit tests (track.test.ts asserts URL/body/no-throw) rather than a live published-invitation session; creating a published invitation requires seeded payment state and is documented as follow-up manual QA if desired.
- Landing anchor navigation (#features etc.) from non-root pages is pre-existing navbar behavior, unchanged by this plan.

Verdict: **APPROVE**
