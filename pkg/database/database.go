package database

import (
	"fmt"

	"migrash/internal/config"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

func getConnection(dbms string, url string) *sqlx.DB {
	var db *sqlx.DB
	var err error

	switch dbms {
	case "mysql", "mariadb":
		db, err = sqlx.Open("mysql", url)
	case "postgres":
		db, err = sqlx.Open("pgx", url)
	case "sqlite3":
		db, err = sqlx.Open("sqlite3", url)
	}

	if err != nil {
		panic(err)
	}

	return db
}

func getDatabaseConnection(dbms string, url string, name string) *sqlx.DB {
	var db *sqlx.DB
	var err error

	switch dbms {
	case "mysql", "mariadb":
		db, err = sqlx.Open("mysql", url+name)
	case "postgres":
		db, err = sqlx.Open("pgx", url+name)
	case "sqlite3":
		db, err = sqlx.Open("sqlite3", url)
	}

	if err != nil {
		panic(err)
	}

	return db
}

func createDatabase(config *config.Config, db *sqlx.DB) {
	switch config.Database.DBMS {
	case "mysql", "mariadb":
		if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + config.Database.Name); err != nil {
			panic(err)
		}
	case "postgres":
		var exists bool
		err := db.Get(&exists, "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)", config.Database.Name)
		if err != nil {
			panic(err)
		}

		if !exists {
			if _, err := db.Exec("CREATE DATABASE " + config.Database.Name); err != nil {
				panic(err)
			}
		}
	case "sqlite3":
		return
	}
}

func getDisplayInfo(config *config.Config) string {
	var info string

	switch config.Database.DBMS {
	case "mysql", "mariadb", "postgres":
		info = fmt.Sprintf("Starting a connection with %s on server %s@%s:%s", config.Database.DBMS, config.Database.User, config.Database.Host, config.Database.Port)
	case "sqlite3":
		info = fmt.Sprintf("Starting a connection with %s on %s", config.Database.DBMS, config.Database.Url)
	}

	return info
}

func Connect(config *config.Config) *sqlx.DB {
	databaseVersion := getDisplayInfo(config)
	fmt.Println(databaseVersion)

	db := getConnection(config.Database.DBMS, config.Database.Url)

	createDatabase(config, db)

	db.Close()

	db = getDatabaseConnection(config.Database.DBMS, config.Database.Url, config.Database.Name)
	fmt.Println("Connected to", config.Database.DBMS)

	return db
}
