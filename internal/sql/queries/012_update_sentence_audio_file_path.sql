-- name: UpdateAudioFilePath :exec
UPDATE book_sentences SET audio_file_path = ? WHERE id = ?;
