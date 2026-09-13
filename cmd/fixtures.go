package cmd

import (
	"fmt"

	"github.com/alanraison/predictions/pkg/model"
	"github.com/spf13/cobra"
)

var (
	roundID     string
	fixturesCmd = &cobra.Command{
		Use:   "fixtures",
		Short: "Manage fixtures",
	}
	seasonID       string
	createRoundCmd = &cobra.Command{
		Use:   "create-round",
		Short: "Creates a Round",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("No round supplied")
			}
			ctx := getAppContext(cmd)
			if ctx == nil {
				return fmt.Errorf("app context not found")
			}
			if err := ctx.fixtureRepo.AddRound(model.RoundID(args[0]), model.SeasonID(seasonID)); err != nil {
				return fmt.Errorf("adding round: %w", err)
			}
			return nil
		},
	}
	addFixturesCmd = &cobra.Command{
		Use:   "add",
		Short: "Add fixtures from a CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := getAppContext(cmd)
			if ctx == nil {
				return fmt.Errorf("app context not found")
			}
			round := model.RoundID(roundID)
			if roundID == "" {
				round, err := ctx.fixtureRepo.GetLatestRoundWithNoResults()
				if err != nil {
					return fmt.Errorf("getting latest round id: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStderr(), "Using latest round id: %s\n", round)
			}
			fs, err := c.ReadFixtureRows(cmd.InOrStdin(), round)
			if err != nil {
				return fmt.Errorf("reading fixture rows for round %q: %w", round, err)
			}

			if err := ctx.fixtureRepo.AddFixtures(fs); err != nil {
				return fmt.Errorf("adding fixtures for round %q: %w", round, err)
			}

			return nil
		},
	}
	listFixturesCmd = &cobra.Command{
		Use:   "list",
		Short: "List fixtures for a round",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := getAppContext(cmd)
			if ctx == nil {
				return fmt.Errorf("app context not found")
			}
			var round = model.RoundID(roundID)
			if roundID == "" {
				round, err := ctx.fixtureRepo.GetLatestRoundWithNoResults()
				if err != nil {
					return fmt.Errorf("getting latest round id: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStderr(), "Using latest round id: %s\n", round)
			}
			f, err := ctx.fixtureRepo.ListFixtures(round)
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
)

func init() {
	createRoundCmd.Flags().StringVarP(&seasonID, "season", "s", "", "season in which the round belongs")
	addFixturesCmd.Flags().StringVarP(&roundID, "round", "r", "", "round identifier for all fixtures in the CSV")
	listFixturesCmd.Flags().StringVarP(&roundID, "round", "r", "", "round identifier for listing fixtures")

	rootCmd.AddCommand(fixturesCmd)
	fixturesCmd.AddCommand(createRoundCmd)
	fixturesCmd.AddCommand(addFixturesCmd)
	fixturesCmd.AddCommand(listFixturesCmd)
}
