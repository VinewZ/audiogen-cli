-- +goose Up
CREATE TABLE book_text (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    book_id TEXT NOT NULL,
    page INTEGER NOT NULL,
    chapter TEXT NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE book_text;
