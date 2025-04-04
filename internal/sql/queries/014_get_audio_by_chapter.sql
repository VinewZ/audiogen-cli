-- name: GetAudioWhereChapter :many
SELECT * FROM book_sentences WHERE chapter_id = ? ORDER BY id;
