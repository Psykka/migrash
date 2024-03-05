package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Up")
	},
}
