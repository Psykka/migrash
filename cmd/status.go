package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status check if there are pending or executed migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Status")
	},
}
