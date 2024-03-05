package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new migration file and a directory if it does not exist",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Create")
	},
}
