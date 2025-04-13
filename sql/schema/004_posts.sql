-- +goose Up
create table posts (
id uuid primary KEY,
created_at timestamp not null,
updated_at timestamp not null,
title text not null, 
url text not null,
description text,
published_at text,
feed_id uuid not null
);

-- +goose Down
drop table posts;