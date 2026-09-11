package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/alanraison/predictions/pkg/model"
)

func (c *Csv) ReadResultRows(input io.Reader, roundID model.RoundID) ([]*model.Result, error) {
	reader := csv.NewReader(input)
	reader.FieldsPerRecord = -1

	results := []*model.Result{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}
		if len(record) != 4 {
			return nil, fmt.Errorf("expected 4 fields per result row, got %d", len(record))
		}

		homeTeam, err := c.teamRepo.FindTeamByKey(model.TeamKey(record[0]))
		if err != nil {
			return nil, fmt.Errorf("failed to find home team: %w", err)
		}
		awayTeam, err := c.teamRepo.FindTeamByKey(model.TeamKey(record[1]))
		if err != nil {
			return nil, fmt.Errorf("failed to find away team: %w", err)
		}
		homeScore, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, fmt.Errorf("failed to parse home score: %w", err)
		}
		awayScore, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, fmt.Errorf("failed to parse away score: %w", err)
		}

		results = append(results, &model.Result{
			FixtureKey: model.FixtureKey{
				RoundID:  roundID,
				HomeTeam: homeTeam.Key,
				AwayTeam: awayTeam.Key,
			},
			HomeGoals: homeScore,
			AwayGoals: awayScore,
		})
	}

	return results, nil
}
