-- name: InsertText :exec
INSERT INTO book_text (book_id, page, chapter, content) 
VALUES (?, ?, ?, ?);
