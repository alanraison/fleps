package cmd

import (
	"fmt"
	"os"
	"time"

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
			fs, err := c.ReadFixtureRows(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading fixture rows: %v\n", err)
				return
			}
			database.AddFixtures(fs)
		},
	}
	listFrom, listTo string
	listTeams        []string
	list             = &cobra.Command{
		Use:   "list",
		Short: "List all fixtures",
		Run: func(cmd *cobra.Command, args []string) {
			fromDate, err := time.Parse("2006-01-02", listFrom)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid from date: %v\n", err)
				return
			}
			toDate, err := time.Parse("2006-01-02", listTo)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid to date: %v\n", err)
				return
			}
			f, err := database.ListFixtures(fromDate, toDate, listTeams)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing fixtures: %v\n", err)
				return
			}
			for _, fixture := range f {
				fmt.Printf("%v %v %v\n", fixture.Date.Local().Format("02/01 15:04"), fixture.HomeTeam, fixture.AwayTeam)
			}
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
	list.Flags().StringVarP(&listFrom, "from", "f", "", "from date")
	list.Flags().StringVarP(&listTo, "to", "t", "", "to date")
	list.Flags().StringArrayVarP(&listTeams, "teams", "m", []string{}, "filter by teams")
	err := list.MarkFlagRequired("from")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'from' flag as required: %v\n", err)
	}
	err = list.MarkFlagRequired("to")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'to' flag as required: %v\n", err)
	}

	fixtures.AddCommand(add)
	fixtures.AddCommand(list)
	fixtures.AddCommand(addResult)
}
