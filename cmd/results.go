package cmd

import (
	"fmt"

	"github.com/alanraison/predictions/pkg/model"
	"github.com/spf13/cobra"
)

var (
	resultsCmd = &cobra.Command{
		Use:   "results",
		Short: "Manage results for fixtures",
	}
	addResult = &cobra.Command{
		Use:   "add",
		Short: "Add a results for a fixture",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := getAppContext(cmd)
			if ctx == nil {
				return fmt.Errorf("app context not found")
			}
			var round = model.RoundID(roundID)
			if roundID == "" {
				var err error
				round, err = ctx.fixtureRepo.GetLatestRoundWithNoResults()
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
				if err := ctx.resultRepo.AddResult(model.FixtureKey{
					RoundID:  result.RoundID,
					HomeTeam: result.HomeTeam,
					AwayTeam: result.AwayTeam,
				}, result.HomeGoals, result.AwayGoals); err != nil {
					return fmt.Errorf("adding result for round %q and fixture %s vs %s: %w", result.RoundID, result.HomeTeam, result.AwayTeam, err)
				}
			}

			return nil
		},
	}
	listResults = &cobra.Command{
		Use:   "list",
		Short: "List all results for a round",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := getAppContext(cmd)
			if ctx == nil {
				return fmt.Errorf("app context not found")
			}
			round := model.RoundID(roundID)
			if roundID == "" {
				var err error
				round, err = ctx.resultRepo.GetLatestRoundWithResults()
				if err != nil {
					return fmt.Errorf("getting latest round id: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStderr(), "Using latest round id: %s\n", round)
			}
			results, err := ctx.resultRepo.ListResults(round)
			if err != nil {
				return fmt.Errorf("listing results for round %q: %w", round, err)
			}
			for _, result := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "%s %d - %d %s\n", result.HomeTeam, result.HomeGoals, result.AwayGoals, result.AwayTeam)
			}
			return nil
		},
	}
)

func init() {
	addResult.Flags().StringVarP(&roundID, "round", "r", "", "round identifier for all results in the CSV")
	listResults.Flags().StringVarP(&roundID, "round", "r", "", "round identifier for all results")

	rootCmd.AddCommand(resultsCmd)
	resultsCmd.AddCommand(addResult)
	resultsCmd.AddCommand(listResults)
}
