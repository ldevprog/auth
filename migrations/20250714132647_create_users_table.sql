-- +goose Up
CREATE TABLE users (
    id bigint PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL,
    role smallint NOT NULL DEFAULT 0,
    password text NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp
);

-- +goose Down
DROP TABLE users;

