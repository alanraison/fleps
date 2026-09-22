package main

import (
	"fmt"

	"github.com/alanraison/fleps/pkg/model"
	"github.com/spf13/cobra"
)

var (
	seasonFlag string
	scoresCmd  = &cobra.Command{
		Use:   "scores",
		Short: "Manage scores",
	}
	getRoundScoresCmd = &cobra.Command{
		Use:   "round",
		Short: "Get round scores",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := getAppContext(cmd)
			if ctx == nil {
				return fmt.Errorf("app context not found")
			}
			var round model.RoundID
			if roundID == "" {
				var err error
				round, err = ctx.resultRepo.GetLatestRoundWithResults()
				if err != nil {
					return fmt.Errorf("getting latest round id: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStderr(), "Using latest round id: %s\n", round)
			}
			scores, err := ctx.scoreService.CalculateRoundScores(round)
			if err != nil {
				return fmt.Errorf("calculating round scores: %w", err)
			}
			for player, score := range scores {
				fmt.Printf("%s: %d\n", player, score)
			}
			return nil
		},
	}
	getSeasonScoresCmd = &cobra.Command{
		Use:   "season",
		Short: "Get season scores",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := getAppContext(cmd)
			if ctx == nil {
				return fmt.Errorf("app context not found")
			}
			var scores map[string]int
			var err error
			if seasonFlag == "" {
				scores, err = ctx.scoreService.CalculateCurrentSeasonScores()
			} else {
				scores, err = ctx.scoreService.CalculateSeasonScores(model.SeasonID(seasonFlag))
			}
			if err != nil {
				return fmt.Errorf("calculating scores: %w", err)
			}
			for player, score := range scores {
				fmt.Printf("%s: %d\n", player, score)
			}
			return nil
		},
	}
)

func init() {
	scoresCmd.AddCommand(getRoundScoresCmd)
	scoresCmd.AddCommand(getSeasonScoresCmd)
	rootCmd.AddCommand(scoresCmd)

	getSeasonScoresCmd.Flags().StringVar(&seasonFlag, "season", "", "Specify the season for which to get scores")
}
