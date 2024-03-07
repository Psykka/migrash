package core

import (
	"fmt"
	"io/fs"
	"migrash/internal/config"
	"migrash/pkg/database"
	"migrash/pkg/utils"
	"os"

	"github.com/jmoiron/sqlx"
)

const (
	up   = "up.sql"
	down = "down.sql"
)

func runSqlFilesInDir(dirs []fs.FileInfo, config *config.Config, db *sqlx.DB, action string) {
	for _, dir := range dirs {
		name := dir.Name()
		if !database.MigrationExists(name, config, db) {
			fmt.Println("Running migration", name)

			filePath := fmt.Sprintf("%s/%s/%s", config.MigrationDir, name, action)
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
}

func Up(config *config.Config, db *sqlx.DB) {
	dirs := utils.ReadDir(config.MigrationDir)

	runSqlFilesInDir(dirs, config, db, up)
}

func Down(config *config.Config, db *sqlx.DB) {
	dirs := utils.ReadDir(config.MigrationDir)

	runSqlFilesInDir(dirs, config, db, down)
}
