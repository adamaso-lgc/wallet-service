-- +goose Up
-- +goose StatementBegin

CREATE TABLE events (
    id            UUID        NOT NULL DEFAULT gen_random_uuid(),
    aggregate_id  UUID        NOT NULL,
    event_type    TEXT        NOT NULL,
    payload       JSONB       NOT NULL,
    version       BIGINT      NOT NULL,
    occurred_at   TIMESTAMPTZ NOT NULL,

    PRIMARY KEY (id),
    -- Optimistic concurrency: two concurrent writers cannot both commit
    -- the same version for the same aggregate.
    UNIQUE (aggregate_id, version)
);

CREATE INDEX idx_events_aggregate_id ON events (aggregate_id);

CREATE TABLE wallet_views (
    id         UUID           NOT NULL,
    owner_id   TEXT           NOT NULL,
    balance    NUMERIC(20, 8) NOT NULL,
    currency   TEXT           NOT NULL,
    status     TEXT           NOT NULL,
    created_at TIMESTAMPTZ    NOT NULL,
    updated_at TIMESTAMPTZ    NOT NULL,

    PRIMARY KEY (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallet_views;
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
