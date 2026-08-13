# Git Workflow

## Branching Strategy

We follow **GitHub Flow** (short-lived feature branches, no long-running develop branch).

```
main  ───●─────────●──────────────●─────────●──
          \         /              \         /
feature    ●──●──●                   ●──●──●
```

### Branch Naming Convention

```
<type>/<short-description>
```

| Type | Purpose | Example |
|------|---------|---------|
| `feat/` | New feature | `feat/invitation-rsvp` |
| `fix/` | Bug fix | `fix/payment-amount-rounding` |
| `chore/` | Tooling, CI, dependencies | `chore/upgrade-go-1.22` |
| `refactor/` | Code restructuring | `refactor/payment-service` |
| `docs/` | Documentation | `docs/api-endpoints` |
| `test/` | Adding/fixing tests | `test/payment-webhook` |

**Rules**:
- Use kebab-case (e.g., `feat/invitation-rsvp` not `feat/InvitationRSVP`).
- Keep names short but descriptive (under 50 characters).
- Delete branch after merging.

## Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat` — New feature.
- `fix` — Bug fix.
- `chore` — Maintenance, deps, CI.
- `refactor` — Code change with no behavior change.
- `test` — Adding or fixing tests.
- `docs` — Documentation only.
- `perf` — Performance improvement.
- `style` — Formatting, linting (no logic change).

### Examples

```
feat(payment): add Midtrans Snap token generation
fix(auth): handle expired refresh token gracefully
chore(deps): upgrade testify to v1.9.0
docs(readme): add setup instructions
```

### Rules

- Limit subject to 72 characters.
- Use imperative mood ("add" not "added" or "adds").
- Do not capitalize the subject.
- Do not end the subject with a period.
- Reference issues in the footer: `Closes #42`.

## PR Process

1. Create branch from `main`.
2. Make changes, commit with conventional commits.
3. Push branch and open a PR against `main`.
4. PR title follows the commit convention format.
5. Add a description explaining what and why (not how).
6. Request review from at least one teammate.
7. Address review feedback with additional commits.
8. Squash-merge into `main` when approved.

### PR Template

```markdown
## Description
Briefly describe the change and why it's needed.

## Type
- [ ] feat
- [ ] fix
- [ ] chore
- [ ] refactor
- [ ] test
- [ ] docs

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing done

## Checklist
- [ ] Code follows project conventions
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No new warnings
```

## Code Review Guidelines

- **Review within 24 hours** on business days.
- Approve only when confident the code is correct.
- Leave constructive comments (explain *why* something should change).
- Distinguish between blockers (must fix) and suggestions (nice to have).
- Use GitHub's "suggest changes" feature for small fixes.