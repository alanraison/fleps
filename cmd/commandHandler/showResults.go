package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alanraison/fleps/pkg/chat"
	"github.com/alanraison/fleps/pkg/model"
	"github.com/alanraison/fleps/pkg/render"
)

type resultsSource interface {
	LatestRound() (model.RoundID, error)
	Table(model.RoundID) (*render.Table, error)
	Upload(ctx context.Context, name string, png []byte) (string, error)
}

type uploader interface {
	UploadPNG(ctx context.Context, name string, data []byte) (string, error)
}

type showResults struct {
	sources  render.Sources
	uploader uploader
}

func (s *showResults) LatestRound() (model.RoundID, error) {
	return s.sources.Results.GetLatestRoundWithResults()
}

func (s *showResults) Table(r model.RoundID) (*render.Table, error) {
	return s.sources.BuildTable(r)
}

func (s *showResults) Upload(ctx context.Context, name string, png []byte) (string, error) {
	return s.uploader.UploadPNG(ctx, name, png)
}

func (h *handler) handleShowResults(ctx context.Context) ([]byte, error) {
	roundID, err := h.results.LatestRound()
	if err != nil {
		return nil, fmt.Errorf("finding latest round with results: %w", err)
	}
	table, err := h.results.Table(roundID)
	if err != nil {
		return nil, fmt.Errorf("building results table for round %v: %w", roundID, err)
	}
	png, err := render.RenderPNG(table)
	if err != nil {
		return nil, fmt.Errorf("rendering results table for round %v: %w", roundID, err)
	}
	name := fmt.Sprintf("results/%v-%d.png", roundID, time.Now().UnixNano())
	url, err := h.results.Upload(ctx, name, png)
	if err != nil {
		return nil, fmt.Errorf("uploading results image for round %v: %w", roundID, err)
	}
	bytes, err := json.Marshal(chat.NewResponseMessageWithImage(
		fmt.Sprintf("Round %v results", roundID), url, "Table of predictions and results"))
	if err != nil {
		return nil, fmt.Errorf("encoding results response: %w", err)
	}
	return bytes, nil
}
