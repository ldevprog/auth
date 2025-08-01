-- +goose Up
CREATE TABLE refresh_tokens (
    token text PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamp DEFAULT now(),
    expires_at timestamp NOT NULL,
    revoked boolean DEFAULT FALSE
);

-- +goose Down
DROP TABLE refresh_tokens;

