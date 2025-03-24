-- +goose Up
CREATE TABLE books (
    id TEXT PRIMARY KEY,
    title TEXT UNIQUE,
    total_pages INTEGER NOT NULL
);

-- +goose Down
DROP TABLE books;
