// Package render draws the results table as an image.
package render

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

// Points awarded for a prediction.
const (
	PointsNone    = 0
	PointsOutcome = 1
	PointsExact   = 3
)

var (
	exactColour   = color.RGBA{R: 0x39, G: 0xff, B: 0x14, A: 0xff}
	outcomeColour = color.RGBA{R: 0x00, G: 0x64, B: 0x00, A: 0xff}
	headerColour  = color.RGBA{R: 0xdd, G: 0xdd, B: 0xdd, A: 0xff}
	summaryColour = color.RGBA{R: 0xee, G: 0xee, B: 0xee, A: 0xff}
	lineColour    = color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff}
	white         = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	black         = color.RGBA{A: 0xff}
)

// scale multiplies every pixel dimension; Chat scales the image to the card
// width, so a higher value gives sharper, larger text.
const (
	scale     = 2
	padding   = 10 * scale
	rowHeight = 34 * scale
	fontSize  = 16 * scale
)

// Cell is a single prediction, with the points it earned.
type Cell struct {
	Text   string
	Points int
}

// Row is one match: its label, each player's prediction and the actual result.
type Row struct {
	Label       string
	Predictions []Cell
	Result      string
}

// Table is the data rendered by RenderPNG. Players, Predictions in each row,
// RoundScores and SeasonScores are all in the same (column) order.
type Table struct {
	Players      []string
	Rows         []Row
	RoundScores  []int
	SeasonScores []int
}

// Validate checks that every row has one value per player.
func (t *Table) Validate() error {
	n := len(t.Players)
	for i, r := range t.Rows {
		if len(r.Predictions) != n {
			return fmt.Errorf("row %d has %d predictions for %d players", i, len(r.Predictions), n)
		}
	}
	if len(t.RoundScores) != n || len(t.SeasonScores) != n {
		return fmt.Errorf("expected %d scores, got %d round and %d season", n, len(t.RoundScores), len(t.SeasonScores))
	}
	return nil
}

type faces struct{ regular, bold font.Face }

func newFaces() (*faces, error) {
	load := func(data []byte) (font.Face, error) {
		f, err := opentype.Parse(data)
		if err != nil {
			return nil, fmt.Errorf("parsing font: %w", err)
		}
		face, err := opentype.NewFace(f, &opentype.FaceOptions{Size: fontSize, DPI: 72})
		if err != nil {
			return nil, fmt.Errorf("creating font face: %w", err)
		}
		return face, nil
	}
	reg, err := load(goregular.TTF)
	if err != nil {
		return nil, err
	}
	bold, err := load(gobold.TTF)
	if err != nil {
		return nil, err
	}
	return &faces{regular: reg, bold: bold}, nil
}

// RenderPNG draws the table and returns it PNG-encoded.
func RenderPNG(t *Table) ([]byte, error) {
	if err := t.Validate(); err != nil {
		return nil, fmt.Errorf("validating table: %w", err)
	}
	f, err := newFaces()
	if err != nil {
		return nil, fmt.Errorf("loading fonts: %w", err)
	}

	// Columns: match label, one per player, result.
	header := append(append([]string{""}, t.Players...), "Result")
	nCols := len(header)
	widths := make([]int, nCols)
	measure := func(col int, s string, face font.Face) {
		if w := font.MeasureString(face, s).Ceil() + 2*padding; w > widths[col] {
			widths[col] = w
		}
	}
	for c, h := range header {
		measure(c, h, f.bold)
	}
	for _, r := range t.Rows {
		measure(0, r.Label, f.regular)
		for i, p := range r.Predictions {
			measure(i+1, p.Text, f.regular)
		}
		measure(nCols-1, r.Result, f.bold)
	}
	measure(0, "Round score", f.bold)
	measure(0, "Season score", f.bold)
	for i := range t.Players {
		measure(i+1, fmt.Sprint(t.RoundScores[i]), f.bold)
		measure(i+1, fmt.Sprint(t.SeasonScores[i]), f.bold)
	}

	nRows := len(t.Rows) + 3
	totalW := 1
	for _, w := range widths {
		totalW += w
	}
	dc := gg.NewContext(totalW, nRows*rowHeight+1)
	dc.SetColor(white)
	dc.Clear()

	y := 0
	drawRow := func(texts []string, fills []color.Color, face font.Face) {
		dc.SetFontFace(face)
		x := 0
		for c, s := range texts {
			if fills[c] != nil {
				dc.SetColor(fills[c])
				dc.DrawRectangle(float64(x), float64(y), float64(widths[c]), rowHeight)
				dc.Fill()
			}
			dc.SetColor(lineColour)
			dc.SetLineWidth(scale)
			dc.DrawRectangle(float64(x)+0.5, float64(y)+0.5, float64(widths[c]), rowHeight)
			dc.Stroke()
			if fills[c] == outcomeColour {
				dc.SetColor(white)
			} else {
				dc.SetColor(black)
			}
			dc.DrawStringAnchored(s, float64(x)+float64(widths[c])/2, float64(y)+rowHeight/2, 0.5, 0.35)
			x += widths[c]
		}
		y += rowHeight
	}
	fillsOf := func(c color.Color) []color.Color {
		fills := make([]color.Color, nCols)
		for i := range fills {
			fills[i] = c
		}
		return fills
	}

	drawRow(header, fillsOf(headerColour), f.bold)
	for _, r := range t.Rows {
		texts := []string{r.Label}
		fills := fillsOf(nil)
		for i, p := range r.Predictions {
			texts = append(texts, p.Text)
			switch p.Points {
			case PointsExact:
				fills[i+1] = exactColour
			case PointsOutcome:
				fills[i+1] = outcomeColour
			}
		}
		texts = append(texts, r.Result)
		drawRow(texts, fills, f.regular)
	}
	summary := func(label string, scores []int) {
		texts := []string{label}
		for _, s := range scores {
			texts = append(texts, fmt.Sprint(s))
		}
		texts = append(texts, "")
		drawRow(texts, fillsOf(summaryColour), f.bold)
	}
	summary("Round score", t.RoundScores)
	summary("Season score", t.SeasonScores)

	var buf bytes.Buffer
	if err := dc.EncodePNG(&buf); err != nil {
		return nil, fmt.Errorf("encoding png: %w", err)
	}
	return buf.Bytes(), nil
}
