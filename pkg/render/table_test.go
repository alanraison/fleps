package render

import (
	"bytes"
	"image/png"
	"testing"
)

func TestRenderPNG(t *testing.T) {
	table := &Table{
		Players: []string{"Ann", "Bob"},
		Rows: []Row{{
			Label:       "ARS v CHE",
			Predictions: []Cell{{Text: "1-0", Points: PointsExact}, {Text: "2-0", Points: PointsOutcome}},
			Result:      "1-0",
		}},
		RoundScores:  []int{3, 1},
		SeasonScores: []int{10, 8},
	}
	data, err := RenderPNG(table)
	if err != nil {
		t.Fatalf("RenderPNG: %v", err)
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decoding png: %v", err)
	}
	var sawBright, sawDark bool
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, _ := img.At(x, y).RGBA()
			c := [3]uint32{r >> 8, g >> 8, bl >> 8}
			sawBright = sawBright || c == [3]uint32{0x39, 0xff, 0x14}
			sawDark = sawDark || c == [3]uint32{0, 0x64, 0}
		}
	}
	if !sawBright || !sawDark {
		t.Errorf("expected bright (%v) and dark (%v) green cells", sawBright, sawDark)
	}
}

func TestRenderPNGRejectsMismatchedRows(t *testing.T) {
	table := &Table{
		Players:      []string{"Ann"},
		Rows:         []Row{{Label: "x"}},
		RoundScores:  []int{0},
		SeasonScores: []int{0},
	}
	if _, err := RenderPNG(table); err == nil {
		t.Fatal("expected error")
	}
}
