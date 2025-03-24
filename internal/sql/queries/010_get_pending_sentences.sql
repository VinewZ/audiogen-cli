-- name: GetPendingSentences :many
SELECT * FROM book_sentences WHERE book_id = ? AND status = 'pending';
