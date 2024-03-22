package cmd

import (
	"migrash/internal/config"
	"migrash/pkg/core"
	"migrash/pkg/database"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status check if there are pending or executed migrations",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.ParseConfig()

		if err != nil {
			panic(err)
		}

		db := database.Connect(config)
		core.Status(config, db)

		db.Close()
	},
}
