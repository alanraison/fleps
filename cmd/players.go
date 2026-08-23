package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	players = &cobra.Command{
		Use:   "players",
		Short: "Manage players",
	}
	addPlayer = &cobra.Command{
		Use:   "add",
		Short: "Add a new player",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Add player")
		},
	}
	listPlayers = &cobra.Command{
		Use:   "list",
		Short: "List all players",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("List players")
		},
	}
)

func init() {
	players.AddCommand(addPlayer)
	players.AddCommand(listPlayers)
}
