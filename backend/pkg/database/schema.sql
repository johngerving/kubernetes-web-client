CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL
);

CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

CREATE TABLE workspaces (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    owner INT REFERENCES users (id) NOT NULL,
    PRIMARY KEY (owner, name)
);