package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Rollback all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Reset")
	},
}
