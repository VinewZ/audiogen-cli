-- name: InsertSentence :exec
INSERT INTO book_sentences (book_id, chapter_id, sentence, audio_file_path) 
VALUES (?, ?, ?, ?);
