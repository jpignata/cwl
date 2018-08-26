package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use: "cwl",
}

func Execute() {
	rootCmd.Execute()
}
