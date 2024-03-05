package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downCmd)
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Down")
	},
}
