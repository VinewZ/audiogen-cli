-- +goose Up
CREATE TABLE book_chapters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    book_id TEXT NOT NULL,
    start_page INTEGER NOT NULL,
    chapter TEXT NOT NULL,
    FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE,
    UNIQUE(book_id, chapter)
);

-- +goose Down
DROP TABLE book_chapters;
