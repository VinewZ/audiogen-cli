package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/iFaceless/godub"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/vinewz/audiogen-cli/sqlc"
)

func Merge(ctx context.Context, tracker *progress.Tracker, db *sql.DB, title, lang, delay string) error {
	tmpDir := os.TempDir()
	queries := sqlc.New(db)

	bookId, err := queries.GetBookID(ctx, sql.NullString{String: title, Valid: true})
	if err != nil {
		return fmt.Errorf("Error getting BookId: %v", err)
	}

	chapters, err := queries.GetChapters(ctx, bookId)
	if err != nil {
		return fmt.Errorf("Error getting chapters: %v", err)
	}
  tracker.Total = int64(len(chapters))

	fDelay, err := strconv.ParseFloat(delay, 64)
	silence, err := createSilentAudio(tmpDir, title, fDelay)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	outputPath := path.Join(wd, title)
	err = os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	silenceSeg, err := godub.NewLoader().Load(silence)
	if err != nil {
		return fmt.Errorf("Error reading silenceSeg: %v", err)
	}

	for _, chap := range chapters {
		audioName := fmt.Sprintf("%s-%s.wav", title, chap.Chapter)
		err = godub.NewExporter(path.Join(outputPath, audioName)).WithDstFormat("wav").WithBitRate(128).Export(silenceSeg)
		if err != nil {
			return err
		}

		stcs, err := queries.GetAudioWhereChapter(ctx, chap.ID)
		if err != nil {
			return fmt.Errorf("Error getting audios from chapter %s: %v", chap.Chapter, err)
		}

		for _, stc := range stcs {
			mainSeg, err := godub.NewLoader().Load(path.Join(outputPath, audioName))
			if err != nil {
				return fmt.Errorf("Error getting main segment: %v", err)
			}

			newSeg, err := godub.NewLoader().Load(stc.AudioFilePath.String)
			if err != nil {
				return fmt.Errorf("Error getting new segment: %v", err)
			}

			mainSeg, err = mainSeg.Append(newSeg, silenceSeg)
			if err != nil {
				return fmt.Errorf("Error appending new segment: %v", err)
			}

			err = godub.NewExporter(path.Join(outputPath, audioName)).WithDstFormat("wav").WithBitRate(128).Export(mainSeg)
      if err != nil {
        return fmt.Errorf("Error exporting segment %d: %v", stc.ID, err)
      }
		}
    tracker.Increment(1)
	}

	return nil
}

func createSilentAudio(tmpDir, title string, duration float64) (string, error) {
	segment, err := godub.NewSilentAudioSegment(duration, 24000)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(path.Join(tmpDir, title), os.ModePerm)
	if err != nil {
		return "", err
	}

	outputPath := path.Join(tmpDir, title, "silence.wav")
	err = godub.NewExporter(outputPath).WithDstFormat("wav").WithBitRate(128).Export(segment)
	if err != nil {
		return "", err
	}
	return outputPath, nil
}
