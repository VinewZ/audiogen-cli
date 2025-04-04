-- name: GetAllAudios :many
SELECT * FROM book_sentences WHERE book_id = ? ORDER BY id;
