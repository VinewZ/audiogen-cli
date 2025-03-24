-- name: GetSentences :many
SELECT * FROM book_sentences WHERE book_id = ?;
