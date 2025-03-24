-- name: InsertBook :exec
INSERT INTO books (id, title, total_pages) 
VALUES (?, ?, ?);
