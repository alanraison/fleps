package cmd

import (
	"fmt"

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
			round := model.RoundID(roundID)
			if roundID == "" {
				var err error
				round, err = database.GetLatestRound()
				if err != nil {
					return fmt.Errorf("getting latest round id: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStderr(), "Using latest round id: %s\n", round)
			}
			fs, err := c.ReadFixtureRows(cmd.InOrStdin(), round)
			if err != nil {
				return fmt.Errorf("reading fixture rows for round %q: %w", round, err)
			}

			if err := database.AddFixtures(fs); err != nil {
				return fmt.Errorf("adding fixtures for round %q: %w", round, err)
			}

			return nil
		},
	}
	list = &cobra.Command{
		Use:   "list",
		Short: "List fixtures for a round",
		RunE: func(cmd *cobra.Command, args []string) error {
			var round = model.RoundID(roundID)
			if roundID == "" {
				var err error
				round, err = database.GetLatestRound()
				if err != nil {
					return fmt.Errorf("getting latest round id: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStderr(), "Using latest round id: %s\n", round)
			}
			f, err := database.ListFixtures(round)
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
	addResult        = &cobra.Command{
		Use:   "add-result",
		Short: "Add a result for a fixture",
		RunE: func(cmd *cobra.Command, args []string) error {
			var round = model.RoundID(roundID)
			if roundID == "" {
				var err error
				round, err = database.GetLatestRound()
				if err != nil {
					return fmt.Errorf("getting latest round id: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStderr(), "Using latest round id: %s\n", round)
			}
			results, err := c.ReadResultRows(cmd.InOrStdin(), round)
			if err != nil {
				return fmt.Errorf("reading result rows for round %q: %w", round, err)
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
	addResult.Flags().StringVarP(&roundID, "round", "r", "", "round identifier for all results in the CSV")
	list.Flags().StringVarP(&roundID, "round", "r", "", "round identifier for listing fixtures")

	fixtures.AddCommand(createRound)
	fixtures.AddCommand(add)
	fixtures.AddCommand(list)
	fixtures.AddCommand(addResult)
}
