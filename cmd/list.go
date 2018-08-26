package cmd

import (
	"github.com/jpignata/cwl/pkg/list"
	"github.com/spf13/cobra"
)

var flagMaxResults int

var listCmd = &cobra.Command{
	Use:   "list [group name]",
	Short: "List log groups and log streams within a group",
	Long:  "List log groups and log streams within a group",
	Args:  cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
	},
	Run: func(cmd *cobra.Command, args []string) {
		list := list.List{
			MaxResults: flagMaxResults,
		}

		if len(args) == 1 {
			list.GroupName = args[0]
		}

		list.Run()
	},
}

func init() {
	listCmd.Flags().IntVarP(&flagMaxResults, "max-results", "l", 50, "Maximum number of results to return")

	rootCmd.AddCommand(listCmd)
}
