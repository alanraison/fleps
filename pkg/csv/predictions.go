package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/alanraison/predictions/pkg/model"
)

func (c *Csv) ReadPredictionRows(r io.Reader, roundID model.RoundID, player string) (model.GamePredictions, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = 4

	predictions := make(model.GamePredictions)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading CSV record: %w", err)
		}
		homeTeam, err := c.teamRepo.FindTeamByKey(model.TeamKey(record[0]))
		if err != nil {
			return nil, fmt.Errorf("failed to find home team: %w", err)
		}
		awayTeam, err := c.teamRepo.FindTeamByKey(model.TeamKey(record[1]))
		if err != nil {
			return nil, fmt.Errorf("failed to find away team: %w", err)
		}
		homeGoals, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, fmt.Errorf("parsing home goals: %w", err)
		}
		awayGoals, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, fmt.Errorf("parsing away goals: %w", err)
		}
		predictions[model.Game{
			HomeTeam: homeTeam.Key,
			AwayTeam: awayTeam.Key,
		}] = model.Prediction{
			HomeGoals: homeGoals,
			AwayGoals: awayGoals,
		}
	}
	return predictions, nil
}
