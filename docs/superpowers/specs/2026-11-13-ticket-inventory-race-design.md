# Ticket Inventory Race-Condition Design

## Goal

Prevent ticket overselling when concurrent requests try to hold the last available
inventory, while returning an error that is clear in logs and API responses.

## Chosen approach

Use PostgreSQL pessimistic row locking for inventory mutations. The transaction
will lock the matching `ticket_inventories` row with `FOR UPDATE`, read the
latest `total_quota`, `sold_count`, and `held_count`, and update `held_count`
only when enough inventory remains. A concurrent request waits for the first
transaction to finish, then evaluates the fresh counts and fails if capacity is
gone.

This follows the existing transaction pattern in checkout and avoids adding a
new dependency or changing the successful API response contract. The existing
`version` column remains useful for tracking mutations but will not be used as
the primary concurrency mechanism in this change.

## Error behavior

Inventory exhaustion will use a stable, readable error message:
`ticket inventory is not available`.

Repository errors will retain useful context when returned from the transaction,
without swallowing database errors. Missing inventory rows and database failures
remain distinguishable from an inventory-capacity failure.

## Scope and data flow

The change is limited to inventory-holding/checkout repository paths and their
focused tests. Inventory availability is checked and reserved in one database
transaction. The booking flow continues to create booking items from active
carts; payment success remains responsible for moving quantities from held to
sold through the existing database trigger.

## Verification

Add a concurrency-focused test that starts two requests against a one-ticket
inventory and verifies exactly one succeeds, the other returns the readable
availability error, and the final inventory never exceeds its quota. Run the
targeted Go tests and the Go race detector for the affected packages.
