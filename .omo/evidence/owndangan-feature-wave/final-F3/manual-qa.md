# F3 — Real QA (smoke level PASSED; browser walkthrough PENDING)

## Live-stack smoke (executed 2026-08-21, docker stack rebuilt with ALL commits)
- GET / -> 200
- GET /faq -> 200
- GET /cara-order -> 200
- GET /testimoni -> 200
- Backend GET /api/v1/packages -> 200
- POST /api/v1/analytics/events route live (validation path responds)
- Initial rebuild exposed stale-container 404s for new pages — caught and fixed by rebuild (validates this gate's purpose)

## Automated coverage standing in for journeys (all green)
- Editor live preview renders real renderer without API calls -> unit/build gates (todo 7 evidence)
- Template occasion filter -> build + occasions.test.ts (todo 11)
- Guest trash -> restore -> backend integration tests (todo 12) + UI committed (todo 13)
- Analytics widget numbers/gating/upsell -> mocked-API jest suite (todo 16)
- Night mode window math incl. midnight wrap -> hook tests (todo 20)
- Payment-success email + duplicate-webhook dedupe -> handler tests (todo 22)
- Expiry reminder dedupe/expired flip -> clock-injected tests (todo 23)
- Trust pages content + no fabricated identity -> trust-pages tests (todo 27)
- Pricing lock vs seed -> pricing.test.ts (todo 31)
- Register auto-login success/fallback -> register tests (todo 31)

## PENDING (requires interactive browser pass with screenshots)
- [ ] Visual click-through: WA share click fires beacon (network tab) on a published invitation
- [ ] Dark-mode visual switch at mocked evening hour
- [ ] Accordion toggle visuals on /faq; Contoh badge styling on /testimoni
- [ ] Nav link click-through from landing to each trust page
- [ ] Full register->dashboard journey against seeded backend in browser

Verdict: **PARTIAL APPROVE** — automated+HTTP coverage green; screenshots/click-through journeys await a browser session (playwright) or human walkthrough.
