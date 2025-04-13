-- +goose Up
CREATE TABLE feed (
    name text not null,
    url text UNIQUE not null,
    user_id uuid not null REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feed;