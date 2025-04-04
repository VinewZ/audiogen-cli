-- +goose Up
CREATE TABLE book_sentences (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    book_id TEXT NOT NULL,
    chapter_id INTEGER NOT NULL,
    sentence TEXT NOT NULL,
    audio_file_path TEXT,
    status TEXT DEFAULT 'pending',
    FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE,
    FOREIGN KEY (chapter_id) REFERENCES book_chapters (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE book_sentences;
