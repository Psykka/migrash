package core

import (
	"fmt"
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

func runSqlFile(filePath string, db *sqlx.DB) {
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

func Up(config *config.Config, db *sqlx.DB) {
	dirs := utils.ReadDir(config.MigrationDir)

	var pendingMigrations int

	for _, dir := range dirs {
		name := dir.Name()
		if !database.MigrationExists(name, config, db) {
			fmt.Println("Applying migration", name)

			filePath := fmt.Sprintf("%s/%s/%s", config.MigrationDir, name, up)
			runSqlFile(filePath, db)

			database.InsertMigration(name, config, db)
			pendingMigrations++
		}
	}

	if pendingMigrations == 0 {
		fmt.Println("No pending migrations!")
	} else {
		fmt.Printf("Applied %d migration(s) successfully!\n", pendingMigrations)
	}
}

func Down(config *config.Config, db *sqlx.DB) {
	migration := database.GetLastMigration(config, db)

	if migration != nil {
		fmt.Println("Rolling back migration", migration.Name)

		filePath := fmt.Sprintf("%s/%s/%s", config.MigrationDir, migration.Name, down)
		runSqlFile(filePath, db)

		database.RemoveMigration(migration.Name, config, db)

		fmt.Println("Migration rolled back successfully!")
	} else {
		fmt.Println("No migrations to rollback!")
	}
}

func Reset(config *config.Config, db *sqlx.DB) {
	dirs := utils.ReadDir(config.MigrationDir)

	var pendingMigrations int

	for _, dir := range dirs {
		name := dir.Name()
		if database.MigrationExists(name, config, db) {
			fmt.Println("Rolling back migration", name)

			filePath := fmt.Sprintf("%s/%s/%s", config.MigrationDir, name, down)
			runSqlFile(filePath, db)

			database.RemoveMigration(name, config, db)

			pendingMigrations++
		}
	}

	if pendingMigrations == 0 {
		fmt.Println("No pending migrations to reset!")
	} else {
		fmt.Printf("Reset %d migration(s) successfully!\n", pendingMigrations)
	}
}

func Status(config *config.Config, db *sqlx.DB) {
	dirs := utils.ReadDir(config.MigrationDir)

	var pendingMigrations int

	for _, dir := range dirs {
		name := dir.Name()
		if !database.MigrationExists(name, config, db) {
			fmt.Println("Pending migration", name)
			pendingMigrations++
		}
	}

	if pendingMigrations == 0 {
		fmt.Println("No pending migrations!")
	} else {
		fmt.Printf("There are %d pending migration(s)\n", pendingMigrations)
	}
}

func Redo(config *config.Config, db *sqlx.DB) {
	Down(config, db)
	Up(config, db)

	fmt.Println("Migrations rolled back and executed successfully!")
}
