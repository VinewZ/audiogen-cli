-- name: InsertChapter :exec
INSERT INTO book_chapters (book_id, start_page, chapter) 
VALUES (?, ?, ?);
