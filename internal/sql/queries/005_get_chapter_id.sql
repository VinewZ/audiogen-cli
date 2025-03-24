-- name: GetChapterID :one
SELECT id FROM book_chapters WHERE chapter = ?;
