package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gen2brain/go-fitz"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/vinewz/audiogen-cli/sqlc"
)

func ExtractText(ctx context.Context, tracker *progress.Tracker, db *sql.DB, title string, file *fitz.Document) error {
	queries := sqlc.New(db)

	if title == "" {
		txts := []string{}

		for i := 1; i < file.NumPage(); i++ {
			txt, err := file.Text(i)
			if err != nil {
				return fmt.Errorf("Error extracting text from page %d: %v", i, err.Error())
			}

			txts = append(txts, txt)
			tracker.Increment(1)
		}

		fmt.Println(strings.Join(txts, ""))
	} else {
		err := saveBook(ctx, queries, file, title)
		if err != nil {
			return err
		}
		err = saveChapters(ctx, queries, file, title)
		if err != nil {
			return err
		}
		err = saveText(ctx, tracker, queries, file, title)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveBook(ctx context.Context, queries *sqlc.Queries, file *fitz.Document, title string) error {
	err := queries.InsertBook(
		ctx,
		sqlc.InsertBookParams{
			ID:         uuid.NewString(),
			Title:      sql.NullString{String: title, Valid: true},
			TotalPages: int64(file.NumPage()),
		},
	)
	if err != nil {
		return fmt.Errorf("Error saving book: %v", err.Error())
	}

	return nil
}

func saveChapters(ctx context.Context, queries *sqlc.Queries, file *fitz.Document, title string) error {
	toc, err := file.ToC()
	if err != nil {
		return fmt.Errorf("Couldn't get TOC: %v", err.Error())
	}

	bookId, err := queries.GetBookID(
		ctx,
		sql.NullString{String: title, Valid: true},
	)
	if err != nil {
		return fmt.Errorf("Couldn't get BookId: %v", err.Error())
	}

	for idx := range toc {
		err := queries.InsertChapter(
			ctx,
			sqlc.InsertChapterParams{
				BookID:    bookId,
				StartPage: int64(toc[idx].Page),
				Chapter:   toc[idx].Title,
			},
		)
		if err != nil {
			return fmt.Errorf("Error inserting Chapter: %v", err.Error())
		}
	}

	return nil
}

func saveText(ctx context.Context, tracker *progress.Tracker, queries *sqlc.Queries, file *fitz.Document, title string) error {
	bookId, err := queries.GetBookID(
		ctx,
		sql.NullString{String: title, Valid: true},
	)
	if err != nil {
		return fmt.Errorf("Error getting book id: %v", err)
	}

	chapters, err := queries.GetChapters(ctx, bookId)
	if err != nil {
		return fmt.Errorf("Error getting chapters: %v", err)
	}

	for i := 0; i <= file.NumPage()-1; i++ {
		txt, err := file.Text(i)
		if err != nil {
			return fmt.Errorf("Error getting text from page %d: %v", i, err)
		}

		sanTxt, err := sanitizeTxt(txt)
		if err != nil {
			continue
		}

		var currChapter string
		for _, ch := range chapters {
			if ch.StartPage <= int64(i) {
				currChapter = ch.Chapter
			}
		}

		err = queries.InsertText(
			ctx,
			sqlc.InsertTextParams{
				BookID:  bookId,
				Page:    int64(i),
				Content: sanTxt,
				Chapter: currChapter,
			},
		)
		if err != nil {
			return fmt.Errorf("Error saving book: %v", err.Error())
		}
		tracker.Increment(1)
	}

	return nil
}

func sanitizeTxt(txt string) (string, error) {
	noBkLine := strings.ReplaceAll(txt, "\n", " ")
	trimmed := strings.TrimSpace(noBkLine)

	if trimmed == "" {
		return "", fmt.Errorf("Empty string")
	}

	return trimmed, nil
}
