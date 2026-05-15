# 0001 — Project layout

Status: accepted
Date: 2026-05-15

## Context
We need a clear directory structure that scales from a single-server pet
project to a real product without restructuring later.

## Decision
Adopt the standard layered Go layout:
- `cmd/server` for the binary entrypoint
- `internal/` for all application code (Go enforces module-private)
- Layers: `http → service → repository → domain`
- `migrations/`, `scripts/`, `docs/` at root

## Consequences
- Clear import direction enforced by package boundaries.
- Slightly more files than a flat layout, worth it for clarity.
- `internal/` prevents accidental external imports if we ever publish modules.
