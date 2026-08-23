package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/alanraison/predictions/pkg/model"
)

func (c *Csv) readFixtures(r io.Reader) ([]model.Fixture, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}
	// skip the header row
	if len(records) > 0 {
		records = records[1:]
	}
	fixtures := make([]model.Fixture, len(records))
	for i, record := range records {
		home := record[0]
		away := record[1]
		date := record[2]
		h, err := c.teamRepo.FindByKey(home)
		if err != nil {
			return nil, fmt.Errorf("failed to find home team: %w", err)
		}
		a, err := c.teamRepo.FindByKey(away)
		if err != nil {
			return nil, fmt.Errorf("failed to find away team: %w", err)
		}
		d, err := time.Parse("2006-01-02T15:04", date)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		fixtures[i] = model.Fixture{
			HomeTeam: h,
			AwayTeam: a,
			Date:     d,
		}
	}
	return fixtures, nil
}

func (c *Csv) AddFixtures(r io.Reader) error {
	fs, err := c.readFixtures(r)
	if err != nil {
		return fmt.Errorf("reading fixtures: %w", err)
	}
	if err := c.fixtureRepo.AddFixtures(fs); err != nil {
		return fmt.Errorf("adding fixtures: %w", err)
	}
	return nil
}
