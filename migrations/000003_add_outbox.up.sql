CREATE TABLE outbox (
    id serial PRIMARY KEY,
    topic text NOT NULL,
    payload bytea NOT NULL,
    published boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_outbox_unpublished ON outbox (published) WHERE published = false;