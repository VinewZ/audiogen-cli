-- name: GetText :many
SELECT * FROM book_text WHERE book_id = ?;
