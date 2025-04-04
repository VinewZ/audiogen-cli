package cli

import (
	"context"
	"fmt"

	"github.com/gen2brain/go-fitz"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/urfave/cli/v3"
	"github.com/vinewz/audiogen-cli/cmd/cli/handlers"
	"github.com/vinewz/audiogen-cli/cmd/db"
)

var Cmds = &cli.Command{
	Commands: []*cli.Command{
		{
			Name:    "extract",
			Aliases: []string{"e"},
			Usage:   "Extract text from the given PDF, save it to the DB or return it in stdout.",
			Action:  handleExtract,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "path",
					Aliases:  []string{"p"},
					Usage:    "Path to the PDF to extract text from.",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "title",
					Aliases: []string{"t"},
					Usage:   "If set, save the extracted text and the book chapters to the DB, else returns it in stdout.",
				},
			},
		},
		{
			Name:    "split",
			Aliases: []string{"s"},
			Usage:   "Split text from the given book into sentences.",
			Action:  handleSplit,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "title",
					Aliases:  []string{"t"},
					Usage:    "Title of the book to split text from.",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "language",
					Aliases:  []string{"l"},
					Usage:    "Language of the book.",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "chapter",
					Aliases: []string{"c"},
					Usage:   "If set split the text from a single Chapter of the book.",
				},
			},
		},
		{
			Name:    "narrate",
			Aliases: []string{"n"},
			Usage:   "Post sentences to the TTS Api and retrieve the audio file.",
			Action:  handleNarrate,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "title",
					Aliases:  []string{"t"},
					Usage:    "Title of the book to split text from.",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "language",
					Aliases:  []string{"l"},
					Usage:    "Language of the book.",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "chapter",
					Aliases: []string{"c"},
					Usage:   "If set post the sentences from a single Chapter of the book.",
				},
			},
		},
		{
			Name:    "merge",
			Aliases: []string{"m"},
			Usage:   "Merge all the audio files into a single mp3 file.",
			Action:  handleMerge,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "title",
					Aliases:  []string{"t"},
					Usage:    "Title of the book to split text from.",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "delay",
					Aliases:  []string{"d"},
					Usage:    "Delay between sentences in milliseconds.",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "language",
					Aliases:  []string{"l"},
					Usage:    "Language of the book.",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "chapter",
					Aliases: []string{"c"},
					Usage:   "If set merge the text from a single Chapter of the book.",
				},
			},
		},
	},
}

func handleExtract(ctx context.Context, cmd *cli.Command) error {
	pdfPath := cmd.String("path")
	title := cmd.String("title")

	db, err := db.Connection()
	if err != nil {
		return err
	}
	defer db.Close()

	file, err := fitz.New(pdfPath)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	pw := NewProgress()
	go pw.Render()
	defer pw.Stop()

	tracker := progress.Tracker{
		Message: "Extracting text from pages",
		Total:   int64(file.NumPage()),
	}
	pw.AppendTracker(&tracker)
	err = handlers.ExtractText(ctx, &tracker, db, title, file)
	if err != nil {
    tracker.MarkAsErrored()
		return err
	}

	tracker.MarkAsDone()
	return nil
}

func handleSplit(ctx context.Context, cmd *cli.Command) error {
	title := cmd.String("title")
	// chapter := cmd.String("chapter")
	lang := cmd.String("language")

	db, err := db.Connection()
	if err != nil {
		return err
	}
	defer db.Close()

	pw := NewProgress()
	go pw.Render()
	defer pw.Stop()

	tracker := progress.Tracker{
		Message: "Splitting to sentences",
	}
	pw.AppendTracker(&tracker)
	err = handlers.SplitToSentences(ctx, &tracker, db, title, lang)
	if err != nil {
    tracker.MarkAsErrored()
		return err
	}

	tracker.MarkAsDone()
	return nil
}

func handleNarrate(ctx context.Context, cmd *cli.Command) error {
	title := cmd.String("title")
	// chapter := cmd.String("chapter")
	lang := cmd.String("language")

	db, err := db.Connection()
	if err != nil {
		return err
	}

	pw := NewProgress()
	go pw.Render()
	defer pw.Stop()

	tracker := progress.Tracker{
		Message: "Generating audios",
	}
	pw.AppendTracker(&tracker)

	err = handlers.Narrate(ctx, &tracker, db, title, lang)
	if err != nil {
    tracker.MarkAsErrored()
		return err
	}

	tracker.MarkAsDone()
	return nil
}

func handleMerge(ctx context.Context, cmd *cli.Command) error {
	title := cmd.String("title")
	// chapter := cmd.String("chapter")
	lang := cmd.String("language")
	delay := cmd.String("delay")

	db, err := db.Connection()
	if err != nil {
		return err
	}

	pw := NewProgress()
	go pw.Render()
	defer pw.Stop()

	tracker := progress.Tracker{
		Message: "Merging audios",
	}
	pw.AppendTracker(&tracker)

	err = handlers.Merge(ctx, &tracker, db, title, lang, delay)
	if err != nil {
    tracker.MarkAsErrored()
		return err
	}

	tracker.MarkAsDone()
	return nil
}
