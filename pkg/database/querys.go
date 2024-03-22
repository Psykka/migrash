package database

import (
	"fmt"
	"migrash/internal/config"
	"time"

	"github.com/jmoiron/sqlx"
)

type Migration struct {
	Id        int
	Name      string
	CreatedAt string `db:"created_at"`
}

var createTable = map[string]string{
	"sqlite3":  "CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY, name TEXT NOT NULL, created_at TEXT NOT NULL);",
	"mysql":    "CREATE TABLE IF NOT EXISTS %s (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, created_at DATETIME NOT NULL);",
	"postgres": "CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, created_at TIMESTAMP NOT NULL);",
}

var insert = map[string]string{
	"sqlite3":  "INSERT INTO %s (name, created_at) VALUES (?, ?);",
	"mysql":    "INSERT INTO %s (name, created_at) VALUES (?, ?);",
	"postgres": "INSERT INTO %s (name, created_at) VALUES ($1, $2);",
}

var delete = map[string]string{
	"sqlite3":  "DELETE FROM %s WHERE name = ?;",
	"mysql":    "DELETE FROM %s WHERE name = ?;",
	"postgres": "DELETE FROM %s WHERE name = $1;",
}

var exists = map[string]string{
	"sqlite3":  "SELECT COUNT(*) FROM %s WHERE name = ?;",
	"mysql":    "SELECT COUNT(*) FROM %s WHERE name = ?;",
	"postgres": "SELECT COUNT(*) FROM %s WHERE name = $1;",
}

func createMigrationTable(config *config.Config, db *sqlx.DB) {
	query := fmt.Sprintf(createTable[config.Database.DBMS], config.MigrationTable)
	db.MustExec(query)
}

func InsertMigration(name string, config *config.Config, db *sqlx.DB) {
	createMigrationTable(config, db)

	query := fmt.Sprintf(insert[config.Database.DBMS], config.MigrationTable)
	timestamp := time.Now()

	db.MustExec(query, name, timestamp)
}

func RemoveMigration(name string, config *config.Config, db *sqlx.DB) {
	createMigrationTable(config, db)

	query := fmt.Sprintf(delete[config.Database.DBMS], config.MigrationTable)

	db.MustExec(query, name)
}

func GetLastMigration(config *config.Config, db *sqlx.DB) *Migration {
	createMigrationTable(config, db)

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id DESC LIMIT 1;", config.MigrationTable)

	var migration Migration
	err := db.Get(&migration, query)

	if err != nil {
		return nil
	}

	return &migration
}

func MigrationExists(name string, config *config.Config, db *sqlx.DB) bool {
	createMigrationTable(config, db)

	query := fmt.Sprintf(exists[config.Database.DBMS], config.MigrationTable)
	row := db.QueryRow(query, name)

	var exists bool
	err := row.Scan(&exists)

	if err != nil {
		panic(err)
	}

	return exists
}
