# Architecture Decision Records (ADRs)

This directory contains Architecture Decision Records (ADRs) for the Owndangan platform. ADRs document significant architectural decisions, including the context, options considered, and the final decision.

## What is an ADR?

An Architecture Decision Record is a short document that captures:
- **Context**: Why the decision was needed.
- **Options**: What alternatives were considered.
- **Decision**: What was chosen and why.
- **Consequences**: What trade-offs were accepted.

## ADR Index

| ID | Title | Status | Date |
|----|-------|--------|------|
| [ADR-001](./ADR-001-initial-architecture.md) | Initial Architecture Decision | Accepted | 2025-01-15 |

## Status Meanings

| Status | Description |
|--------|-------------|
| **Proposed** | Under discussion, not yet approved. |
| **Accepted** | Approved and implemented. |
| **Deprecated** | Superseded by a newer ADR. |
| **Superseded** | Replaced by a newer ADR. |

## Process for New ADRs

1. Create a new ADR file following the naming convention: `ADR-<NNN>-<short-description>.md`.
2. Use the template below.
3. Submit as a PR with the label `adr`.
4. Discuss and iterate with the team.
5. Once consensus is reached, update status to `Accepted` and merge.

### Template

```markdown
# ADR-<NNN>: <Title>

- **Status**: Proposed
- **Date**: <YYYY-MM-DD>
- **Author**: <Name>

## Context

[Describe the problem or opportunity that prompted this decision.]

## Options

### Option 1: <Name>
[Description, pros, cons]

### Option 2: <Name>
[Description, pros, cons]

## Decision

We chose **Option <N>** because [reasoning].

## Consequences

- [Positive consequence]
- [Negative consequence / trade-off]
- [Migration plan if applicable]
```

## References

- [Documenting Architecture Decisions](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions) by Michael Nygard.
- [ADR GitHub organization](https://adr.github.io/).