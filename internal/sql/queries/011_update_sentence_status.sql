-- name: UpdateSentenceStatus :exec
UPDATE book_sentences SET status = ? WHERE id = ? AND status IN ('pending', 'done');
