package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	playerName  string
	playerEmail string
	playersCmd  = &cobra.Command{
		Use:   "players",
		Short: "Manage players",
	}
	addPlayer = &cobra.Command{
		Use:   "add",
		Short: "Add a new player",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := getAppContext(cmd)
			if ctx == nil {
				return fmt.Errorf("app context not found")
			}
			if playerName == "" || playerEmail == "" {
				return fmt.Errorf("name and email are required")
			}
			if err := ctx.playerRepo.AddPlayer(playerName, playerEmail); err != nil {
				return fmt.Errorf("adding player: %w", err)
			}
			return nil
		},
	}
	listPlayers = &cobra.Command{
		Use:   "list",
		Short: "List all players",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := getAppContext(cmd)
			if ctx == nil {
				return fmt.Errorf("app context not found")
			}
			players, err := ctx.playerRepo.ListPlayers()
			if err != nil {
				return fmt.Errorf("listing players: %w", err)
			}
			for _, player := range players {
				fmt.Printf("%s <%s>\n", player.Name, player.Email)
			}
			return nil
		},
	}
)

func init() {
	addPlayer.Flags().StringVarP(&playerName, "name", "n", "", "name of the player")
	addPlayer.Flags().StringVarP(&playerEmail, "email", "e", "", "email of the player")
	addPlayer.MarkFlagRequired("name")
	addPlayer.MarkFlagRequired("email")

	rootCmd.AddCommand(playersCmd)
	playersCmd.AddCommand(addPlayer)
	playersCmd.AddCommand(listPlayers)
}
