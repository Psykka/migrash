package cmd

import (
	"fmt"
	"migrash/internal/config"
	"migrash/pkg/database"
	"migrash/pkg/utils"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upCmd)
}

var db *sqlx.DB

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	PostRun: func(cmd *cobra.Command, args []string) {
		db.Close()
	},
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.ParseConfig()

		if err != nil {
			panic(err)
		}

		db = database.Connect(config)

		if err != nil {
			panic(err)
		}

		dirs := utils.ReadDir(config.MigrationDir)

		for _, dir := range dirs {
			name := dir.Name()
			if !database.MigrationExists(name, config, db) {
				fmt.Println("Running migration", name)

				filePath := fmt.Sprintf("%s/%s/up.sql", config.MigrationDir, name)
				fileString, err := os.ReadFile(filePath)
				if err != nil {
					panic(err)
				}

				tx, err := db.Begin()
				if err != nil {
					panic(err)
				}

				_, err = tx.Exec(string(fileString))
				if err != nil {
					panic(err)
				}

				err = tx.Commit()
				if err != nil {
					panic(err)
				}
			}
		}

		fmt.Println("Migrations executed successfully")
	},
}
