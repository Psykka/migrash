package cmd

import (
	"migrash/internal/config"
	"migrash/pkg/core"
	"migrash/pkg/database"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(redoCmd)
}

var redoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Rollback the last migration and then apply it again",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.ParseConfig()

		if err != nil {
			panic(err)
		}

		db := database.Connect(config)
		core.Redo(config, db)

		db.Close()
	},
}
