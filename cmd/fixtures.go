package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/alanraison/predictions/pkg/model"
	"github.com/spf13/cobra"
)

var (
	addRoundID string
	fixtures   = &cobra.Command{
		Use:   "fixtures",
		Short: "Manage fixtures",
	}
	add = &cobra.Command{
		Use:   "add",
		Short: "Add fixtures from a CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {
			fs, err := c.ReadFixtureRows(cmd.InOrStdin(), model.RoundID(addRoundID))
			if err != nil {
				return fmt.Errorf("reading fixture rows for round %q: %w", addRoundID, err)
			}

			if err := database.AddFixtures(fs); err != nil {
				return fmt.Errorf("adding fixtures for round %q: %w", addRoundID, err)
			}

			return nil
		},
	}
	listFrom, listTo string
	listTeams        []string
	list             = &cobra.Command{
		Use:   "list",
		Short: "List all fixtures",
		RunE: func(cmd *cobra.Command, args []string) error {
			fromDate, err := time.Parse("2006-01-02", listFrom)
			if err != nil {
				return fmt.Errorf("parsing from date %q: %w", listFrom, err)
			}
			toDate, err := time.Parse("2006-01-02", listTo)
			if err != nil {
				return fmt.Errorf("parsing to date %q: %w", listTo, err)
			}
			f, err := database.ListFixtures(fromDate, toDate, listTeams)
			if err != nil {
				return fmt.Errorf("listing fixtures: %w", err)
			}
			for _, fixture := range f {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %v %v %v\n", fixture.RoundID, fixture.Date.Local().Format("02/01 15:04"), fixture.HomeTeam, fixture.AwayTeam); err != nil {
					return fmt.Errorf("writing fixture output: %w", err)
				}
			}

			return nil
		},
	}
	addResultRoundID string
	addResult        = &cobra.Command{
		Use:   "add-result",
		Short: "Add a result for a fixture",
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := c.ReadResultRows(cmd.InOrStdin(), model.RoundID(addResultRoundID))
			if err != nil {
				return fmt.Errorf("reading result rows for round %q: %w", addResultRoundID, err)
			}

			for _, result := range results {
				if err := database.AddResult(result.RoundID, result.HomeTeam, result.AwayTeam, result.HomeScore, result.AwayScore); err != nil {
					return fmt.Errorf("adding result for round %q and fixture %s vs %s: %w", result.RoundID, result.HomeTeam, result.AwayTeam, err)
				}
			}

			return nil
		},
	}
)

func init() {
	add.Flags().StringVarP(&addRoundID, "round", "r", "", "round identifier for all fixtures in the CSV")
	err := add.MarkFlagRequired("round")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'round' flag as required on fixtures add: %v\n", err)
	}

	list.Flags().StringVarP(&listFrom, "from", "f", "", "from date")
	list.Flags().StringVarP(&listTo, "to", "t", "", "to date")
	list.Flags().StringArrayVarP(&listTeams, "teams", "m", []string{}, "filter by teams")
	err = list.MarkFlagRequired("from")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'from' flag as required: %v\n", err)
	}
	err = list.MarkFlagRequired("to")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'to' flag as required: %v\n", err)
	}

	addResult.Flags().StringVarP(&addResultRoundID, "round", "r", "", "round identifier for all results in the CSV")
	err = addResult.MarkFlagRequired("round")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'round' flag as required on fixtures add-result: %v\n", err)
	}

	fixtures.AddCommand(add)
	fixtures.AddCommand(list)
	fixtures.AddCommand(addResult)
}
