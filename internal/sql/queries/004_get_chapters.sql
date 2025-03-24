-- name: GetChapters :many
SELECT * FROM book_chapters WHERE book_id = ?;
