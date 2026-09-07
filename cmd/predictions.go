package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	predictionsCmd = &cobra.Command{
		Use:   "predictions",
		Short: "Manage predictions",
	}
	addPrediction = &cobra.Command{
		Use:   "add",
		Short: "Add a new prediction",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Add prediction")
		},
	}
	listPredictions = &cobra.Command{
		Use:   "list",
		Short: "List all predictions",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("List predictions")
		},
	}
)

func init() {
	rootCmd.AddCommand(predictionsCmd)
	predictionsCmd.AddCommand(addPrediction)
	predictionsCmd.AddCommand(listPredictions)
}
