package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/alanraison/predictions/pkg/model"
)

func (c *Csv) ReadFixtureRows(r io.Reader) ([]model.Fixture, error) {
	cr := csv.NewReader(r)
	record, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV record: %w", err)
	}
	fixtures := []model.Fixture{}
	for {
		date, err := time.Parse("2006-01-02 15:04", record[0])
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
			Date:     date,
			HomeTeam: homeTeam.Key,
			AwayTeam: awayTeam.Key,
		})
		record, err = cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}
	}
	return fixtures, nil
}
