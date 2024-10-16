CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email text NOT NULL
);

CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);