-- +goose Up
CREATE TABLE feed (
    id uuid PRIMARY KEY,
    name text not null,
    url text UNIQUE not null,
    user_id uuid not null REFERENCES users(id) ON DELETE CASCADE,
    last_fetched_at timestamp
);

-- +goose Down
DROP TABLE feed;