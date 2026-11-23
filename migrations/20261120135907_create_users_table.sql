-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    username VARCHAR(20) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS users;