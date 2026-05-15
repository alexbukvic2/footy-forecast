# 0002 — No ORM

Status: accepted
Date: 2026-05-15

## Context
Go ORMs (GORM, ent, etc.) hide SQL, generate inefficient queries, and complicate
debugging. The schema for footy-forecast is small and stable.

## Decision
Use `pgx/v5` directly. SQL queries live in repository files as plain strings or
in `.sql` files compiled by `sqlc` if we adopt it later.

## Consequences
- SQL is visible and reviewable.
- Slightly more boilerplate for scans — acceptable.
- Re-evaluate sqlc once we have 10+ tables.
