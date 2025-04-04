package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/vinewz/audiogen-cli/sqlc"
)

func Narrate(ctx context.Context, tracker *progress.Tracker, db *sql.DB, title, lang string) error {
	queries := sqlc.New(db)

	bookId, err := queries.GetBookID(ctx, sql.NullString{String: title, Valid: true})
	if err != nil {
		return fmt.Errorf("Error getting BookId: %v", err)
	}

	stcs, err := queries.GetPendingSentences(ctx, bookId)
	if err != nil {
		return fmt.Errorf("Error getting pending sentences: %v", err)
	}

	tracker.Total = int64(len(stcs))

	for idx, stc := range stcs {
		formData := url.Values{
			"text_input":            {strings.ReplaceAll(stc.Sentence, ".", "\n")},
			"text_filtering":        {"standard"},
			"character_voice_gen":   {"female_01.wav"},
			"narrator_enabled":      {"false"},
			"narrator_voice_gen":    {"male_01.wav"},
			"text_not_inside":       {"character"},
			"language":              {lang},
			"output_file_name":      {fmt.Sprintf("%s_%04d", title, idx)},
			"output_file_timestamp": {"false"},
			"autoplay":              {"false"},
			"autoplay_volume":       {"0.1"},
		}

		postSent, err := http.PostForm(
			"http://127.0.0.1:7851/api/tts-generate",
			formData,
		)
		if err != nil {
			return fmt.Errorf("Error posting to API: %v", err)
		}
		defer postSent.Body.Close()

		bd, err := io.ReadAll(postSent.Body)
		if err != nil {
			return fmt.Errorf("Error reading response body: %v", err)
		}

		type TTSResponse struct {
			Status         string `json:"status"`
			OutputFilePath string `json:"output_file_path"`
			OutputFileURL  string `json:"output_file_url"`
			OutputCacheURL string `json:"output_cache_url"`
		}

		var resJson TTSResponse
		err = json.Unmarshal(bd, &resJson)
		if err != nil {
			return fmt.Errorf("Error unmarshaling JSON: %v", err)
		}

		err = queries.UpdateAudioFilePath(
			ctx,
			sqlc.UpdateAudioFilePathParams{
				ID:            stc.ID,
				AudioFilePath: sql.NullString{String: resJson.OutputFilePath, Valid: true},
			},
		)
		if err != nil {
			return fmt.Errorf("Error updating audio file path status:\n---\n%s\n---\n%v", stc.Sentence, err)
		}

		err = queries.UpdateSentenceStatus(
			ctx,
			sqlc.UpdateSentenceStatusParams{
				ID:     stc.ID,
				Status: sql.NullString{String: "done", Valid: true},
			},
		)
		if err != nil {
			return fmt.Errorf("Error updating sentence status:\n---\n%s\n---\n%v", stc.Sentence, err)
		}
		tracker.Increment(1)
	}

	return nil
}
