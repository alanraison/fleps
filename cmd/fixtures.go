package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	fixtures = &cobra.Command{
		Use:   "fixtures",
		Short: "Manage fixtures",
	}
	add = &cobra.Command{
		Use:   "add",
		Short: "Add fixtures from a CSV file",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Add fixtures")
		},
	}
	list = &cobra.Command{
		Use:   "list",
		Short: "List all fixtures",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("List fixtures")
		},
	}
	addResult = &cobra.Command{
		Use:   "add-result",
		Short: "Add a result for a fixture",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Add result")
		},
	}
)

func init() {
	fixtures.AddCommand(add)
	fixtures.AddCommand(list)
	fixtures.AddCommand(addResult)
}