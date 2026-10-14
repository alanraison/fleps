package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/alanraison/fleps/pkg/model"
	"github.com/alanraison/fleps/pkg/render"
)

type fakeSource struct{ uploaded []byte }

func (f *fakeSource) LatestRound() (model.RoundID, error) { return "r1", nil }
func (f *fakeSource) Table(model.RoundID) (*render.Table, error) {
	return &render.Table{Players: []string{"Ann"}, RoundScores: []int{0}, SeasonScores: []int{0}}, nil
}
func (f *fakeSource) Upload(_ context.Context, name string, png []byte) (string, error) {
	f.uploaded = png
	return "https://example.com/" + name, nil
}

func TestHandleShowResults(t *testing.T) {
	src := &fakeSource{}
	out, err := (&handler{results: src}).handleShowResults(context.Background())
	if err != nil {
		t.Fatalf("handleShowResults: %v", err)
	}
	if len(src.uploaded) == 0 {
		t.Error("nothing uploaded")
	}
	var v map[string]any
	if err := json.Unmarshal(out, &v); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if !strings.Contains(string(out), `"imageUrl":"https://example.com/results/r1-`) {
		t.Errorf("response lacks image url: %s", out)
	}
}
