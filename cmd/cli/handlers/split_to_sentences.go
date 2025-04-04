package handlers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jdkato/prose/v2"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/vinewz/audiogen-cli/sqlc"
)

func SplitToSentences(ctx context.Context, tracker *progress.Tracker, db *sql.DB, title, lang string) error {
	queries := sqlc.New(db)

	bookId, err := queries.GetBookID(ctx, sql.NullString{String: title, Valid: true})
	if err != nil {
		return fmt.Errorf("Error getting BookId: %v", err)
	}

	txts, err := queries.GetText(ctx, bookId)
	if err != nil {
		return fmt.Errorf("Error getting texts from book %s: %v", title, err)
	}

  tracker.Total = int64(len(txts))

	for _, txt := range txts {
		doc, _ := prose.NewDocument(txt.Content)
    if txt.Chapter == "" {
      continue
    }
		chpId, err := queries.GetChapterID(ctx, txt.Chapter)
		if err != nil {
			return fmt.Errorf("Error getting chapter from book %s: %v", title, err)
		}

		for _, stc := range doc.Sentences() {
			queries.InsertSentence(
				ctx,
				sqlc.InsertSentenceParams{
					BookID:    bookId,
					ChapterID: chpId,
					Sentence:  stc.Text,
				},
			)
		}
    tracker.Increment(1)
	}

	return nil
}
