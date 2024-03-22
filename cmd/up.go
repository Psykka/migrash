package cmd

import (
	"migrash/internal/config"
	"migrash/pkg/core"
	"migrash/pkg/database"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.ParseConfig()

		if err != nil {
			panic(err)
		}

		db := database.Connect(config)
		core.Up(config, db)

		db.Close()
	},
}
