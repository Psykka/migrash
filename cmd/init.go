package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create config file for migrash",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Init")
	},
}
