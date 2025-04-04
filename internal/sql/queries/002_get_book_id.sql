-- name: GetBookID :one
SELECT id FROM books WHERE title = ?;
