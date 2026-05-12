-- name: InsertEvents :exec
-- Inserts a batch of events for one or more aggregates.
-- The UNIQUE(aggregate_id, version) constraint enforces optimistic concurrency:
-- if another writer already committed the same version, this will error.
INSERT INTO events (aggregate_id, event_type, payload, version, occurred_at)
SELECT
    unnest(@aggregate_ids::uuid[]),
    unnest(@event_types::text[]),
    unnest(@payloads::jsonb[]),
    unnest(@versions::bigint[]),
    unnest(@occurred_ats::timestamptz[]);

-- name: LoadEventsByAggregateID :many
SELECT id, aggregate_id, event_type, payload, version, occurred_at
FROM   events
WHERE  aggregate_id = @aggregate_id
ORDER  BY version ASC;
