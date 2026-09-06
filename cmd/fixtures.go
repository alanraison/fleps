package cmd

import (
	"fmt"
	"os"

	"github.com/alanraison/predictions/pkg/model"
	"github.com/spf13/cobra"
)

var (
	roundID  string
	fixtures = &cobra.Command{
		Use:   "fixtures",
		Short: "Manage fixtures",
	}
	seasonID    string
	createRound = &cobra.Command{
		Use:   "create-round",
		Short: "Creates a Round",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("No round supplied")
			}
			if err := database.AddRound(model.RoundID(args[0]), model.SeasonID(seasonID)); err != nil {
				return fmt.Errorf("adding round: %w", err)
			}
			return nil
		},
	}
	add = &cobra.Command{
		Use:   "add",
		Short: "Add fixtures from a CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {
			fs, err := c.ReadFixtureRows(cmd.InOrStdin(), model.RoundID(roundID))
			if err != nil {
				return fmt.Errorf("reading fixture rows for round %q: %w", roundID, err)
			}

			if err := database.AddFixtures(fs); err != nil {
				return fmt.Errorf("adding fixtures for round %q: %w", roundID, err)
			}

			return nil
		},
	}
	list = &cobra.Command{
		Use:   "list",
		Short: "List fixtures for a round",
		RunE: func(cmd *cobra.Command, args []string) error {

			f, err := database.ListFixtures(model.RoundID(roundID))
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
	createRound.Flags().StringVarP(&seasonID, "season", "s", "", "season in which the round belongs")
	add.Flags().StringVarP(&roundID, "round", "r", "", "round identifier for all fixtures in the CSV")
	err := add.MarkFlagRequired("round")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'round' flag as required on fixtures add: %v\n", err)
	}

	addResult.Flags().StringVarP(&addResultRoundID, "round", "r", "", "round identifier for all results in the CSV")
	err = addResult.MarkFlagRequired("round")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'round' flag as required on fixtures add-result: %v\n", err)
	}

	fixtures.AddCommand(createRound)
	fixtures.AddCommand(add)
	fixtures.AddCommand(list)
	fixtures.AddCommand(addResult)
}
