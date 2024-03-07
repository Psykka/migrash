package cmd

import (
	"fmt"
	"migrash/internal/config"
	"migrash/pkg/core"
	"migrash/pkg/database"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downCmd)
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.ParseConfig()

		if err != nil {
			panic(err)
		}

		db := database.Connect(config)
		core.Down(config, db)

		fmt.Println("Migrations executed successfully")
		db.Close()
	},
}
