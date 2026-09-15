package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/alanraison/fleps/pkg/model"
)

func (c *Csv) ReadFixtureRows(r io.Reader, roundID model.RoundID) ([]model.Fixture, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = 3

	fixtures := []model.Fixture{}
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}
		if len(record) != 3 {
			return nil, fmt.Errorf("expected 3 fields per fixture row, got %d", len(record))
		}

		date, err := time.ParseInLocation("2006-01-02 15:04", record[0], time.Local)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		homeTeam, err := c.teamRepo.FindTeamByKey(model.TeamKey(record[1]))
		if err != nil {
			return nil, fmt.Errorf("failed to find home team: %w", err)
		}
		awayTeam, err := c.teamRepo.FindTeamByKey(model.TeamKey(record[2]))
		if err != nil {
			return nil, fmt.Errorf("failed to find away team: %w", err)
		}

		fixtures = append(fixtures, model.Fixture{
			FixtureKey: model.FixtureKey{
				RoundID:  roundID,
				HomeTeam: homeTeam.Key,
				AwayTeam: awayTeam.Key,
			},
			Date: date,
		})
	}
	return fixtures, nil
}
