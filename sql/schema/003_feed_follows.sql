-- +goose Up
CREATE TABLE feed_follows (
    id uuid PRIMARY KEY,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id uuid NOT NULL REFERENCES users(id) on DELETE CASCADE,
    feed_id uuid NOT NULL REFERENCES feed(id) on DELETE CASCADE
);

-- +goose Down
DROP TABLE feed_follows;