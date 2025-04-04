package cli

import (
	"math/rand"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"
)

var (
	messageColors = []text.Color{
		text.FgRed,
		text.FgGreen,
		text.FgYellow,
		text.FgBlue,
		text.FgMagenta,
		text.FgCyan,
		text.FgWhite,
	}
	rng       = rand.New(rand.NewSource(time.Now().UnixNano()))
	timeStart = time.Now()
)

func NewProgress() (progress.Writer) {
  pw := progress.NewWriter()

	pw.SetSortBy(progress.SortByPercentDsc)
	pw.SetStyle(progress.StyleDefault)
	pw.SetTrackerLength(40)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.Style().Colors = progress.StyleColorsExample
	pw.Style().Options.PercentFormat = "%4.1f%%"
  pw.Style().Visibility.ETA = true

	return pw
}
