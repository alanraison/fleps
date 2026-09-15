package cmd

import (
	"fmt"

	"github.com/alanraison/predictions/pkg/model"
	"github.com/spf13/cobra"
)

var (
	player         string
	predictionsCmd = &cobra.Command{
		Use:   "predictions",
		Short: "Manage predictions",
	}
	addPrediction = &cobra.Command{
		Use:   "add",
		Short: "Add new predictions",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := getAppContext(cmd)
			if ctx == nil {
				return fmt.Errorf("app context not found")
			}
			round := model.RoundID(roundID)
			if roundID == "" {
				var err error
				round, err = ctx.fixtureRepo.GetLatestRoundWithNoResults()
				if err != nil {
					return fmt.Errorf("getting latest round id: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStderr(), "Using latest round id: %s\n", round)
			}
			predictions, err := c.ReadPredictionRows(cmd.InOrStdin(), round, player)
			if err != nil {
				return fmt.Errorf("reading prediction rows: %w", err)
			}
			if err := ctx.predictionRepo.AddPredictions(player, round, predictions); err != nil {
				return fmt.Errorf("adding predictions: %w", err)
			}
			return nil
		},
	}
	listPredictions = &cobra.Command{
		Use:   "list",
		Short: "List all predictions",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := getAppContext(cmd)
			if ctx == nil {
				return fmt.Errorf("app context not found")
			}
			round := model.RoundID(roundID)
			if roundID == "" {
				var err error
				round, err = ctx.fixtureRepo.GetLatestRoundWithNoResults()
				if err != nil {
					return fmt.Errorf("getting latest round id: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStderr(), "Using latest round id: %s\n", round)
			}
			predictions, err := ctx.predictionRepo.ListPredictions(round)
			if err != nil {
				return fmt.Errorf("listing predictions for round %q: %w", round, err)
			}
			for player, games := range predictions {
				for game, prediction := range games {
					fmt.Fprintf(cmd.OutOrStdout(), "Player: %s, Game: %s vs %s, Prediction: %d - %d\n", player, game.HomeTeam, game.AwayTeam, prediction.HomeGoals, prediction.AwayGoals)
				}
			}
			return nil
		},
	}
)

func init() {
	addPrediction.Flags().StringVarP(&player, "player", "p", "", "email address of the player making the prediction")
	addPrediction.Flags().StringVarP(&roundID, "round", "r", "", "round identifier for the prediction")
	addPrediction.MarkFlagRequired("player")

	listPredictions.Flags().StringVarP(&roundID, "round", "r", "", "round identifier for the predictions")

	rootCmd.AddCommand(predictionsCmd)
	predictionsCmd.AddCommand(addPrediction)
	predictionsCmd.AddCommand(listPredictions)
}
