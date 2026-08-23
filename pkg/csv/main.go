package csv

import (
	encodingcsv "encoding/csv"
	"io"

	"github.com/alanraison/predictions/pkg/model"
)

type Csv struct {
	teamRepo    model.TeamRepository
	fixtureRepo model.FixtureRepository
}

type PredictionRow struct {
	Match       string
	Predictions map[string]string
}

func NewCsv(teamRepo model.TeamRepository, fixtureRepo model.FixtureRepository) *Csv {
	return &Csv{
		teamRepo:    teamRepo,
		fixtureRepo: fixtureRepo,
	}
}

func readPredictionRows(r io.Reader) ([]PredictionRow, error) {
	reader := encodingcsv.NewReader(r)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}

	headers := records[0]
	rows := make([]PredictionRow, 0, len(records)-1)

	for _, record := range records[1:] {
		row := PredictionRow{Predictions: make(map[string]string)}
		for i, value := range record {
			if i >= len(headers) {
				break
			}
			header := headers[i]
			if header == "match" {
				row.Match = value
				continue
			}
			row.Predictions[header] = value
		}
		rows = append(rows, row)
	}

	return rows, nil
}
