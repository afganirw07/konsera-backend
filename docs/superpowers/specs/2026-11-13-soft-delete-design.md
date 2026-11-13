# Soft Delete Design

## Goal

Ensure business data is never physically deleted. Delete endpoints must mark
rows with `deleted_at`, and normal reads must exclude rows where
`deleted_at IS NOT NULL`.

## Scope

- Add `deleted_at TIMESTAMPTZ` to business tables that do not have it.
- Convert repository-level business `DELETE FROM` statements to
  `UPDATE ... SET deleted_at = NOW()` with active-row guards.
- Add active-row filters to reads and related existence checks.
- Replace business foreign-key `ON DELETE CASCADE` behavior so deleting a
  parent cannot physically remove child records.
- Keep non-business in-memory cleanup (for example rate-limit cache entries)
  unchanged.
- Preserve existing HTTP routes, response codes, authorization, and API
  shapes.

## Data flow

1. A DELETE endpoint validates the request through the existing handler and
   service layers.
2. The repository updates `deleted_at` only when the target row is active.
3. Subsequent public reads and mutation lookups ignore soft-deleted rows.
4. Historical rows remain available to database administrators and future
   restore/audit workflows.

## Error handling

Repositories continue returning their existing not-found and database errors.
A repeated delete behaves as not found because the active-row guard affects zero
rows. No broad catches or success fallbacks are introduced.

## Testing

Add focused repository tests for representative parent and child entities:

- delete updates `deleted_at` instead of issuing a physical delete;
- active reads exclude soft-deleted rows;
- repeated deletes do not affect already deleted rows.

Run `gofmt` and the repository's Go test suite after implementation.

## Non-goals

- No endpoint renaming or response contract changes.
- No purge job or restore endpoint.
- No changes to unrelated cache/map deletion.
